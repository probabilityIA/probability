package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// IHandlers define la interfaz de handlers del módulo inventory
type IHandlers interface {
	GetProductInventory(c *gin.Context)
	ListWarehouseInventory(c *gin.Context)
	AdjustStock(c *gin.Context)
	TransferStock(c *gin.Context)
	BulkLoadInventory(c *gin.Context)
	ListMovements(c *gin.Context)
	ListMovementTypes(c *gin.Context)
	CreateMovementType(c *gin.Context)
	UpdateMovementType(c *gin.Context)
	DeleteMovementType(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

// handlers contiene el use case
type handlers struct {
	uc     app.IUseCase
	rabbit rabbitmq.IQueue
}

// New crea una nueva instancia de los handlers
func New(uc app.IUseCase, rabbit rabbitmq.IQueue) IHandlers {
	return &handlers{uc: uc, rabbit: rabbit}
}

// resolveBusinessID obtiene el business_id efectivo.
func (h *handlers) resolveBusinessID(c *gin.Context) (uint, bool) {
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
