package errors

import "errors"

var (
	ErrPayUConfigNotFound = errors.New("payu integration type configuration not found")
	ErrPayUAPIError       = errors.New("payu API error")
	ErrInvalidCredentials = errors.New("invalid payu credentials")
)
