package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas HTTP del módulo WhatsApp
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	whatsapp := router.Group("/whatsapp")
	{
		// Endpoints protegidos con JWT
		whatsapp.POST("/send-template", middleware.JWT(), h.SendTemplate)
		whatsapp.POST("/conversations/:id/reply", middleware.JWT(), h.SendManualReply)
		whatsapp.POST("/conversations/:id/pause-ai", middleware.JWT(), h.PauseAI)
		whatsapp.POST("/conversations/:id/resume-ai", middleware.JWT(), h.ResumeAI)

		// Webhook endpoints - SIN JWT (Meta usa su propia autenticación: verify_token + X-Hub-Signature-256)
		whatsapp.GET("/webhook", h.VerifyWebhook)
		whatsapp.POST("/webhook", h.ReceiveWebhook)
	}
}
