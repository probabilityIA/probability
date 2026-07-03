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
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "mock": "woocommerce"})
}

func (h *Handler) handleSystemStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"environment": gin.H{
			"home_url":            "http://woocommerce-mock.local",
			"site_url":            "http://woocommerce-mock.local",
			"wc_version":          "9.0.0-mock",
			"wp_version":          "6.5-mock",
			"external_object_cache": nil,
		},
		"settings": gin.H{"currency": "COP"},
	})
}

func sampleOrder(id int64) gin.H {
	return gin.H{
		"id":           id,
		"number":       strconv.FormatInt(id, 10),
		"status":       "processing",
		"currency":     "COP",
		"date_created": time.Now().UTC().Format(time.RFC3339),
		"total":        "150000.00",
		"billing": gin.H{
			"first_name": "Cliente",
			"last_name":  "Prueba",
			"email":      "cliente.prueba@example.com",
			"phone":      "3001234567",
			"address_1":  "Calle 123",
			"city":       "Bogota",
		},
		"line_items": []gin.H{
			{
				"id":         1,
				"name":       "Producto de prueba",
				"product_id": 100,
				"quantity":   2,
				"sku":        "MOCK-SKU-001",
				"total":      "150000.00",
			},
		},
	}
}

func (h *Handler) handleListOrders(c *gin.Context) {
	c.JSON(http.StatusOK, []gin.H{sampleOrder(5001), sampleOrder(5002)})
}

func (h *Handler) handleGetOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id == 0 {
		id = 5001
	}
	c.JSON(http.StatusOK, sampleOrder(id))
}

func (h *Handler) handleCreateProduct(c *gin.Context) {
	var body map[string]interface{}
	_ = c.ShouldBindJSON(&body)

	h.mu.Lock()
	id := h.nextProdID
	h.nextProdID++
	h.mu.Unlock()

	if body == nil {
		body = map[string]interface{}{}
	}
	body["id"] = id
	c.JSON(http.StatusCreated, body)
}

func (h *Handler) handleUpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var body map[string]interface{}
	_ = c.ShouldBindJSON(&body)
	if body == nil {
		body = map[string]interface{}{}
	}
	body["id"] = id
	c.JSON(http.StatusOK, body)
}

func (h *Handler) handleListWebhooks(c *gin.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()
	list := make([]*webhook, 0, len(h.webhooks))
	for _, w := range h.webhooks {
		list = append(list, w)
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) handleCreateWebhook(c *gin.Context) {
	var body struct {
		Name        string `json:"name"`
		Topic       string `json:"topic"`
		DeliveryURL string `json:"delivery_url"`
		Secret      string `json:"secret"`
		Status      string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "invalid", "message": err.Error()})
		return
	}

	h.mu.Lock()
	id := h.nextID
	h.nextID++
	status := body.Status
	if status == "" {
		status = "active"
	}
	w := &webhook{
		ID:          id,
		Name:        body.Name,
		Status:      status,
		Topic:       body.Topic,
		DeliveryURL: body.DeliveryURL,
		Secret:      body.Secret,
		DateCreated: time.Now().UTC().Format(time.RFC3339),
	}
	h.webhooks[id] = w
	h.mu.Unlock()

	c.JSON(http.StatusCreated, w)
}

func (h *Handler) handleDeleteWebhook(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.mu.Lock()
	w, ok := h.webhooks[id]
	if ok {
		delete(h.webhooks, id)
	}
	h.mu.Unlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"code": "woocommerce_rest_webhook_invalid_id", "message": "webhook no existe"})
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *Handler) handleSimulateOrder(c *gin.Context) {
	topic := c.DefaultQuery("topic", "order.created")

	h.mu.Lock()
	id := h.nextOrder
	h.nextOrder++
	targets := make([]*webhook, 0)
	for _, w := range h.webhooks {
		if w.Topic == topic && w.Status == "active" {
			targets = append(targets, w)
		}
	}
	h.mu.Unlock()

	order := sampleOrder(id)
	body, _ := json.Marshal(order)

	fired := 0
	errs := []string{}
	client := &http.Client{Timeout: 15 * time.Second}
	for _, w := range targets {
		req, err := http.NewRequest(http.MethodPost, w.DeliveryURL, bytes.NewReader(body))
		if err != nil {
			errs = append(errs, w.DeliveryURL+": "+err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-WC-Webhook-Topic", topic)
		req.Header.Set("X-WC-Webhook-Source", "http://woocommerce-mock.local")
		if w.Secret != "" {
			mac := hmac.New(sha256.New, []byte(w.Secret))
			mac.Write(body)
			req.Header.Set("X-WC-Webhook-Signature", base64.StdEncoding.EncodeToString(mac.Sum(nil)))
		}
		resp, err := client.Do(req)
		if err != nil {
			errs = append(errs, w.DeliveryURL+": "+err.Error())
			continue
		}
		resp.Body.Close()
		fired++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("orden %d simulada (%s)", id, topic),
		"order_id": id,
		"fired":    fired,
		"errors":   errs,
	})
}
