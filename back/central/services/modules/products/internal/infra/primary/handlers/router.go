package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas del módulo products
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		// CRUD básico
		products.GET("", middleware.JWT(), h.ListProducts)
		products.GET("/:id", middleware.JWT(), h.GetProductByID)
		products.POST("", middleware.JWT(), h.CreateProduct)
		products.PUT("/:id", middleware.JWT(), h.UpdateProduct)
		products.DELETE("/:id", middleware.JWT(), h.DeleteProduct)

		// Upload de imagen
		products.POST("/:id/image", middleware.JWT(), h.UploadProductImage)

		// Gestión de integraciones
		products.POST("/:id/integrations", middleware.JWT(), h.AddProductIntegration)
		products.GET("/:id/integrations", middleware.JWT(), h.GetProductIntegrations)
		products.DELETE("/:id/integrations/:integration_id", middleware.JWT(), h.RemoveProductIntegration)
	}
}
