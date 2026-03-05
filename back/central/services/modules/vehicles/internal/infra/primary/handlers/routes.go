package handlers

import "github.com/gin-gonic/gin"

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	vehicles := router.Group("/vehicles")
	{
		vehicles.GET("", h.ListVehicles)
		vehicles.GET("/:id", h.GetVehicle)
		vehicles.POST("", h.CreateVehicle)
		vehicles.PUT("/:id", h.UpdateVehicle)
		vehicles.DELETE("/:id", h.DeleteVehicle)
	}
}
