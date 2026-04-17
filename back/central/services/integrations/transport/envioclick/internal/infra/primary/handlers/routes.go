package handlers

import "github.com/gin-gonic/gin"

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/integrations/transport/envioclick")
	group.POST("/webhook", h.Webhook)

	router.POST("/webhooks/envioclick", h.Webhook)
}
