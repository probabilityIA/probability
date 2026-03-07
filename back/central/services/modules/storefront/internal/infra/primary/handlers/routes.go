package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registers all storefront routes
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	sf := router.Group("/storefront")
	{
		// Public (no auth required)
		sf.POST("/register", h.Register)

		// Protected (JWT required, role validation in handlers)
		sf.GET("/catalog", middleware.JWT(), h.ListCatalog)
		sf.GET("/catalog/:id", middleware.JWT(), h.GetProduct)
		sf.POST("/orders", middleware.JWT(), h.CreateOrder)
		sf.GET("/orders", middleware.JWT(), h.ListMyOrders)
		sf.GET("/orders/:id", middleware.JWT(), h.GetMyOrder)
	}
}
