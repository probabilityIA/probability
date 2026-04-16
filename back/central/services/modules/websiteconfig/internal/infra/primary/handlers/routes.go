package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	wc := router.Group("/website-config", middleware.JWT())
	{
		wc.GET("", h.GetConfig)
		wc.PUT("", h.UpdateConfig)
	}
}
