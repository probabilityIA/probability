package notification_event_type

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	events := router.Group("/notification-event-types")
	events.Use(middleware.JWT())
	{
		events.GET("", h.List)
		events.GET("/:id", h.GetByID)

		events.POST("", middleware.RequireSuperAdmin(), h.Create)
		events.PUT("/:id", middleware.RequireSuperAdmin(), h.Update)
		events.PATCH("/:id/toggle-active", middleware.RequireSuperAdmin(), h.ToggleActive)
		events.DELETE("/:id", middleware.RequireSuperAdmin(), h.Delete)
	}
}
