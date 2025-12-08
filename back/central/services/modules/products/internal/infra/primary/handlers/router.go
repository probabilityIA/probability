package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra todas las rutas del módulo products
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		// CRUD básico
		products.GET("", h.ListProducts)
		products.GET("/:id", h.GetProductByID)
		products.POST("", h.CreateProduct)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
	}
}

