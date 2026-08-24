package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	mobile := router.Group("/mobile")
	mobile.Use(middleware.JWT())
	{
		mobile.GET("/orders/:id/full", h.GetOrderFull)
	}
}
