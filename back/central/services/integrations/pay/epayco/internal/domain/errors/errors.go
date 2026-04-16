package errors

import "errors"

var (
	ErrEPaycoConfigNotFound = errors.New("epayco integration type configuration not found")
	ErrEPaycoAPIError       = errors.New("epayco API error")
	ErrInvalidCredentials   = errors.New("invalid epayco credentials")
)
