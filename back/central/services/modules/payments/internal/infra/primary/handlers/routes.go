package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas del módulo de payments
func (h *PaymentHandlers) RegisterRoutes(router *gin.RouterGroup) {
	// Payment Statuses (catálogo)
	paymentStatuses := router.Group("/payment-statuses")
	paymentStatuses.Use(middleware.JWT())
	{
		paymentStatuses.GET("", h.ListPaymentStatuses) // GET /api/v1/payment-statuses
	}

	payments := router.Group("/payments")
	payments.Use(middleware.JWT())
	{
		// Channel Payment Methods routes
		channelMethods := payments.Group("/channel-methods")
		{
			channelMethods.GET("", h.ListChannelPaymentMethods) // GET /api/v1/payments/channel-methods
		}

		// Payment Methods routes
		methods := payments.Group("/methods")
		{
			methods.GET("", h.ListPaymentMethods)               // GET /api/v1/payments/methods
			methods.GET("/:id", h.GetPaymentMethod)             // GET /api/v1/payments/methods/:id
			methods.POST("", h.CreatePaymentMethod)             // POST /api/v1/payments/methods
			methods.PUT("/:id", h.UpdatePaymentMethod)          // PUT /api/v1/payments/methods/:id
			methods.DELETE("/:id", h.DeletePaymentMethod)       // DELETE /api/v1/payments/methods/:id
			methods.PATCH("/:id/toggle", h.TogglePaymentMethod) // PATCH /api/v1/payments/methods/:id/toggle
		}

		// Payment Mappings routes
		mappings := payments.Group("/mappings")
		{
			mappings.GET("", h.ListPaymentMappings)                               // GET /api/v1/payments/mappings
			mappings.GET("/:id", h.GetPaymentMapping)                             // GET /api/v1/payments/mappings/:id
			mappings.GET("/integration/:type", h.GetPaymentMappingsByIntegration) // GET /api/v1/payments/mappings/integration/:type
			mappings.POST("", h.CreatePaymentMapping)                             // POST /api/v1/payments/mappings
			mappings.PUT("/:id", h.UpdatePaymentMapping)                          // PUT /api/v1/payments/mappings/:id
			mappings.DELETE("/:id", h.DeletePaymentMapping)                       // DELETE /api/v1/payments/mappings/:id
			mappings.PATCH("/:id/toggle", h.TogglePaymentMapping)                 // PATCH /api/v1/payments/mappings/:id/toggle
		}
	}
}
