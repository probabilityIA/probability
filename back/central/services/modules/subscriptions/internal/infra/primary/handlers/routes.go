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
		g.GET("/me/usage", h.GetMySubscriptionUsage)
		g.POST("/overage/accept", h.AcceptShipmentOverage)
		g.POST("/overage/pay", h.PayShipmentOverage)
		g.POST("/purchase", h.PurchaseSubscription)
		g.GET("/module-codes", h.GetModuleCodes)
		g.GET("/module-catalog", h.GetModuleCatalog)
		g.GET("/my-modules", h.GetMyModules)

		g.GET("/types", h.ListSubscriptionTypes)
		g.GET("/types/:id", h.GetSubscriptionType)
		g.POST("/types", h.CreateSubscriptionType)
		g.PUT("/types/:id", h.UpdateSubscriptionType)
		g.DELETE("/types/:id", h.DeleteSubscriptionType)

		g.GET("/custom-plans", h.ListCustomPlans)
		g.POST("/custom-plans", h.CreateCustomPlan)
		g.PUT("/custom-plans/:id", h.UpdateCustomPlan)
		g.DELETE("/custom-plans/:id", h.DeleteCustomPlan)

		g.POST("/register-payment", h.RegisterPayment)
		g.PUT("/edit-dates", h.EditSubscriptionDates)
		g.POST("/disable", h.DisableSubscription)
		g.POST("/reactivate", h.ReactivateSubscription)
		g.POST("/extend-days", h.ExtendCourtesy)
		g.POST("/payments/:id/revert", h.RevertPayment)
		g.GET("/payments/:businessId", h.ListPaymentHistory)
		g.GET("/audit-logs/:businessId", h.ListAuditLogs)

		g.GET("/overrides/:businessId", h.ListOverrides)
		g.POST("/overrides", h.GrantOverride)
		g.DELETE("/overrides/:businessId/:moduleCode", h.RevokeOverride)

		g.GET("/admin/businesses", h.ListAdminBusinesses)
		g.GET("/admin/kpis", h.GetAdminKPIs)
	}
}
