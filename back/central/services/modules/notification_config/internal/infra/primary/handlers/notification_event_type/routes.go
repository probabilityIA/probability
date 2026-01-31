package notification_event_type

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas HTTP del módulo NotificationEventType
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	events := router.Group("/notification-event-types")
	{
		events.POST("", h.Create)
		events.GET("", h.List) // Soporta ?notification_type_id=X
		events.GET("/:id", h.GetByID)
		events.PUT("/:id", h.Update)
		events.DELETE("/:id", h.Delete)
	}
}
