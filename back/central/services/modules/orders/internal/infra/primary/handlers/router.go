package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra todas las rutas del módulo orders
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		// CRUD básico
		orders.GET("", h.ListOrders)
		orders.GET("/:id", h.GetOrderByID)
		orders.GET("/:id/raw", h.GetOrderRaw)
		orders.POST("", h.CreateOrder)
		orders.PUT("/:id", h.UpdateOrder)
		orders.DELETE("/:id", h.DeleteOrder)

		// Mapeo de órdenes canónicas (para integraciones)
		orders.POST("/map", h.MapAndSaveOrder)
	}
}
