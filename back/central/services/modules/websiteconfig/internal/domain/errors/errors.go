package errors

import "errors"

var (
	ErrConfigNotFound   = errors.New("configuracion de sitio web no encontrada")
	ErrBusinessRequired = errors.New("business_id es requerido")
)
