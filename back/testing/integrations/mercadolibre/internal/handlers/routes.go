package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "mock": "mercadolibre"})
}

func requireBearer(c *gin.Context) bool {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") || strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid_token",
			"error":   "not_found",
			"status":  401,
		})
		return false
	}
	return true
}

func (h *Handler) handleToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"access_token":  MockAccessToken,
		"token_type":    "bearer",
		"expires_in":    21600,
		"scope":         "offline_access read write",
		"user_id":       MockSellerID,
		"refresh_token": MockRefreshToken,
	})
}

func (h *Handler) handleUserMe(c *gin.Context) {
	if !requireBearer(c) {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":                MockSellerID,
		"nickname":          "MOCKSELLER",
		"email":             "seller.mock@example.com",
		"first_name":        "Vendedor",
		"last_name":         "Mock",
		"site_id":           MockSiteID,
		"country_id":        "CO",
		"user_type":         "normal",
		"seller_reputation": gin.H{"level_id": "5_green"},
	})
}

func (h *Handler) handleSearchItems(c *gin.Context) {
	if !requireBearer(c) {
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if offset < 0 {
		offset = 0
	}

	h.mu.Lock()
	ids := h.sortedItemIDs()
	h.mu.Unlock()

	total := len(ids)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}

	c.JSON(http.StatusOK, gin.H{
		"seller_id": strconv.Itoa(MockSellerID),
		"results":   ids[offset:end],
		"paging": gin.H{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	})
}

func (h *Handler) handleMultigetItems(c *gin.Context) {
	if !requireBearer(c) {
		return
	}

	raw := strings.TrimSpace(c.Query("ids"))
	if raw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ids es requerido", "error": "bad_request", "status": 400})
		return
	}

	out := make([]gin.H, 0)
	h.mu.Lock()
	for _, id := range strings.Split(raw, ",") {
		id = strings.TrimSpace(id)
		it, ok := h.items[id]
		if !ok {
			out = append(out, gin.H{"code": http.StatusNotFound, "body": gin.H{"message": "Item not found", "error": "not_found"}})
			continue
		}
		out = append(out, gin.H{"code": http.StatusOK, "body": it})
	}
	h.mu.Unlock()

	c.JSON(http.StatusOK, out)
}

func (h *Handler) handleGetItem(c *gin.Context) {
	if !requireBearer(c) {
		return
	}

	h.mu.Lock()
	it, ok := h.items[c.Param("id")]
	h.mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found", "error": "not_found", "status": 404})
		return
	}
	c.JSON(http.StatusOK, it)
}

func (h *Handler) handleCreateItem(c *gin.Context) {
	if !requireBearer(c) {
		return
	}

	var body struct {
		Title             string          `json:"title"`
		Price             float64         `json:"price"`
		CurrencyID        string          `json:"currency_id"`
		AvailableQuantity int             `json:"available_quantity"`
		SiteID            string          `json:"site_id"`
		ListingTypeID     string          `json:"listing_type_id"`
		SellerCustomField string          `json:"seller_custom_field"`
		Attributes        []itemAttribute `json:"attributes"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "error": "bad_request", "status": 400})
		return
	}

	sku := body.SellerCustomField
	if sku == "" {
		for _, a := range body.Attributes {
			if a.ID == "SELLER_SKU" {
				sku = a.ValueName
			}
		}
	}

	h.mu.Lock()
	it := h.seed(body.Title, sku, body.Price, body.AvailableQuantity)
	if body.CurrencyID != "" {
		it.CurrencyID = body.CurrencyID
	}
	if body.ListingTypeID != "" {
		it.ListingTypeID = body.ListingTypeID
	}
	if body.SiteID != "" {
		it.SiteID = body.SiteID
	}
	h.mu.Unlock()

	h.logger.Info().Msgf("Meli mock: item creado %s (sku %s)", it.ID, sku)
	c.JSON(http.StatusCreated, it)
}

func (h *Handler) handleUpdateItem(c *gin.Context) {
	if !requireBearer(c) {
		return
	}

	var body struct {
		AvailableQuantity *int     `json:"available_quantity"`
		Price             *float64 `json:"price"`
		Status            *string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "error": "bad_request", "status": 400})
		return
	}

	id := c.Param("id")

	h.mu.Lock()
	it, ok := h.items[id]
	if ok {
		if body.AvailableQuantity != nil {
			it.AvailableQuantity = *body.AvailableQuantity
		}
		if body.Price != nil {
			it.Price = *body.Price
		}
		if body.Status != nil {
			it.Status = *body.Status
		}
	}
	h.mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found", "error": "not_found", "status": 404})
		return
	}

	h.logger.Info().Msgf("Meli mock: item %s queda con available_quantity %d", id, it.AvailableQuantity)
	c.JSON(http.StatusOK, it)
}

func (h *Handler) handleGetUserProductStock(c *gin.Context) {
	if !requireBearer(c) {
		return
	}

	h.mu.Lock()
	stock, ok := h.stocks[c.Param("id")]
	h.mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "user product not found", "error": "not_found", "status": 404})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_product_id": stock.UserProductID,
		"version":         stock.Version,
		"locations":       stock.Locations,
	})
}

func (h *Handler) handleUpdateUserProductStock(c *gin.Context) {
	if !requireBearer(c) {
		return
	}

	var body struct {
		Locations []map[string]interface{} `json:"locations"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "error": "bad_request", "status": 400})
		return
	}

	id := c.Param("id")

	h.mu.Lock()
	stock, ok := h.stocks[id]
	if ok {
		stock.Locations = body.Locations
		stock.Version = "v" + strconv.FormatInt(time.Now().Unix(), 10)
	}
	h.mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "user product not found", "error": "not_found", "status": 404})
		return
	}

	h.logger.Info().Msgf("Meli mock: user product %s actualizado con %d locations", id, len(body.Locations))
	c.JSON(http.StatusOK, gin.H{"user_product_id": id, "version": stock.Version, "locations": stock.Locations})
}

func (h *Handler) sortedItemIDs() []string {
	ids := make([]string, 0, len(h.items))
	for id := range h.items {
		ids = append(ids, id)
	}
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if ids[j] < ids[i] {
				ids[i], ids[j] = ids[j], ids[i]
			}
		}
	}
	return ids
}

func (h *Handler) sampleOrder(orderID int64) gin.H {
	h.mu.Lock()
	ids := h.sortedItemIDs()
	var first *item
	if len(ids) > 0 {
		first = h.items[ids[0]]
	}
	h.mu.Unlock()

	itemID := "MCO1000000001"
	title := "Producto Mock Meli"
	sku := "MELI-MOCK-001"
	price := 75000.0
	if first != nil {
		itemID = first.ID
		title = first.Title
		sku = first.SellerCustomField
		price = first.Price
	}

	now := time.Now().UTC().Format(time.RFC3339)

	return gin.H{
		"id":            orderID,
		"status":        "paid",
		"status_detail": nil,
		"date_created":  now,
		"date_closed":   now,
		"last_updated":  now,
		"total_amount":  price,
		"paid_amount":   price,
		"currency_id":   "COP",
		"buyer": gin.H{
			"id":         987654321,
			"nickname":   "COMPRADORMOCK",
			"email":      "comprador.mock@example.com",
			"first_name": "Comprador",
			"last_name":  "Mock",
			"phone":      gin.H{"area_code": "57", "number": "3001234567"},
		},
		"seller": gin.H{"id": MockSellerID, "nickname": "MOCKSELLER"},
		"order_items": []gin.H{
			{
				"item": gin.H{
					"id":                  itemID,
					"title":               title,
					"seller_custom_field": sku,
					"seller_sku":          sku,
					"variation_id":        nil,
				},
				"quantity":        1,
				"unit_price":      price,
				"full_unit_price": price,
				"currency_id":     "COP",
			},
		},
		"payments": []gin.H{
			{
				"id":                 orderID + 1,
				"status":             "approved",
				"status_detail":      "accredited",
				"transaction_amount": price,
				"total_paid_amount":  price,
				"shipping_cost":      0,
				"currency_id":        "COP",
				"date_approved":      now,
				"payment_method_id":  "visa",
				"payment_type":       "credit_card",
			},
		},
		"shipping": gin.H{"id": orderID + 2},
		"tags":     []string{"paid"},
	}
}

func (h *Handler) handleGetOrder(c *gin.Context) {
	if !requireBearer(c) {
		return
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id == 0 {
		id = 2000000001
	}
	c.JSON(http.StatusOK, h.sampleOrder(id))
}

func (h *Handler) handleSearchOrders(c *gin.Context) {
	if !requireBearer(c) {
		return
	}

	h.mu.Lock()
	id := h.nextOrderID
	h.mu.Unlock()

	orders := []gin.H{h.sampleOrder(id), h.sampleOrder(id + 10)}

	c.JSON(http.StatusOK, gin.H{
		"results": orders,
		"paging":  gin.H{"total": len(orders), "offset": 0, "limit": 50},
	})
}

func (h *Handler) handleSimulateNotification(c *gin.Context) {
	topic := c.DefaultQuery("topic", "orders_v2")

	h.mu.Lock()
	orderID := h.nextOrderID
	h.nextOrderID++
	h.mu.Unlock()

	resource := "/orders/" + strconv.FormatInt(orderID, 10)
	if topic == "items" {
		h.mu.Lock()
		ids := h.sortedItemIDs()
		h.mu.Unlock()
		if len(ids) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no hay items en el mock"})
			return
		}
		resource = "/items/" + ids[0]
	}

	payload := map[string]interface{}{
		"resource":       resource,
		"user_id":        MockSellerID,
		"topic":          topic,
		"application_id": 1234567890123456,
		"attempts":       1,
		"sent":           time.Now().UTC().Format(time.RFC3339),
		"received":       time.Now().UTC().Format(time.RFC3339),
	}

	target := c.Query("target")
	if target == "" {
		target = h.webhookURL
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "no se pudo entregar la notificacion", "error": err.Error(), "target": target})
		return
	}
	defer resp.Body.Close()

	h.logger.Info().Msgf("Meli mock: notificacion %s enviada a %s (status %d)", topic, target, resp.StatusCode)

	c.JSON(http.StatusOK, gin.H{
		"message":         "notificacion enviada",
		"topic":           topic,
		"resource":        resource,
		"target":          target,
		"response_status": resp.StatusCode,
	})
}
