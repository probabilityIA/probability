package message_audit

import "github.com/gin-gonic/gin"

// RegisterRoutes registra las rutas HTTP del módulo de auditoría de mensajes
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	audit := router.Group("/notification-configs/message-audit")
	{
		audit.GET("", h.List)
		audit.GET("/stats", h.Stats)
	}
}
