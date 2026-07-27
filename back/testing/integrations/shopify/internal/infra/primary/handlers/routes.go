package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/app/usecases"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/admin/api/:version/orders.json", h.handleGetOrders)

	router.GET("/admin/api/:version/shop.json", h.handleGetShop)

	router.GET("/admin/api/:version/webhooks.json", h.handleListWebhooks)
	router.POST("/admin/api/:version/webhooks.json", h.handleCreateWebhook)
	router.DELETE("/admin/api/:version/webhooks/:webhook_id.json", h.handleDeleteWebhook)

	router.GET("/admin/api/:version/products.json", h.handleListProducts)
	router.POST("/admin/api/:version/products.json", h.handleCreateProduct)
	router.GET("/admin/api/:version/products/:product_id", h.handleGetProduct)

	router.GET("/admin/api/:version/locations.json", h.handleGetLocations)
	router.POST("/admin/api/:version/inventory_levels/set.json", h.handleSetInventoryLevel)

	router.GET("/health", h.handleHealth)

	router.GET("/mock/info", h.handleMockInfo)

	router.POST("/mock/generate", h.handleGenerate)
}

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "shopify-mock-api",
		"orders":  h.mockAPI.GetTotalOrders(),
	})
}

func (h *Handler) handleMockInfo(c *gin.Context) {
	config := h.mockAPI.GetBusinessConfig()
	currencyMode := "single-currency"
	if config != nil && config.IsDualCurrency() {
		currencyMode = config.ShopCurrency + "/" + config.PresentmentCurrency
	}

	c.JSON(http.StatusOK, gin.H{
		"total_orders":   h.mockAPI.GetTotalOrders(),
		"service":        "shopify-mock-api",
		"api_version":    "2024-10",
		"currency_mode":  currencyMode,
		"shop_currency":  config.ShopCurrency,
		"taxes_included": config.TaxesIncluded,
		"exchange_rate":  config.ExchangeRate,
		"descripcion":    "Mock del API REST de Shopify para pruebas de sincronizacion por lotes",
	})
}

func (h *Handler) handleGenerate(c *gin.Context) {
	var req struct {
		Count    int    `json:"count" binding:"required"`
		DateFrom string `json:"date_from"`
		DateTo   string `json:"date_to"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere 'count'. Opcional: 'date_from', 'date_to' (YYYY-MM-DD o RFC3339)"})
		return
	}

	if req.Count <= 0 || req.Count > 10000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "count debe estar entre 1 y 10000"})
		return
	}

	dateFrom := time.Now().AddDate(0, -6, 0)
	dateTo := time.Now()

	if req.DateFrom != "" {
		if t, err := parseFlexDate(req.DateFrom); err == nil {
			dateFrom = t
		}
	}
	if req.DateTo != "" {
		if t, err := parseFlexDate(req.DateTo); err == nil {
			dateTo = t
		}
	}

	h.mockAPI.GenerateOrders(req.Count, dateFrom, dateTo)

	c.JSON(http.StatusOK, gin.H{
		"mensaje":      fmt.Sprintf("Se generaron %d ordenes exitosamente", req.Count),
		"total_orders": h.mockAPI.GetTotalOrders(),
		"date_from":    dateFrom.Format(time.RFC3339),
		"date_to":      dateTo.Format(time.RFC3339),
	})
}

func (h *Handler) handleGetOrders(c *gin.Context) {
	accessToken := c.GetHeader("X-Shopify-Access-Token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": "[API] Clave de API no valida o token de acceso (verificar encabezado)",
		})
		return
	}

	params := h.parseQueryParams(c)

	h.logger.Info().
		Int("limit", params.Limit).
		Str("status", params.Status).
		Str("financial_status", params.FinancialStatus).
		Str("fulfillment_status", params.FulfillmentStatus).
		Msg(" GET /orders.json recibido")

	orders, hasNextPage := h.mockAPI.QueryOrders(params)

	c.Header("Content-Type", "application/json")

	if hasNextPage && len(orders) > 0 {
		nextOffset := params.PageOffset + params.Limit
		pageInfo := encodePageInfo(nextOffset)
		version := c.Param("version")
		host := c.Request.Host
		scheme := "https"
		if c.Request.TLS == nil {
			scheme = "http"
		}
		linkURL := fmt.Sprintf("%s://%s/admin/api/%s/orders.json?page_info=%s&limit=%d",
			scheme, host, version, pageInfo, params.Limit)
		c.Header("Link", fmt.Sprintf(`<%s>; rel="next"`, linkURL))
	}

	response := map[string]interface{}{
		"orders": orders,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) parseQueryParams(c *gin.Context) usecases.OrderQueryParams {
	params := usecases.OrderQueryParams{
		Limit:  250,
		Status: "any",
	}

	if pageInfo := c.Query("page_info"); pageInfo != "" {
		if offset, err := decodePageInfo(pageInfo); err == nil {
			params.PageOffset = offset
		}
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 250 {
				params.Limit = l
			}
		}
		return params
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 250 {
			params.Limit = l
		}
	}

	if status := c.Query("status"); status != "" {
		params.Status = status
	}

	if fs := c.Query("financial_status"); fs != "" {
		params.FinancialStatus = fs
	}

	if ffs := c.Query("fulfillment_status"); ffs != "" {
		params.FulfillmentStatus = ffs
	}

	if sinceIDStr := c.Query("since_id"); sinceIDStr != "" {
		if id, err := strconv.ParseInt(sinceIDStr, 10, 64); err == nil {
			params.SinceID = id
		}
	}

	if minStr := c.Query("created_at_min"); minStr != "" {
		if t, err := parseFlexDate(minStr); err == nil {
			params.CreatedAtMin = &t
		}
	}

	if maxStr := c.Query("created_at_max"); maxStr != "" {
		if t, err := parseFlexDate(maxStr); err == nil {
			params.CreatedAtMax = &t
		}
	}

	return params
}

func parseFlexDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formato de fecha no reconocido: %s", dateStr)
}

func encodePageInfo(offset int) string {
	data := map[string]int{"offset": offset}
	b, _ := json.Marshal(data)
	return base64.URLEncoding.EncodeToString(b)
}

func decodePageInfo(pageInfo string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(pageInfo)
	if err != nil {
		return 0, err
	}
	var data struct {
		Offset int `json:"offset"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return 0, err
	}
	return data.Offset, nil
}
