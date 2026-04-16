package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas del módulo pricing
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	pricing := router.Group("/pricing")
	{
		// Client Pricing Rules
		rules := pricing.Group("/rules")
		rules.GET("", middleware.JWT(), h.ListClientPricingRules)
		rules.GET("/:id", middleware.JWT(), h.GetClientPricingRule)
		rules.POST("", middleware.JWT(), h.CreateClientPricingRule)
		rules.PUT("/:id", middleware.JWT(), h.UpdateClientPricingRule)
		rules.DELETE("/:id", middleware.JWT(), h.DeleteClientPricingRule)

		// Quantity Discounts
		discounts := pricing.Group("/quantity-discounts")
		discounts.GET("", middleware.JWT(), h.ListQuantityDiscounts)
		discounts.GET("/:id", middleware.JWT(), h.GetQuantityDiscount)
		discounts.POST("", middleware.JWT(), h.CreateQuantityDiscount)
		discounts.PUT("/:id", middleware.JWT(), h.UpdateQuantityDiscount)
		discounts.DELETE("/:id", middleware.JWT(), h.DeleteQuantityDiscount)

		// Price Calculator (preview)
		pricing.POST("/calculate", middleware.JWT(), h.CalculatePrice)
	}
}
