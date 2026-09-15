package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai")
	{
		group.GET("/recommendation", middleware.JWT(), h.GetRecommendation)

		assistant := group.Group("/assistant")
		assistant.POST("/chat", middleware.JWT(), h.AssistantChat)
		assistant.GET("/state", middleware.JWT(), h.GetAssistantState)
		assistant.POST("/intro-seen", middleware.JWT(), h.MarkAssistantIntroSeen)
	}
}
