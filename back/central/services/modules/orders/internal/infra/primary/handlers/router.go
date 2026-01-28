package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas del módulo orders
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		// CRUD básico
		orders.GET("", middleware.JWT(), h.ListOrders)
		orders.GET("/:id", middleware.JWT(), h.GetOrderByID)
		orders.GET("/:id/raw", middleware.JWT(), h.GetOrderRaw)
		orders.POST("", middleware.JWT(), h.CreateOrder)
		orders.PUT("/:id", middleware.JWT(), h.UpdateOrder)
		orders.DELETE("/:id", middleware.JWT(), h.DeleteOrder)

		// Mapeo de órdenes canónicas (para integraciones)
		orders.POST("/map", middleware.JWT(), h.MapAndSaveOrder)

		// Confirmación de órdenes vía WhatsApp
		orders.POST("/:id/request-confirmation", middleware.JWT(), h.RequestConfirmation)
	}
}
