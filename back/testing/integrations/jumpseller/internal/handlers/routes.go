package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const jumpsellerTimeLayout = "2006-01-02 15:04:05 MST"

func pathID(c *gin.Context, param string) int64 {
	raw := strings.TrimSuffix(c.Param(param), ".json")
	id, _ := strconv.ParseInt(raw, 10, 64)
	return id
}

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "mock": "jumpseller"})
}

func (h *Handler) handleStoreInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"store": gin.H{
			"code":        MockStoreCode,
			"name":        "Tienda Mock Jumpseller",
			"url":         "http://jumpseller-mock.local",
			"country":     "CO",
			"currency":    "COP",
			"hooks_token": MockHooksToken,
			"weight_unit": "kg",
		},
	})
}

func sampleOrder(id int64) gin.H {
	return gin.H{
		"order": gin.H{
			"id":                     id,
			"created_at":             time.Now().UTC().Format(jumpsellerTimeLayout),
			"status":                 "Paid",
			"currency":               "COP",
			"subtotal":               75000.0,
			"tax":                    14250.0,
			"shipping_tax":           0.0,
			"shipping":               8000.0,
			"shipping_required":      true,
			"total":                  97250.0,
			"discount":               0.0,
			"shipping_discount":      0.0,
			"shipment_status":        "requested",
			"shipping_method_id":     1,
			"shipping_method_name":   "Envio estandar",
			"payment_method_name":    "Contra entrega",
			"payment_method_type":    "cod",
			"payment_information":    "Pago contra entrega",
			"additional_information": "Orden generada por el simulador",
			"customer": gin.H{
				"id":    "900",
				"name":  "Cliente Prueba",
				"email": "cliente.prueba@example.com",
				"phone": "3001234567",
				"ip":    "127.0.0.1",
			},
			"billing_address":  sampleAddress(),
			"shipping_address": sampleAddress(),
			"products": []gin.H{
				{
					"id":         100,
					"variant_id": 0,
					"sku":        "JS-MOCK-001",
					"name":       "Producto de prueba Jumpseller",
					"qty":        1,
					"price":      75000.0,
					"tax":        14250.0,
					"discount":   0.0,
					"weight":     1.5,
				},
			},
			"additional_fields": []gin.H{},
		},
	}
}

func sampleAddress() gin.H {
	return gin.H{
		"name":          "Cliente",
		"surname":       "Prueba",
		"taxid":         "1020304050",
		"address":       "Calle 123",
		"street_number": "45",
		"city":          "Bogota",
		"postal":        "110111",
		"region":        "Bogota D.C.",
		"country":       "Colombia",
		"country_code":  "CO",
		"region_code":   "DC",
	}
}

func (h *Handler) handleListOrders(c *gin.Context) {
	c.JSON(http.StatusOK, []gin.H{sampleOrder(5001), sampleOrder(5002)})
}

func (h *Handler) handleGetOrder(c *gin.Context) {
	id := pathID(c, "id")
	if id == 0 {
		id = 5001
	}
	c.JSON(http.StatusOK, sampleOrder(id))
}

func (h *Handler) handleUpdateOrder(c *gin.Context) {
	id := pathID(c, "id")
	var body map[string]interface{}
	_ = c.ShouldBindJSON(&body)

	h.logger.Info().Msgf("Jumpseller mock: orden %d actualizada con %v", id, body)
	c.JSON(http.StatusOK, sampleOrder(id))
}

func (h *Handler) handleListHooks(c *gin.Context) {
	h.mu.Lock()
	out := make([]gin.H, 0, len(h.hooks))
	for _, item := range h.hooks {
		out = append(out, gin.H{"hook": item})
	}
	h.mu.Unlock()

	c.JSON(http.StatusOK, out)
}

func (h *Handler) handleCreateHook(c *gin.Context) {
	var body struct {
		Hook struct {
			Event string `json:"event"`
			URL   string `json:"url"`
		} `json:"hook"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload invalido"})
		return
	}

	h.mu.Lock()
	id := h.nextHookID
	h.nextHookID++
	item := &hook{
		ID:        id,
		Name:      "probability-" + body.Hook.Event,
		Event:     body.Hook.Event,
		URL:       body.Hook.URL,
		CreatedAt: time.Now().UTC().Format(jumpsellerTimeLayout),
	}
	h.hooks[id] = item
	h.mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{"hook": item})
}

func (h *Handler) handleDeleteHook(c *gin.Context) {
	id := pathID(c, "id")

	h.mu.Lock()
	item, ok := h.hooks[id]
	if ok {
		delete(h.hooks, id)
	}
	h.mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "hook no existe"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"hook": item})
}

func (h *Handler) productsSnapshot() []gin.H {
	h.mu.Lock()
	defer h.mu.Unlock()

	out := make([]gin.H, 0, len(h.products))
	for _, item := range h.products {
		out = append(out, gin.H{"product": item})
	}
	return out
}

func (h *Handler) handleListProducts(c *gin.Context) {
	c.JSON(http.StatusOK, h.productsSnapshot())
}

func (h *Handler) handleSearchProducts(c *gin.Context) {
	query := strings.ToLower(strings.TrimSpace(c.Query("query")))

	h.mu.Lock()
	out := make([]gin.H, 0)
	for _, item := range h.products {
		if query == "" || strings.Contains(strings.ToLower(item.SKU), query) {
			out = append(out, gin.H{"product": item})
		}
	}
	h.mu.Unlock()

	c.JSON(http.StatusOK, out)
}

func (h *Handler) handleCreateProduct(c *gin.Context) {
	var body struct {
		Product struct {
			Name           string  `json:"name"`
			SKU            string  `json:"sku"`
			Price          float64 `json:"price"`
			Stock          int     `json:"stock"`
			StockUnlimited bool    `json:"stock_unlimited"`
			Status         string  `json:"status"`
			Weight         float64 `json:"weight"`
			Height         float64 `json:"height"`
			Width          float64 `json:"width"`
			Length         float64 `json:"length"`
		} `json:"product"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload invalido"})
		return
	}

	h.mu.Lock()
	id := h.nextProdID
	h.nextProdID++
	item := &product{
		ID:             id,
		Name:           body.Product.Name,
		SKU:            body.Product.SKU,
		Price:          body.Product.Price,
		Stock:          body.Product.Stock,
		StockUnlimited: body.Product.StockUnlimited,
		Status:         "available",
		Weight:         body.Product.Weight,
		Height:         body.Product.Height,
		Width:          body.Product.Width,
		Length:         body.Product.Length,
	}
	h.products[id] = item
	h.mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{"product": item})
}

func (h *Handler) handleUpdateProduct(c *gin.Context) {
	id := pathID(c, "id")

	var body struct {
		Product struct {
			Stock          int  `json:"stock"`
			StockUnlimited bool `json:"stock_unlimited"`
		} `json:"product"`
	}
	_ = c.ShouldBindJSON(&body)

	h.mu.Lock()
	item, ok := h.products[id]
	if ok {
		item.Stock = body.Product.Stock
		item.StockUnlimited = body.Product.StockUnlimited
	}
	h.mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "producto no existe"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": item})
}

func (h *Handler) handleUpdateVariant(c *gin.Context) {
	productID := pathID(c, "id")
	variantID := pathID(c, "vid")

	var body struct {
		Variant struct {
			Stock int `json:"stock"`
		} `json:"variant"`
	}
	_ = c.ShouldBindJSON(&body)

	h.logger.Info().Msgf("Jumpseller mock: variante %d del producto %d con stock %d", variantID, productID, body.Variant.Stock)

	c.JSON(http.StatusOK, gin.H{
		"variant": gin.H{
			"id":    variantID,
			"stock": body.Variant.Stock,
		},
	})
}

func (h *Handler) handleSimulateOrder(c *gin.Context) {
	event := c.DefaultQuery("event", "order_paid")

	h.mu.Lock()
	id := h.nextOrder
	h.nextOrder++
	targets := make([]*hook, 0)
	for _, item := range h.hooks {
		if item.Event == event {
			targets = append(targets, item)
		}
	}
	h.mu.Unlock()

	body, _ := json.Marshal(sampleOrder(id))

	fired := 0
	errs := make([]string, 0)
	client := &http.Client{Timeout: 15 * time.Second}

	for _, target := range targets {
		req, err := http.NewRequest(http.MethodPost, target.URL, bytes.NewReader(body))
		if err != nil {
			errs = append(errs, target.URL+": "+err.Error())
			continue
		}

		mac := hmac.New(sha256.New, []byte(MockHooksToken))
		mac.Write(body)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Jumpseller-Event", event)
		req.Header.Set("Jumpseller-Store-Code", MockStoreCode)
		req.Header.Set("Jumpseller-Hmac-Sha256", base64.StdEncoding.EncodeToString(mac.Sum(nil)))

		resp, err := client.Do(req)
		if err != nil {
			errs = append(errs, target.URL+": "+err.Error())
			continue
		}
		resp.Body.Close()
		fired++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("orden %d simulada (%s)", id, event),
		"order_id": id,
		"fired":    fired,
		"errors":   errs,
	})
}
