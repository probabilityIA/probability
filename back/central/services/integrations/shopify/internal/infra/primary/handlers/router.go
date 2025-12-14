package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/shared/log"
)

// RegisterRoutes registra las rutas del handler de Shopify
func (h *ShopifyHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	shopifyGroup := router.Group("/integrations/shopify")
	{
		// Webhook endpoint - sin autenticación JWT (Shopify valida con HMAC)
		shopifyGroup.POST("/webhook", h.WebhookHandler)
		shopifyGroup.POST("/webhook/:integration_id", h.WebhookHandler) // Alternativa con ID en path
	}
}
