package errors

import "errors"

var (
	ErrInvalidLead = errors.New("nombre, correo y telefono son requeridos")
)
