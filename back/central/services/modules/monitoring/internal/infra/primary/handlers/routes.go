package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra las rutas HTTP del módulo de monitoreo (sin JWT - autenticado por HMAC)
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	monitoring := router.Group("/monitoring")
	monitoring.POST("/alerts/grafana", h.WebhookGrafana)
}
