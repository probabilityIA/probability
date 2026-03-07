package errors

import "errors"

var (
	ErrProductNotFound     = errors.New("producto no encontrado")
	ErrOrderNotFound       = errors.New("orden no encontrada")
	ErrClientNotFound      = errors.New("cliente no encontrado")
	ErrBusinessNotFound    = errors.New("negocio no encontrado")
	ErrEmailAlreadyExists  = errors.New("ya existe un usuario con este email")
	ErrRoleNotFound        = errors.New("rol cliente_final no encontrado")
	ErrForbidden           = errors.New("acceso denegado: rol invalido")
	ErrNoItems             = errors.New("la orden debe tener al menos un item")
	ErrIntegrationNotFound = errors.New("integracion platform no encontrada para el negocio")
	ErrInvalidQuantity     = errors.New("la cantidad del item debe ser mayor a cero")
)
