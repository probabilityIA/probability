package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	drivers := router.Group("/drivers")
	drivers.Use(middleware.JWT())
	{
		drivers.GET("", h.ListDrivers)
		drivers.GET("/:id", h.GetDriver)
		drivers.POST("", h.CreateDriver)
		drivers.PUT("/:id", h.UpdateDriver)
		drivers.DELETE("/:id", h.DeleteDriver)
	}
}
