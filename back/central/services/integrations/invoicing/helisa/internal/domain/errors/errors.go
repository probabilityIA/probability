package errors

import "errors"

var (
	// ErrAuthentication indica un error de autenticación con Helisa
	ErrAuthentication = errors.New("helisa authentication failed")

	// ErrInvalidCredentials indica credenciales inválidas
	ErrInvalidCredentials = errors.New("invalid helisa credentials")

	// ErrInvoiceCreation indica un error al crear la factura en Helisa
	ErrInvoiceCreation = errors.New("failed to create invoice in helisa")

	// ErrMissingRequiredField indica que falta un campo requerido
	ErrMissingRequiredField = errors.New("missing required field for helisa")

	// ErrNotImplemented indica que la funcionalidad no está implementada aún
	ErrNotImplemented = errors.New("helisa: functionality not yet implemented")
)
