package notification_config

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas HTTP del módulo
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	configs := router.Group("/notification-configs")
	configs.Use(middleware.JWT())
	{
		configs.POST("", h.Create)
		configs.GET("", h.List)
		configs.PUT("/sync", h.SyncByIntegration)
		configs.GET("/:id", h.GetByID)
		configs.PUT("/:id", h.Update)
		configs.DELETE("/:id", h.Delete)
	}
}
