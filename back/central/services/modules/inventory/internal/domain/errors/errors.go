package errors

import "errors"

var (
	ErrProductNotFound        = errors.New("producto no encontrado")
	ErrWarehouseNotFound      = errors.New("bodega no encontrada")
	ErrInsufficientStock      = errors.New("stock insuficiente para esta operación")
	ErrInvalidQuantity        = errors.New("la cantidad no puede ser cero")
	ErrSameWarehouse          = errors.New("la bodega de origen y destino deben ser diferentes")
	ErrTransferQtyNeg         = errors.New("la cantidad a transferir debe ser positiva")
	ErrProductNoTracking      = errors.New("el producto no tiene seguimiento de inventario habilitado")
	ErrMovementTypeNotFound   = errors.New("tipo de movimiento no encontrado")
	ErrMovementTypeCodeExists = errors.New("el código del tipo de movimiento ya existe")
	ErrNoDefaultWarehouse    = errors.New("no se encontró bodega por defecto para el negocio")
)
