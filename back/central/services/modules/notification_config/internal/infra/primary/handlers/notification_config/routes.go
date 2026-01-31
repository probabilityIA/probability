package notification_config

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas HTTP del módulo
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	configs := router.Group("/notification-configs")
	{
		configs.POST("", h.Create)
		configs.GET("", h.List)
		configs.GET("/:id", h.GetByID)
		configs.PUT("/:id", h.Update)
		configs.DELETE("/:id", h.Delete)
	}
}
