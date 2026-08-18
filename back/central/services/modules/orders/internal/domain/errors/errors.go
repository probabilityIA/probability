package errors

import "errors"

var (
	ErrOrderAlreadyExists = errors.New("order with this external_id already exists for this integration")
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidStatus = errors.New("invalid status code")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrOrderInTerminalState = errors.New("order is in a terminal state and cannot be changed")
	ErrInsufficientStock = errors.New("inventario insuficiente para mover la orden a Seleccionando productos")
	ErrOrderBusinessDeleted = errors.New("business is deleted or does not exist")
)
