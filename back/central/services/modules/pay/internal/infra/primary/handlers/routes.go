package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra las rutas del módulo de pagos
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	pay := router.Group("/pay")
	{
		pay.POST("/transactions", h.CreatePayment)
		pay.GET("/transactions", h.ListPayments)
		pay.GET("/transactions/:id", h.GetPayment)
	}
}
