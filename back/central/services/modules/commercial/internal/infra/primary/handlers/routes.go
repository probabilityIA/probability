package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	commercial := router.Group("/commercial")
	commercial.Use(middleware.JWT(), middleware.RequireSuperAdmin())
	{
		commercial.GET("/prospects", h.ListProspects)
		commercial.GET("/prospects/stats", h.GetStats)
		commercial.PATCH("/prospects/:id/seen", h.UpdateSeen)
	}
}
