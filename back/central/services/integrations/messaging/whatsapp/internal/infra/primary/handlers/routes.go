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

		whatsapp.PUT("/connection", middleware.JWT(), h.SaveConnection)

		whatsapp.GET("/numbers", middleware.JWT(), h.GetNumberState)
		whatsapp.POST("/numbers", middleware.JWT(), h.AddNumber)
		whatsapp.POST("/numbers/code", middleware.JWT(), h.RequestNumberCode)
		whatsapp.POST("/numbers/verify", middleware.JWT(), h.VerifyNumberCode)
		whatsapp.POST("/numbers/register", middleware.JWT(), h.RegisterNumber)

		whatsapp.GET("/embedded-signup/config", middleware.JWT(), h.GetEmbeddedSignupConfig)
		whatsapp.POST("/embedded-signup", middleware.JWT(), h.CompleteEmbeddedSignup)
		whatsapp.GET("/templates/status", middleware.JWT(), h.GetTemplatesStatus)
		whatsapp.GET("/templates/preview", middleware.JWT(), h.PreviewTemplates)
		whatsapp.POST("/templates/provision", middleware.JWT(), h.ProvisionTemplates)

		whatsapp.GET("/webhook", h.VerifyWebhook)
		whatsapp.POST("/webhook", h.ReceiveWebhook)
	}
}
