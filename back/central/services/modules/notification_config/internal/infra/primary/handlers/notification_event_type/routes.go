package notification_event_type

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas HTTP del módulo NotificationEventType
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	events := router.Group("/notification-event-types")
	events.Use(middleware.JWT())
	{
		events.POST("", h.Create)
		events.GET("", h.List) // Soporta ?notification_type_id=X
		events.GET("/:id", h.GetByID)
		events.PUT("/:id", h.Update)
		events.PATCH("/:id/toggle-active", h.ToggleActive)
		events.DELETE("/:id", h.Delete)
	}
}
