package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/subscriptions")
	g.Use(middleware.JWT())
	{
		g.GET("/me", h.GetCurrentSubscription)
		g.POST("/purchase", h.PurchaseSubscription)
		g.GET("/module-codes", h.GetModuleCodes)

		g.GET("/types", h.ListSubscriptionTypes)
		g.GET("/types/:id", h.GetSubscriptionType)
		g.POST("/types", h.CreateSubscriptionType)
		g.PUT("/types/:id", h.UpdateSubscriptionType)
		g.DELETE("/types/:id", h.DeleteSubscriptionType)

		g.POST("/register-payment", h.RegisterPayment)
		g.POST("/disable", h.DisableSubscription)

		g.GET("/overrides/:businessId", h.ListOverrides)
		g.POST("/overrides", h.GrantOverride)
		g.DELETE("/overrides/:businessId/:moduleCode", h.RevokeOverride)
	}
}
