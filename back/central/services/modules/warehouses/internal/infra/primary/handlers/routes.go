package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas del módulo warehouses
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	warehouses := router.Group("/warehouses")
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
