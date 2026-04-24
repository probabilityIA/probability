package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/app"
)

type IHandlers interface {
	ListWarehouses(c *gin.Context)
	GetWarehouse(c *gin.Context)
	CreateWarehouse(c *gin.Context)
	UpdateWarehouse(c *gin.Context)
	DeleteWarehouse(c *gin.Context)
	ListLocations(c *gin.Context)
	CreateLocation(c *gin.Context)
	UpdateLocation(c *gin.Context)
	DeleteLocation(c *gin.Context)

	GetWarehouseTree(c *gin.Context)
	CreateZone(c *gin.Context)
	GetZone(c *gin.Context)
	ListZones(c *gin.Context)
	UpdateZone(c *gin.Context)
	DeleteZone(c *gin.Context)
	CreateAisle(c *gin.Context)
	GetAisle(c *gin.Context)
	ListAisles(c *gin.Context)
	UpdateAisle(c *gin.Context)
	DeleteAisle(c *gin.Context)
	CreateRack(c *gin.Context)
	GetRack(c *gin.Context)
	ListRacks(c *gin.Context)
	UpdateRack(c *gin.Context)
	DeleteRack(c *gin.Context)
	CreateRackLevel(c *gin.Context)
	GetRackLevel(c *gin.Context)
	ListRackLevels(c *gin.Context)
	UpdateRackLevel(c *gin.Context)
	DeleteRackLevel(c *gin.Context)

	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc app.IUseCase
}

func New(uc app.IUseCase) IHandlers {
	return &Handlers{uc: uc}
}

func (h *Handlers) resolveBusinessID(c *gin.Context) (uint, bool) {
	businessID := c.GetUint("business_id")
	if businessID > 0 {
		return businessID, true
	}
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			return uint(id), true
		}
	}
	return 0, false
}
