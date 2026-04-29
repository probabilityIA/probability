package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/shipping-margins")
	g.Use(middleware.JWT())
	{
		g.GET("", h.List)
		g.GET("/:id", h.Get)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
	}
}
