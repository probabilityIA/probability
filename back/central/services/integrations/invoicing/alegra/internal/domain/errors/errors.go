package errors

import "errors"

var (
	// ErrAuthentication indica un error de autenticación con Alegra
	ErrAuthentication = errors.New("alegra authentication failed")

	// ErrInvalidCredentials indica credenciales inválidas
	ErrInvalidCredentials = errors.New("invalid alegra credentials")

	// ErrInvoiceCreation indica un error al crear la factura en Alegra
	ErrInvoiceCreation = errors.New("failed to create invoice in alegra")

	// ErrMissingRequiredField indica que falta un campo requerido
	ErrMissingRequiredField = errors.New("missing required field for alegra")

	// ErrNotImplemented indica que la funcionalidad no está implementada aún
	ErrNotImplemented = errors.New("alegra: functionality not yet implemented")
)
