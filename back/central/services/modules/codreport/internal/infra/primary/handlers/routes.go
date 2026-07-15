package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/cod-report", middleware.JWT())
	{
		g.GET("/summary", h.Summary)
		g.GET("/orders", h.ListOrders)
		g.GET("/cuts", h.ListCuts)
		g.GET("/cuts/selectable", h.SelectableCutOrders)
		g.POST("/cuts/confirm", h.ConfirmCut)
		g.GET("/carrier-config", h.CarrierConfigs)
		g.PUT("/carrier-config", h.SaveCarrierConfig)
	}
}
