package message_audit

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra las rutas HTTP del módulo de auditoría de mensajes
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	audit := router.Group("/notification-configs/message-audit")
	audit.Use(middleware.JWT())
	{
		audit.GET("", h.List)
		audit.GET("/stats", h.Stats)
		audit.GET("/conversations", h.ListConversations)
		audit.GET("/conversations/:id/messages", h.GetConversationMessages)
	}
}
