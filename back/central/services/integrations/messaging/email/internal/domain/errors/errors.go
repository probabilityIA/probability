package errors

import "errors"

var (
	// ErrMissingRecipient se retorna cuando no se proporciona un email de destino
	ErrMissingRecipient = errors.New("missing recipient email address")
)
