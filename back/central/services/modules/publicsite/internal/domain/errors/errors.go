package errors

import "errors"

var (
	ErrBusinessNotFound    = errors.New("negocio no encontrado")
	ErrProductNotFound     = errors.New("producto no encontrado")
	ErrInvalidContact      = errors.New("nombre y mensaje son requeridos")
	ErrPublicSiteNotActive = errors.New("el sitio web público no está activo para este negocio")
)
