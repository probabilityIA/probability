package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	warehouses := router.Group("/warehouses")
	warehouses.Use(middleware.JWT())
	{
		warehouses.GET("", h.ListWarehouses)
		warehouses.POST("", h.CreateWarehouse)
		warehouses.GET("/:id", h.GetWarehouse)
		warehouses.PUT("/:id", h.UpdateWarehouse)
		warehouses.DELETE("/:id", h.DeleteWarehouse)

		warehouses.GET("/:id/locations", h.ListLocations)
		warehouses.POST("/:id/locations", h.CreateLocation)
		warehouses.PUT("/:id/locations/:locationId", h.UpdateLocation)
		warehouses.DELETE("/:id/locations/:locationId", h.DeleteLocation)

		warehouses.GET("/:id/tree", h.GetWarehouseTree)
		warehouses.GET("/:id/zones", h.ListZones)
	}

	zones := router.Group("/zones")
	zones.Use(middleware.JWT())
	{
		zones.POST("", h.CreateZone)
		zones.GET("/:zoneId", h.GetZone)
		zones.PUT("/:zoneId", h.UpdateZone)
		zones.DELETE("/:zoneId", h.DeleteZone)
		zones.GET("/:zoneId/aisles", h.ListAisles)
	}

	aisles := router.Group("/aisles")
	aisles.Use(middleware.JWT())
	{
		aisles.POST("", h.CreateAisle)
		aisles.GET("/:aisleId", h.GetAisle)
		aisles.PUT("/:aisleId", h.UpdateAisle)
		aisles.DELETE("/:aisleId", h.DeleteAisle)
		aisles.GET("/:aisleId/racks", h.ListRacks)
	}

	racks := router.Group("/racks")
	racks.Use(middleware.JWT())
	{
		racks.POST("", h.CreateRack)
		racks.GET("/:rackId", h.GetRack)
		racks.PUT("/:rackId", h.UpdateRack)
		racks.DELETE("/:rackId", h.DeleteRack)
		racks.GET("/:rackId/levels", h.ListRackLevels)
	}

	rackLevels := router.Group("/rack-levels")
	rackLevels.Use(middleware.JWT())
	{
		rackLevels.POST("", h.CreateRackLevel)
		rackLevels.GET("/:levelId", h.GetRackLevel)
		rackLevels.PUT("/:levelId", h.UpdateRackLevel)
		rackLevels.DELETE("/:levelId", h.DeleteRackLevel)
	}
}
