package handlerintegrations

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IIntegrationHandler define la interfaz del handler de integraciones

// RegisterRoutes registra las rutas del handler de integraciones
func (h *IntegrationHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	integrationsGroup := router.Group("/integrations")
	{
		// CRUD básico
		integrationsGroup.GET("", middleware.JWT(), h.GetIntegrationsHandler)
		integrationsGroup.GET("/simple", middleware.JWT(), h.GetIntegrationsSimpleHandler)
		integrationsGroup.GET("/:id", middleware.JWT(), h.GetIntegrationByIDHandler) // Devuelve credenciales solo si es super admin
		integrationsGroup.GET("/type/:type", middleware.JWT(), h.GetIntegrationByTypeHandler)
		integrationsGroup.POST("", middleware.JWT(), h.CreateIntegrationHandler)
		integrationsGroup.PUT("/:id", middleware.JWT(), h.UpdateIntegrationHandler)
		integrationsGroup.DELETE("/:id", middleware.JWT(), h.DeleteIntegrationHandler)

		// Acciones específicas
		integrationsGroup.POST("/test", middleware.JWT(), h.TestConnectionRawHandler)
		integrationsGroup.POST("/sync-orders/business/:business_id", middleware.JWT(), h.SyncOrdersByBusinessHandler)
		integrationsGroup.POST("/:id/test", middleware.JWT(), h.TestIntegrationHandler)
		integrationsGroup.POST("/:id/sync", middleware.JWT(), h.SyncOrdersByIntegrationIDHandler)
		integrationsGroup.PUT("/:id/activate", middleware.JWT(), h.ActivateIntegrationHandler)
		integrationsGroup.PUT("/:id/deactivate", middleware.JWT(), h.DeactivateIntegrationHandler)
		integrationsGroup.PUT("/:id/set-default", middleware.JWT(), h.SetAsDefaultHandler)
		integrationsGroup.GET("/:id/webhook", middleware.JWT(), h.GetWebhookURLHandler)
		integrationsGroup.GET("/:id/webhooks", middleware.JWT(), h.ListWebhooksHandler)
		integrationsGroup.GET("/:id/webhooks/verify", middleware.JWT(), h.VerifyWebhooksHandler)
		integrationsGroup.POST("/:id/webhooks/create", middleware.JWT(), h.CreateWebhookHandler)
		integrationsGroup.DELETE("/:id/webhooks/:webhook_id", middleware.JWT(), h.DeleteWebhookHandler)
	}
}
