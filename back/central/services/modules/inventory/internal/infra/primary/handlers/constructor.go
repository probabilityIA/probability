package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// friendlyValidationError traduce errores de validación de Gin a mensajes legibles
func friendlyValidationError(err error) string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return err.Error()
	}

	fieldNames := map[string]string{
		"ProductID":       "producto",
		"WarehouseID":     "bodega",
		"FromWarehouseID": "bodega origen",
		"ToWarehouseID":   "bodega destino",
		"Quantity":        "cantidad",
		"Reason":          "razon",
		"Notes":           "notas",
		"SKU":             "SKU",
		"Items":           "items",
	}

	var msgs []string
	for _, fe := range verrs {
		name := fe.Field()
		if friendly, ok := fieldNames[name]; ok {
			name = friendly
		}
		switch fe.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("El campo %s es obligatorio", name))
		case "min":
			msgs = append(msgs, fmt.Sprintf("El campo %s no cumple el minimo requerido", name))
		case "max":
			msgs = append(msgs, fmt.Sprintf("El campo %s excede el maximo permitido", name))
		default:
			msgs = append(msgs, fmt.Sprintf("El campo %s no es valido", name))
		}
	}
	return strings.Join(msgs, ". ")
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
