package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas HTTP del módulo WhatsApp
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	whatsapp := router.Group("/whatsapp")
	{
		// Endpoint protegido con JWT
		whatsapp.POST("/send-template", middleware.JWT(), h.SendTemplate)

		// Webhook endpoints - SIN JWT (Meta usa su propia autenticación: verify_token + X-Hub-Signature-256)
		whatsapp.GET("/webhook", h.VerifyWebhook)
		whatsapp.POST("/webhook", h.ReceiveWebhook)
	}
}
