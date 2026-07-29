package handlers

import "github.com/gin-gonic/gin"

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/integrations/transport/shipit")
	group.POST("/webhook", h.Webhook)
	group.PUT("/webhook", h.Webhook)

	router.POST("/webhooks/shipit", h.Webhook)
	router.PUT("/webhooks/shipit", h.Webhook)
}
