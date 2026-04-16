package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	vehicles := router.Group("/vehicles")
	vehicles.Use(middleware.JWT())
	{
		vehicles.GET("", h.ListVehicles)
		vehicles.GET("/:id", h.GetVehicle)
		vehicles.POST("", h.CreateVehicle)
		vehicles.PUT("/:id", h.UpdateVehicle)
		vehicles.DELETE("/:id", h.DeleteVehicle)
	}
}
