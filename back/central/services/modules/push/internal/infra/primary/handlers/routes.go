package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	push := router.Group("/push")
	push.Use(middleware.JWT())
	{
		push.POST("/devices", h.RegisterDevice)
		push.GET("/devices", h.ListDevices)
		push.DELETE("/devices", h.UnregisterDevice)
	}
}
