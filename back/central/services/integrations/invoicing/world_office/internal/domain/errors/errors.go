package errors

import "errors"

var (
	// ErrAuthentication indica un error de autenticación con World Office
	ErrAuthentication = errors.New("world office authentication failed")

	// ErrInvalidCredentials indica credenciales inválidas
	ErrInvalidCredentials = errors.New("invalid world office credentials")

	// ErrInvoiceCreation indica un error al crear la factura en World Office
	ErrInvoiceCreation = errors.New("failed to create invoice in world office")

	// ErrMissingRequiredField indica que falta un campo requerido
	ErrMissingRequiredField = errors.New("missing required field for world office")

	// ErrNotImplemented indica que la funcionalidad no está implementada aún
	ErrNotImplemented = errors.New("world_office: functionality not yet implemented")
)
