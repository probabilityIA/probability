package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/app"
)

type IHandlers interface {
	ListVehicles(c *gin.Context)
	GetVehicle(c *gin.Context)
	CreateVehicle(c *gin.Context)
	UpdateVehicle(c *gin.Context)
	DeleteVehicle(c *gin.Context)
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
