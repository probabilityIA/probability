package errors

import "errors"

var (
	// ErrAuthFailed se retorna cuando la autenticación con Factus falla
	ErrAuthFailed = errors.New("factus: authentication failed")

	// ErrInvoiceCreationFailed se retorna cuando la creación de factura falla
	ErrInvoiceCreationFailed = errors.New("factus: invoice creation failed")

	// ErrMissingCredentials se retorna cuando faltan credenciales requeridas
	ErrMissingCredentials = errors.New("factus: missing required credentials")

	// ErrTokenExpired se retorna cuando el token ha expirado y el refresh también falló
	ErrTokenExpired = errors.New("factus: token expired and refresh failed")
)
