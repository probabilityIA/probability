package errors

import "errors"

var (
	ErrOrderNotFound   = errors.New("orden no encontrada")
	ErrOrderNotInScope = errors.New("la orden no pertenece al negocio")
)
