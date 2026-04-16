package errors

import "errors"

var (
	ErrWompiConfigNotFound = errors.New("wompi integration type configuration not found")
	ErrWompiAPIError       = errors.New("wompi API error")
	ErrInvalidCredentials  = errors.New("invalid wompi credentials")
)
