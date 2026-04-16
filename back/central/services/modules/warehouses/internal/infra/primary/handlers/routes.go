package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas del módulo warehouses
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	warehouses := router.Group("/warehouses")
	warehouses.Use(middleware.JWT())
	{
		warehouses.GET("", h.ListWarehouses)
		warehouses.POST("", h.CreateWarehouse)
		warehouses.GET("/:id", h.GetWarehouse)
		warehouses.PUT("/:id", h.UpdateWarehouse)
		warehouses.DELETE("/:id", h.DeleteWarehouse)

		// Locations
		warehouses.GET("/:id/locations", h.ListLocations)
		warehouses.POST("/:id/locations", h.CreateLocation)
		warehouses.PUT("/:id/locations/:locationId", h.UpdateLocation)
		warehouses.DELETE("/:id/locations/:locationId", h.DeleteLocation)
	}
}
