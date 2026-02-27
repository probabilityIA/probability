package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app"
)

// IHandlers define la interfaz de handlers del módulo inventory
type IHandlers interface {
	GetProductInventory(c *gin.Context)
	ListWarehouseInventory(c *gin.Context)
	AdjustStock(c *gin.Context)
	TransferStock(c *gin.Context)
	ListMovements(c *gin.Context)
	ListMovementTypes(c *gin.Context)
	CreateMovementType(c *gin.Context)
	UpdateMovementType(c *gin.Context)
	DeleteMovementType(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

// Handlers contiene el use case
type Handlers struct {
	uc app.IUseCase
}

// New crea una nueva instancia de los handlers
func New(uc app.IUseCase) IHandlers {
	return &Handlers{uc: uc}
}

// resolveBusinessID obtiene el business_id efectivo.
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
