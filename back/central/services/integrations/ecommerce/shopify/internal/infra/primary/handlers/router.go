package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

func (h *ShopifyHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	shopifyGroup := router.Group("/integrations/shopify")
	{
		shopifyGroup.GET("/config", h.GetConfigHandler)

		shopifyGroup.POST("/auth/login", h.LoginWithSessionTokenHandler)

		shopifyGroup.POST("/connect", middleware.JWT(), h.InitiateOAuthHandler)
		shopifyGroup.POST("/connect/custom", middleware.JWT(), h.InitiateCustomOAuthHandler)

		shopifyGroup.GET("/oauth/token", h.GetOAuthTokenHandler)

		shopifyGroup.POST("/carrier-service/:integration_id/enable", middleware.JWT(), h.EnableCarrierServiceHandler)
		shopifyGroup.POST("/carrier-service/:integration_id/disable", middleware.JWT(), h.DisableCarrierServiceHandler)
		shopifyGroup.POST("/auto-guide/:integration_id", middleware.JWT(), h.SetAutoGuideHandler)

		shopifyGroup.POST("/webhook", h.WebhookHandler)
		shopifyGroup.POST("/webhook/:integration_id", h.WebhookHandler)

		shopifyGroup.POST("/webhooks/compliance", h.ComplianceWebhookHandler)

		shopifyGroup.POST("/webhooks/customers/data_request", h.CustomerDataRequestHandler)
		shopifyGroup.POST("/webhooks/customers/redact", h.CustomerRedactHandler)
		shopifyGroup.POST("/webhooks/shop/redact", h.ShopRedactHandler)
	}

	router.GET("/shopify/callback", h.OAuthCallbackHandler)
}
