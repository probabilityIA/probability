package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas HTTP del módulo WhatsApp
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	whatsapp := router.Group("/whatsapp")
	{
		// Endpoint para envío de plantillas
		whatsapp.POST("/send-template", h.SendTemplate)

		// Endpoints de webhook
		whatsapp.GET("/webhook", h.VerifyWebhook)
		whatsapp.POST("/webhook", h.ReceiveWebhook)
	}
}
