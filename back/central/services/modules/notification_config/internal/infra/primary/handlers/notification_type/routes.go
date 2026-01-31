package notification_type

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas HTTP del módulo NotificationType
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	types := router.Group("/notification-types")
	{
		types.POST("", h.Create)
		types.GET("", h.List)
		types.GET("/:id", h.GetByID)
		types.PUT("/:id", h.Update)
		types.DELETE("/:id", h.Delete)
	}
}
