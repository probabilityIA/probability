package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra todas las rutas del módulo orders
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		orders.GET("", h.ListOrders)
		orders.GET("/:id", h.GetOrderByID)
		orders.POST("", h.CreateOrder)
		orders.PUT("/:id", h.UpdateOrder)
		orders.DELETE("/:id", h.DeleteOrder)
	}
}
