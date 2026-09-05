package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	whatsapp := router.Group("/whatsapp")
	{

		whatsapp.POST("/send-template", middleware.JWT(), h.SendTemplate)
		whatsapp.POST("/conversations/:id/reply", middleware.JWT(), h.SendManualReply)
		whatsapp.POST("/conversations/:id/pause-ai", middleware.JWT(), h.PauseAI)
		whatsapp.POST("/conversations/:id/resume-ai", middleware.JWT(), h.ResumeAI)

		whatsapp.GET("/templates/status", middleware.JWT(), h.GetTemplatesStatus)
		whatsapp.POST("/templates/provision", middleware.JWT(), h.ProvisionTemplates)

		whatsapp.GET("/webhook", h.VerifyWebhook)
		whatsapp.POST("/webhook", h.ReceiveWebhook)
	}
}
