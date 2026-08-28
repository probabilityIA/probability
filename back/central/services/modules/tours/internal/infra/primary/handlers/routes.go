package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	tours := router.Group("/tours")
	tours.Use(middleware.JWT())
	{
		tours.GET("/progress", h.ListProgress)
		tours.PUT("/progress", h.SaveProgress)
		tours.DELETE("/progress/:tour_key", h.ResetTour)
		tours.POST("/progress/reset", h.ResetAll)
		tours.POST("/progress/skip-all", h.SkipAll)
	}
}
