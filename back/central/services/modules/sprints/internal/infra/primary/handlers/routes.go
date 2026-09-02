package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/sprints", middleware.JWT())
	{
		g.GET("", h.List)
		g.POST("", h.Create)
		g.GET(":id", h.Get)
		g.PUT(":id", h.Update)
		g.DELETE(":id", h.Delete)
		g.PATCH(":id/status", h.ChangeStatus)
	}
}
