package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	legal := router.Group("/legal")
	legal.Use(middleware.JWT())
	{
		legal.GET("/pending", h.GetPending)
		legal.POST("/accept", h.Accept)
	}
}
