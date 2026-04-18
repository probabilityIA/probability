package errors

import "errors"

var (
	ErrProductNotFound        = errors.New("producto no encontrado")
	ErrWarehouseNotFound      = errors.New("bodega no encontrada")
	ErrInsufficientStock      = errors.New("stock insuficiente para esta operacion")
	ErrInvalidQuantity        = errors.New("la cantidad no puede ser cero")
	ErrSameWarehouse          = errors.New("la bodega de origen y destino deben ser diferentes")
	ErrTransferQtyNeg         = errors.New("la cantidad a transferir debe ser positiva")
	ErrProductNoTracking      = errors.New("el producto no tiene seguimiento de inventario habilitado")
	ErrMovementTypeNotFound   = errors.New("tipo de movimiento no encontrado")
	ErrMovementTypeCodeExists = errors.New("el codigo del tipo de movimiento ya existe")
	ErrNoDefaultWarehouse     = errors.New("no se encontro bodega por defecto para el negocio")

	ErrLotNotFound         = errors.New("lote no encontrado")
	ErrDuplicateLotCode    = errors.New("el codigo de lote ya existe para este producto")
	ErrLotExpired          = errors.New("el lote esta vencido")
	ErrSerialNotFound      = errors.New("numero de serie no encontrado")
	ErrSerialNotAvailable  = errors.New("numero de serie no disponible")
	ErrDuplicateSerial     = errors.New("el numero de serie ya existe para este producto")
	ErrUomNotFound         = errors.New("unidad de medida no encontrada")
	ErrUomConversion       = errors.New("error en conversion de unidad de medida")
	ErrStateNotFound       = errors.New("estado de inventario no encontrado")
	ErrStateTransition     = errors.New("transicion de estado no permitida")
	ErrProductUoMNotFound  = errors.New("unidad de medida del producto no encontrada")
)
