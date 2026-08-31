package errors

import "errors"

var (
	ErrUserRequired    = errors.New("usuario no identificado")
	ErrBusinessRequired = errors.New("business_id es requerido para super admin")
	ErrTourKeyRequired = errors.New("tour_key es requerido")
	ErrInvalidStatus   = errors.New("estado de tour invalido")
	ErrInvalidVersion  = errors.New("la version del tour debe ser mayor a cero")
)
