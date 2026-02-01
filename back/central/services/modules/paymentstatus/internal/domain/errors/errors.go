package errors

import "errors"

var (
	// ErrPaymentStatusNotFound cuando no se encuentra un estado de pago
	ErrPaymentStatusNotFound = errors.New("payment status not found")

	// ErrInvalidPaymentStatusCode cuando el código del estado de pago es inválido
	ErrInvalidPaymentStatusCode = errors.New("invalid payment status code")
)
