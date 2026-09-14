package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *SSEHandler) RegisterRoutes(router *gin.RouterGroup) {
	notifyGroup := router.Group("/notify")
	{
		notifyGroup.GET("/sse/order-notify/:businessID", middleware.JWT(), h.HandleSSE)
		notifyGroup.GET("/sse/order-notify", middleware.JWT(), h.HandleSSE)
	}
}
