package whatsapp_template

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/whatsapp-templates")
	group.Use(middleware.JWT())
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/variables", h.Variables)
		group.GET("/:id", h.GetByID)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
		group.POST("/:id/resubmit", h.Resubmit)
	}
}
