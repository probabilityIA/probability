package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// mockWebhook representa un webhook almacenado en memoria.
type mockWebhook struct {
	ID        int64  `json:"id"`
	Address   string `json:"address"`
	Topic     string `json:"topic"`
	Format    string `json:"format"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var (
	webhookStore   = make([]mockWebhook, 0)
	webhookMu      sync.RWMutex
	webhookSeq     int64 = 900000000
)

// handleGetShop simula GET /admin/api/2024-10/shop.json (validación de token).
func (h *Handler) handleGetShop(c *gin.Context) {
	accessToken := c.GetHeader("X-Shopify-Access-Token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": "[API] Invalid API key or access token",
		})
		return
	}

	// Obtener currency de la config del business (USD para dual-currency, COP por defecto)
	shopCurrency := "COP"
	moneyFormat := "$ {{amount_no_decimals}}"
	config := h.mockAPI.GetBusinessConfig()
	if config != nil && config.ShopCurrency != "" {
		shopCurrency = config.ShopCurrency
		if shopCurrency == "USD" {
			moneyFormat = "${{amount}}"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"shop": gin.H{
			"id":                      1,
			"name":                    "Tienda Mock Probability",
			"email":                   "mock@probability.test",
			"domain":                  "mock-probability.myshopify.com",
			"province":                "Cundinamarca",
			"country":                 "CO",
			"address1":                "Calle 123",
			"zip":                     "110111",
			"city":                    "Bogotá",
			"phone":                   "+573001234567",
			"currency":                shopCurrency,
			"money_format":            moneyFormat,
			"plan_name":               "developer",
			"myshopify_domain":        "mock-probability.myshopify.com",
			"iana_timezone":           "America/Bogota",
			"primary_locale":          "es",
			"created_at":              "2024-01-01T00:00:00-05:00",
			"updated_at":              time.Now().Format(time.RFC3339),
		},
	})
}

// handleListWebhooks simula GET /admin/api/2024-10/webhooks.json.
func (h *Handler) handleListWebhooks(c *gin.Context) {
	accessToken := c.GetHeader("X-Shopify-Access-Token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": "[API] Invalid API key or access token",
		})
		return
	}

	webhookMu.RLock()
	defer webhookMu.RUnlock()

	h.logger.Info().
		Int("count", len(webhookStore)).
		Msg("📋 GET /webhooks.json - listando webhooks mock")

	c.JSON(http.StatusOK, gin.H{
		"webhooks": webhookStore,
	})
}

// handleCreateWebhook simula POST /admin/api/2024-10/webhooks.json.
func (h *Handler) handleCreateWebhook(c *gin.Context) {
	accessToken := c.GetHeader("X-Shopify-Access-Token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": "[API] Invalid API key or access token",
		})
		return
	}

	var req struct {
		Webhook struct {
			Topic   string `json:"topic"`
			Address string `json:"address"`
			Format  string `json:"format"`
		} `json:"webhook"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": gin.H{"webhook": []string{"is required"}},
		})
		return
	}

	webhookMu.Lock()
	defer webhookMu.Unlock()

	webhookSeq++
	now := time.Now().Format(time.RFC3339)

	wh := mockWebhook{
		ID:        webhookSeq + int64(rand.Intn(999999)),
		Address:   req.Webhook.Address,
		Topic:     req.Webhook.Topic,
		Format:    req.Webhook.Format,
		CreatedAt: now,
		UpdatedAt: now,
	}

	webhookStore = append(webhookStore, wh)

	h.logger.Info().
		Int64("webhook_id", wh.ID).
		Str("topic", wh.Topic).
		Str("address", wh.Address).
		Msg("✅ POST /webhooks.json - webhook creado")

	c.JSON(http.StatusCreated, gin.H{
		"webhook": wh,
	})
}

// handleDeleteWebhook simula DELETE /admin/api/2024-10/webhooks/:id.json.
func (h *Handler) handleDeleteWebhook(c *gin.Context) {
	accessToken := c.GetHeader("X-Shopify-Access-Token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": "[API] Invalid API key or access token",
		})
		return
	}

	idStr := c.Param("webhook_id")
	// Limpiar sufijo .json si gin lo incluye en el param
	idStr = stripJSONSuffix(idStr)

	webhookID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"errors": "Not Found",
		})
		return
	}

	webhookMu.Lock()
	defer webhookMu.Unlock()

	found := false
	for i, wh := range webhookStore {
		if wh.ID == webhookID {
			webhookStore = append(webhookStore[:i], webhookStore[i+1:]...)
			found = true
			h.logger.Info().
				Int64("webhook_id", webhookID).
				Msg("🗑️ DELETE /webhooks/:id.json - webhook eliminado")
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"errors": fmt.Sprintf("Webhook %d not found", webhookID),
		})
		return
	}

	c.Status(http.StatusOK)
}

func stripJSONSuffix(s string) string {
	if len(s) > 5 && s[len(s)-5:] == ".json" {
		return s[:len(s)-5]
	}
	return s
}
