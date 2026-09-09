package errors

import "errors"

var (
	ErrProductNotFound      = errors.New("producto no encontrado")
	ErrOrderNotFound        = errors.New("orden no encontrada")
	ErrClientNotFound       = errors.New("cliente no encontrado")
	ErrBusinessNotFound     = errors.New("negocio no encontrado")
	ErrEmailAlreadyExists   = errors.New("ya existe un usuario con este email")
	ErrPasswordMismatch     = errors.New("este email ya tiene una cuenta; usa tu contrasena actual para vincularla a esta tienda")
	ErrRoleNotFound         = errors.New("rol cliente_final no encontrado")
	ErrRoleNotAllowed       = errors.New("tu rol no tiene permiso para crear clientes del catalogo")
	ErrForbidden            = errors.New("acceso denegado: rol invalido")
	ErrNoItems              = errors.New("la orden debe tener al menos un item")
	ErrIntegrationNotFound  = errors.New("integracion platform no encontrada para el negocio")
	ErrInvalidQuantity      = errors.New("la cantidad del item debe ser mayor a cero")
	ErrStorefrontNotActive  = errors.New("el catalogo no esta activo para este negocio")
	ErrInsufficientStock    = errors.New("stock insuficiente para completar el pedido")
)
