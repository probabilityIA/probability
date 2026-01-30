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
		// Config endpoint - público (reemplaza variable de entorno)
		shopifyGroup.GET("/config", h.GetConfigHandler)

		// Auth endpoints
		shopifyGroup.POST("/auth/login", h.LoginWithSessionTokenHandler)

		// OAuth endpoints - requieren autenticación JWT
		shopifyGroup.POST("/connect", middleware.JWT(), h.InitiateOAuthHandler)

		// Webhook endpoint - sin autenticación JWT (Shopify valida con HMAC)
		shopifyGroup.POST("/webhook", h.WebhookHandler)
		shopifyGroup.POST("/webhook/:integration_id", h.WebhookHandler) // Alternativa con ID en path

		// Compliance webhook unificado (OBLIGATORIO para Shopify App Store)
		// Maneja todos los webhooks de GDPR/CCPA en un solo endpoint
		shopifyGroup.POST("/webhooks/compliance", h.ComplianceWebhookHandler)

		// Endpoints individuales de compliance (opcional, para compatibilidad)
		shopifyGroup.POST("/webhooks/customers/data_request", h.CustomerDataRequestHandler)
		shopifyGroup.POST("/webhooks/customers/redact", h.CustomerRedactHandler)
		shopifyGroup.POST("/webhooks/shop/redact", h.ShopRedactHandler)
	}

	// Callback endpoint - sin autenticación JWT (validación por state y HMAC)
	router.GET("/shopify/callback", h.OAuthCallbackHandler)
}
