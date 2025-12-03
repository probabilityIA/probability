package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IIntegrationHandler define la interfaz del handler de integraciones
type IIntegrationHandler interface {
	GetIntegrationsHandler(c *gin.Context)
	GetIntegrationByIDHandler(c *gin.Context)
	GetIntegrationByTypeHandler(c *gin.Context)
	CreateIntegrationHandler(c *gin.Context)
	UpdateIntegrationHandler(c *gin.Context)
	DeleteIntegrationHandler(c *gin.Context)
	TestIntegrationHandler(c *gin.Context)
	ActivateIntegrationHandler(c *gin.Context)
	DeactivateIntegrationHandler(c *gin.Context)
	SetAsDefaultHandler(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, handler IIntegrationHandler, logger log.ILogger)
}

// RegisterRoutes registra las rutas del handler de integraciones
func (h *IntegrationHandler) RegisterRoutes(router *gin.RouterGroup, handler IIntegrationHandler, logger log.ILogger) {
	integrationsGroup := router.Group("/integrations")
	{
		// CRUD básico
		integrationsGroup.GET("", middleware.JWT(), handler.GetIntegrationsHandler)
		integrationsGroup.GET("/:id", middleware.JWT(), handler.GetIntegrationByIDHandler)
		integrationsGroup.GET("/type/:type", middleware.JWT(), handler.GetIntegrationByTypeHandler)
		integrationsGroup.POST("", middleware.JWT(), handler.CreateIntegrationHandler)
		integrationsGroup.PUT("/:id", middleware.JWT(), handler.UpdateIntegrationHandler)
		integrationsGroup.DELETE("/:id", middleware.JWT(), handler.DeleteIntegrationHandler)

		// Acciones específicas
		integrationsGroup.POST("/:id/test", middleware.JWT(), handler.TestIntegrationHandler)
		integrationsGroup.PUT("/:id/activate", middleware.JWT(), handler.ActivateIntegrationHandler)
		integrationsGroup.PUT("/:id/deactivate", middleware.JWT(), handler.DeactivateIntegrationHandler)
		integrationsGroup.PUT("/:id/set-default", middleware.JWT(), handler.SetAsDefaultHandler)
	}
}
