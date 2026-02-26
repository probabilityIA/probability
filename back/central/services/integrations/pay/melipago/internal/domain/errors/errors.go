package errors

import "errors"

var (
	ErrMeliPagoConfigNotFound = errors.New("melipago integration type configuration not found")
	ErrMeliPagoAPIError       = errors.New("melipago API error")
	ErrInvalidCredentials     = errors.New("invalid melipago credentials")
)
