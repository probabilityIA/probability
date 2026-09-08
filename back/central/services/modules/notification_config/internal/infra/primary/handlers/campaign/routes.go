package campaign

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/whatsapp-campaigns")
	group.Use(middleware.JWT())
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.POST("/audience-preview", h.PreviewAudience)
		group.GET("/:id", h.GetByID)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
		group.POST("/:id/launch", h.Launch)
		group.POST("/:id/pause", h.Pause)
		group.POST("/:id/resume", h.Resume)
		group.POST("/:id/cancel", h.Cancel)
		group.GET("/:id/sends", h.ListSends)
	}
}
