package errors

import stderrs "errors"

var (
	ErrBoldConfigNotFound = stderrs.New("bold integration type configuration not found")
	ErrBoldAPIError       = stderrs.New("bold API error")
	ErrInvalidCredentials = stderrs.New("invalid bold credentials")
)
