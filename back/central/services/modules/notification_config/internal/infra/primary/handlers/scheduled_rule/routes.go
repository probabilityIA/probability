package scheduled_rule

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/scheduled-notification-rules")
	group.Use(middleware.JWT())
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.GetByID)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
		group.POST("/:id/run", h.RunNow)
		group.GET("/:id/runs", h.ListRuns)
	}
}
