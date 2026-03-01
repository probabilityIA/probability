package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra las rutas SSE
func (h *SSEHandler) RegisterRoutes(router *gin.RouterGroup) {
	notifyGroup := router.Group("/notify")
	{
		notifyGroup.GET("/sse/order-notify/:businessID", h.HandleSSE)
		notifyGroup.GET("/sse/order-notify", h.HandleSSE)
	}
}
