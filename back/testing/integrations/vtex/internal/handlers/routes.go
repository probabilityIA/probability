package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "mock": "vtex", "account": MockAccountName})
}

func requireVTEXAuth(c *gin.Context) bool {
	if c.GetHeader("X-VTEX-API-AppKey") == "" || c.GetHeader("X-VTEX-API-AppToken") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "AUT0001",
				"message": "Credenciales VTEX invalidas o ausentes",
			},
		})
		return false
	}
	return true
}

func (h *Handler) sortedSKUIDs() []int {
	ids := make([]int, 0, len(h.skus))
	for id := range h.skus {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	return ids
}

func (h *Handler) handleGetProductAndSkuIds(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	h.mu.Lock()
	ids := h.sortedSKUIDs()
	data := make(map[string][]int, len(ids))
	for _, id := range ids {
		item := h.skus[id]
		data[strconv.Itoa(item.ProductID)] = []int{item.ID}
	}
	total := len(ids)
	h.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"range": gin.H{
			"total": total,
			"from":  0,
			"to":    total,
		},
	})
}

func (h *Handler) handleListSKUIDs(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "50"))
	if pageSize <= 0 {
		pageSize = 50
	}

	h.mu.Lock()
	ids := h.sortedSKUIDs()
	h.mu.Unlock()

	start := (page - 1) * pageSize
	if start >= len(ids) {
		c.JSON(http.StatusOK, []int{})
		return
	}
	end := start + pageSize
	if end > len(ids) {
		end = len(ids)
	}

	c.JSON(http.StatusOK, ids[start:end])
}

func (h *Handler) handleGetSKU(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	h.mu.Lock()
	item, ok := h.skus[id]
	h.mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "SKU0001", "message": "SKU no encontrado"}})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) handleSKUIDByRefID(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	refID := strings.TrimSpace(c.Param("refid"))

	h.mu.Lock()
	var found *sku
	for _, id := range h.sortedSKUIDs() {
		if strings.EqualFold(h.skus[id].RefID, refID) && refID != "" {
			found = h.skus[id]
			break
		}
	}
	h.mu.Unlock()

	if found == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "SKU0002", "message": "RefId no encontrado"}})
		return
	}
	c.JSON(http.StatusOK, found.ID)
}

func (h *Handler) handleSellerSKUSearch(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	externalID := strings.TrimSpace(c.Query("externalid"))

	out := make([]gin.H, 0)
	h.mu.Lock()
	for _, id := range h.sortedSKUIDs() {
		item := h.skus[id]
		if item.RefID != "" && strings.EqualFold(item.RefID, externalID) {
			out = append(out, gin.H{"id": strconv.Itoa(item.ID), "externalId": item.RefID})
		}
	}
	h.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *Handler) handleListWarehouses(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	h.mu.Lock()
	list := h.warehouses
	h.mu.Unlock()

	c.JSON(http.StatusOK, list)
}

func (h *Handler) handleUpdateInventory(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	skuID, _ := strconv.Atoi(c.Param("sku"))
	warehouseID := c.Param("warehouse")

	var body struct {
		Quantity          int  `json:"quantity"`
		UnlimitedQuantity bool `json:"unlimitedQuantity"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INV0001", "message": err.Error()}})
		return
	}

	h.mu.Lock()
	_, exists := h.skus[skuID]
	if exists {
		h.stock[stockKey(skuID, warehouseID)] = body.Quantity
	}
	h.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "SKU0001", "message": "SKU no encontrado"}})
		return
	}

	h.logger.Info().Msgf("VTEX mock: sku %d en bodega %s queda en %d", skuID, warehouseID, body.Quantity)
	c.Status(http.StatusOK)
}

func (h *Handler) sampleOrderSummary(orderID int) gin.H {
	now := time.Now().UTC().Format(time.RFC3339)
	return gin.H{
		"orderId":      fmt.Sprintf("%s-01", mockOrderCode(orderID)),
		"sequence":     strconv.Itoa(orderID),
		"status":       "ready-for-handling",
		"creationDate": now,
		"lastChange":   now,
		"totalValue":   9725000,
		"currencyCode": "COP",
		"origin":       "Marketplace",
	}
}

func mockOrderCode(orderID int) string {
	return fmt.Sprintf("%s-%d", MockAccountName, orderID)
}

func (h *Handler) sampleOrderDetail(orderID string) gin.H {
	now := time.Now().UTC().Format(time.RFC3339)

	h.mu.Lock()
	var first *sku
	ids := h.sortedSKUIDs()
	if len(ids) > 0 {
		first = h.skus[ids[0]]
	}
	h.mu.Unlock()

	refID := "VTEX-MOCK-001"
	name := "Camiseta Mock VTEX"
	skuID := "1001"
	productID := "91001"
	if first != nil {
		refID = first.RefID
		name = first.Name
		skuID = strconv.Itoa(first.ID)
		productID = strconv.Itoa(first.ProductID)
	}

	return gin.H{
		"orderId":            orderID,
		"sequence":           "1",
		"marketplaceOrderId": orderID,
		"status":             "ready-for-handling",
		"statusDescription":  "Listo para manejo",
		"value":              9725000,
		"totalItems":         1,
		"totalDiscount":      0,
		"totalFreight":       800000,
		"creationDate":       now,
		"lastChange":         now,
		"items": []gin.H{
			{
				"uniqueId":        "u-" + skuID,
				"id":              skuID,
				"productId":       productID,
				"ean":             "",
				"refId":           refID,
				"name":            name,
				"skuName":         name,
				"quantity":        1,
				"price":           8925000,
				"listPrice":       8925000,
				"sellingPrice":    8925000,
				"tax":             0,
				"measurementUnit": "un",
				"unitMultiplier":  1.0,
			},
		},
		"shippingData": gin.H{
			"address": gin.H{
				"addressType":  "residential",
				"receiverName": "Cliente Prueba",
				"street":       "Calle 123",
				"number":       "45",
				"city":         "Bogota",
				"state":        "DC",
				"country":      "COL",
				"postalCode":   "110111",
				"neighborhood": "Centro",
			},
			"logisticsInfo": []gin.H{
				{"itemIndex": 0, "selectedSla": "Envio estandar", "shippingEstimate": "3bd"},
			},
			"selectedSla": "Envio estandar",
		},
		"paymentData": gin.H{
			"transactions": []gin.H{
				{
					"isActive": true,
					"payments": []gin.H{
						{"id": "pay-" + orderID, "paymentSystemName": "Contra entrega", "value": 9725000, "installments": 1},
					},
				},
			},
		},
		"clientProfileData": gin.H{
			"email":            "cliente.prueba@example.com",
			"firstName":        "Cliente",
			"lastName":         "Prueba",
			"document":         "1020304050",
			"documentType":     "cpf",
			"phone":            "+573001234567",
			"corporateName":    "",
			"isCorporate":      false,
			"userProfileId":    "up-mock",
			"customerClass":    "",
			"stateInscription": "",
		},
		"sellers": []gin.H{
			{"id": "1", "name": "Tienda Mock VTEX", "logo": ""},
		},
		"totals": []gin.H{
			{"id": "Items", "name": "Total de itens", "value": 8925000},
			{"id": "Shipping", "name": "Total do frete", "value": 800000},
		},
	}
}

func (h *Handler) handleListOrders(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))
	if perPage <= 0 {
		perPage = 15
	}

	list := make([]gin.H, 0)
	if page == 1 {
		h.mu.Lock()
		base := h.nextOrderID
		h.mu.Unlock()
		list = append(list, h.sampleOrderSummary(base), h.sampleOrderSummary(base+1))
	}

	c.JSON(http.StatusOK, gin.H{
		"list": list,
		"paging": gin.H{
			"total":       2,
			"pages":       1,
			"currentPage": page,
			"perPage":     perPage,
		},
	})
}

func (h *Handler) handleGetOrder(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}
	c.JSON(http.StatusOK, h.sampleOrderDetail(c.Param("id")))
}

func (h *Handler) handleGetHook(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	h.mu.Lock()
	hook := h.hook
	h.mu.Unlock()

	if hook == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "HOOK0001", "message": "hook no configurado"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filter": gin.H{"type": "FromWorkflow", "status": hook.Statuses},
		"hook":   gin.H{"url": hook.URL, "headers": hook.Headers},
	})
}

func (h *Handler) handleSetHook(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	var body struct {
		Filter struct {
			Type   string   `json:"type"`
			Status []string `json:"status"`
		} `json:"filter"`
		Hook struct {
			URL     string            `json:"url"`
			Headers map[string]string `json:"headers"`
		} `json:"hook"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "HOOK0002", "message": err.Error()}})
		return
	}

	h.mu.Lock()
	h.hook = &hookConfig{URL: body.Hook.URL, Headers: body.Hook.Headers, Statuses: body.Filter.Status}
	h.mu.Unlock()

	h.logger.Info().Msgf("VTEX mock: hook configurado hacia %s (%d estados)", body.Hook.URL, len(body.Filter.Status))
	c.JSON(http.StatusOK, gin.H{
		"filter": gin.H{"type": body.Filter.Type, "status": body.Filter.Status},
		"hook":   gin.H{"url": body.Hook.URL, "headers": body.Hook.Headers},
	})
}

func (h *Handler) handleDeleteHook(c *gin.Context) {
	if !requireVTEXAuth(c) {
		return
	}

	h.mu.Lock()
	h.hook = nil
	h.mu.Unlock()

	c.Status(http.StatusNoContent)
}

func (h *Handler) handleSimulateOrder(c *gin.Context) {
	state := c.DefaultQuery("state", "ready-for-handling")

	h.mu.Lock()
	orderID := h.nextOrderID
	h.nextOrderID++
	hook := h.hook
	h.mu.Unlock()

	target := c.Query("target")
	headers := map[string]string{}
	if target == "" && hook != nil {
		target = hook.URL
		headers = hook.Headers
	}
	if target == "" {
		target = h.webhookURL
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no hay hook configurado ni target indicado"})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	payload := map[string]interface{}{
		"Domain":        "Fulfillment",
		"OrderId":       fmt.Sprintf("%s-01", mockOrderCode(orderID)),
		"State":         state,
		"LastState":     "order-created",
		"LastChange":    now,
		"CurrentChange": now,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "no se pudo entregar el webhook", "error": err.Error(), "target": target})
		return
	}
	defer resp.Body.Close()

	h.logger.Info().Msgf("VTEX mock: webhook %s enviado a %s (status %d)", state, target, resp.StatusCode)

	c.JSON(http.StatusOK, gin.H{
		"message":         "webhook enviado",
		"order_id":        payload["OrderId"],
		"state":           state,
		"target":          target,
		"response_status": resp.StatusCode,
	})
}

func (h *Handler) handleSeedProducts(c *gin.Context) {
	var body struct {
		Reset    bool `json:"reset"`
		Products []struct {
			ID    int    `json:"id"`
			RefID string `json:"ref_id"`
			Name  string `json:"name"`
			Stock int    `json:"stock"`
		} `json:"products"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	created := 0

	h.mu.Lock()
	if body.Reset {
		h.skus = make(map[int]*sku)
		h.stock = make(map[string]int)
	}
	for _, p := range body.Products {
		id := p.ID
		if id == 0 {
			id = h.nextSKUID
			h.nextSKUID++
		}
		h.skus[id] = &sku{
			ID:              id,
			ProductID:       id + 90000,
			Name:            p.Name,
			RefID:           p.RefID,
			IsActive:        true,
			WeightKg:        1.2,
			Length:          30,
			Width:           20,
			Height:          10,
			MeasurementUnit: "un",
		}
		h.stock[stockKey(id, h.warehouses[0].ID)] = p.Stock
		created++
		if id >= h.nextSKUID {
			h.nextSKUID = id + 1
		}
	}
	total := len(h.skus)
	h.mu.Unlock()

	h.logger.Info().Msgf("VTEX mock: seed con %d skus (total %d)", created, total)
	c.JSON(http.StatusOK, gin.H{"seeded_skus": created, "total_skus": total})
}
