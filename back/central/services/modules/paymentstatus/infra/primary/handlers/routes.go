package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas del módulo de payment statuses
func (h *PaymentStatusHandlers) RegisterRoutes(router *gin.RouterGroup) {
	// Rutas para estados de pago
	statuses := router.Group("/payment-statuses")
	{
		statuses.GET("", h.ListPaymentStatuses) // GET /api/v1/payment-statuses
	}
}
