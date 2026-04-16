package notification_type

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas HTTP del módulo NotificationType
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	types := router.Group("/notification-types")
	types.Use(middleware.JWT())
	{
		types.POST("", h.Create)
		types.GET("", h.List)
		types.GET("/:id", h.GetByID)
		types.PUT("/:id", h.Update)
		types.DELETE("/:id", h.Delete)
	}
}
