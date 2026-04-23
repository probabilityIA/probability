package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/notification-backfill")
	g.Use(middleware.JWT())
	{
		g.GET("/events", h.ListEvents)
		g.POST("/preview", h.Preview)
		g.POST("/run", h.Run)
		g.GET("/jobs/:job_id", h.GetJob)
	}
}
