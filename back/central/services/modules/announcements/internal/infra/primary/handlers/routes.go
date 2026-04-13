package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	announcements := router.Group("/announcements")
	announcements.Use(middleware.JWT())
	{
		announcements.GET("", h.List)
		announcements.POST("", h.Create)
		announcements.GET("/categories", h.ListCategories)
		announcements.GET("/active", h.GetActive)
		announcements.GET("/:id", h.Get)
		announcements.PUT("/:id", h.Update)
		announcements.DELETE("/:id", h.Delete)
		announcements.PATCH("/:id/status", h.ChangeStatus)
		announcements.POST("/:id/view", h.RegisterView)
		announcements.GET("/:id/stats", h.GetStats)
		announcements.POST("/:id/force-redisplay", h.ForceRedisplay)
		announcements.POST("/:id/image", h.UploadImage)
		announcements.DELETE("/:id/image/:imageId", h.DeleteImage)
	}
}
