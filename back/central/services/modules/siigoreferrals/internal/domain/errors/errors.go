package errors

import "errors"

var (
	ErrInvalidReferral = errors.New("nombre, correo, telefono y rango de ordenes son requeridos")
)
