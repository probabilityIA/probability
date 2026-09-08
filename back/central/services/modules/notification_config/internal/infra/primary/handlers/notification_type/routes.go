package notification_type

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	types := router.Group("/notification-types")
	types.Use(middleware.JWT())
	{
		types.GET("", h.List)
		types.GET("/:id", h.GetByID)

		types.POST("", middleware.RequireSuperAdmin(), h.Create)
		types.PUT("/:id", middleware.RequireSuperAdmin(), h.Update)
		types.DELETE("/:id", middleware.RequireSuperAdmin(), h.Delete)
	}
}
