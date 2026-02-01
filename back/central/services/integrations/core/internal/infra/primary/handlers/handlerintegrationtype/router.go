package handlerintegrationtype

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

// RegisterRoutes registra las rutas del handler de tipos de integración
func (h *IntegrationTypeHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	integrationTypesGroup := router.Group("/integration-types")
	{
		// CRUD básico
		integrationTypesGroup.GET("", middleware.JWT(), h.ListIntegrationTypesHandler)
		integrationTypesGroup.GET("/active", middleware.JWT(), h.ListActiveIntegrationTypesHandler)
		integrationTypesGroup.GET("/:id", middleware.JWT(), h.GetIntegrationTypeByIDHandler)
		integrationTypesGroup.GET("/code/:code", middleware.JWT(), h.GetIntegrationTypeByCodeHandler)
		integrationTypesGroup.POST("", middleware.JWT(), h.CreateIntegrationTypeHandler)
		integrationTypesGroup.PUT("/:id", middleware.JWT(), h.UpdateIntegrationTypeHandler)
		integrationTypesGroup.DELETE("/:id", middleware.JWT(), h.DeleteIntegrationTypeHandler)
	}

	// Integration Categories
	integrationCategoriesGroup := router.Group("/integration-categories")
	{
		integrationCategoriesGroup.GET("", middleware.JWT(), h.ListIntegrationCategoriesHandler)
	}
}
