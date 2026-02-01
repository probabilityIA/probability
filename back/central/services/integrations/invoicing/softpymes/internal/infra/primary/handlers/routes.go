package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra las rutas HTTP del módulo Softpymes
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	// Grupo de proveedores
	providers := router.Group("/invoicing/providers")
	{
		providers.POST("", middleware.JWT(), h.CreateProvider)
		providers.GET("", middleware.JWT(), h.ListProviders)
		providers.GET("/:id", middleware.JWT(), h.GetProvider)
		providers.PATCH("/:id", middleware.JWT(), h.UpdateProvider)
		providers.DELETE("/:id", middleware.JWT(), h.DeleteProvider)
		providers.POST("/:id/test", middleware.JWT(), h.TestProvider)
	}

	// Tipos de proveedores
	router.GET("/invoicing/provider-types", h.ListProviderTypes)
}
