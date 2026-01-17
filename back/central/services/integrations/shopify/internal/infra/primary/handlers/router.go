package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

// RegisterRoutes registra las rutas del handler de Shopify
func (h *ShopifyHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	shopifyGroup := router.Group("/integrations/shopify")
	{
		// OAuth endpoints - requieren autenticación JWT
		shopifyGroup.POST("/connect", middleware.JWT(), h.InitiateOAuthHandler)

		// Webhook endpoint - sin autenticación JWT (Shopify valida con HMAC)
		shopifyGroup.POST("/webhook", h.WebhookHandler)
		shopifyGroup.POST("/webhook/:integration_id", h.WebhookHandler) // Alternativa con ID en path
	}

	// Callback endpoint - sin autenticación JWT (validación por state y HMAC)
	router.GET("/shopify/callback", h.OAuthCallbackHandler)
}
