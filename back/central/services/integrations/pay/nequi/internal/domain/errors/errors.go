package errors

import "errors"

var (
	ErrNequiConfigNotFound = errors.New("nequi integration type configuration not found")
	ErrNequiAPIError       = errors.New("nequi API error")
	ErrInvalidCredentials  = errors.New("invalid nequi credentials")
)
