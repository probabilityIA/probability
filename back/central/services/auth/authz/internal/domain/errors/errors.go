package errors

import "errors"

var (
	ErrNoBusinessRelation = errors.New("el usuario no pertenece al negocio")
	ErrBusinessRequired   = errors.New("business_id es requerido para super admin")
)
