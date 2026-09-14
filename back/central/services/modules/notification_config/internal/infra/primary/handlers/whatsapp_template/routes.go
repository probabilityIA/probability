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
		group.POST("/media", h.UploadMedia)
		group.POST("/sync-status", h.SyncStatuses)
		group.GET("/flows", h.ListAllFlows)
		group.GET("/:id/flows", h.GetFlows)
		group.PUT("/:id/flows", h.ReplaceFlows)
		group.GET("/:id", h.GetByID)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
		group.POST("/:id/submit", h.SubmitForReview)
	}

	flows := router.Group("/whatsapp-flows")
	flows.Use(middleware.JWT())
	{
		flows.GET("", h.ListFlowGroups)
		flows.POST("", h.CreateFlowGroup)
		flows.PUT("/:id", h.UpdateFlowGroup)
		flows.DELETE("/:id", h.DeleteFlowGroup)
		flows.GET("/:id/transitions", h.ListFlowTransitions)
	}
}
