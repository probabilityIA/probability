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
		assistant.POST("/messages/:id/feedback", middleware.JWT(), h.SubmitAssistantFeedback)
		assistant.POST("/messages/:id/click", middleware.JWT(), h.MarkAssistantClick)
		assistant.GET("/alerts", middleware.JWT(), h.ListAssistantAlerts)
		assistant.GET("/alerts/unread", middleware.JWT(), h.GetAssistantAlertsUnread)
		assistant.POST("/alerts/seen", middleware.JWT(), h.MarkAssistantAlertsSeen)

		admin := assistant.Group("/admin", middleware.JWT(), middleware.RequireSuperAdmin())
		admin.GET("/messages", h.ListAssistantReviewMessages)
		admin.GET("/summary", h.GetAssistantReviewSummary)
	}
}
