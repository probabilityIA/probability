package domain

import "errors"

var (
	// Errores de recursos
	ErrResourceNotFound   = errors.New("recurso no encontrado")
	ErrResourceNameExists = errors.New("ya existe un recurso con este nombre")
)
