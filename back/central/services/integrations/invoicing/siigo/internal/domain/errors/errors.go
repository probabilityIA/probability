package errors

import "errors"

var (
	// ErrAuthentication indica un error de autenticación con Siigo
	ErrAuthentication = errors.New("siigo authentication failed")

	// ErrInvalidCredentials indica credenciales inválidas
	ErrInvalidCredentials = errors.New("invalid siigo credentials")

	// ErrCustomerNotFound indica que el cliente no fue encontrado en Siigo
	ErrCustomerNotFound = errors.New("customer not found in siigo")

	// ErrInvoiceCreation indica un error al crear la factura en Siigo
	ErrInvoiceCreation = errors.New("failed to create invoice in siigo")

	// ErrMissingRequiredField indica que falta un campo requerido
	ErrMissingRequiredField = errors.New("missing required field for siigo")
)
