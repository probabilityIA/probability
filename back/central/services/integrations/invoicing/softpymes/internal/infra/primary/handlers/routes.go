package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra las rutas HTTP del módulo Softpymes
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	// Grupo de proveedores
	providers := router.Group("/invoicing/providers")
	{
		providers.POST("", h.CreateProvider)
		providers.GET("", h.ListProviders)
		providers.GET("/:id", h.GetProvider)
		providers.PATCH("/:id", h.UpdateProvider)
		providers.DELETE("/:id", h.DeleteProvider)
		providers.POST("/:id/test", h.TestProvider)
	}

	// Tipos de proveedores
	router.GET("/invoicing/provider-types", h.ListProviderTypes)
}
