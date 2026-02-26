package errors

import "errors"

var (
	ErrStripeConfigNotFound = errors.New("stripe integration type configuration not found")
	ErrStripeAPIError       = errors.New("stripe API error")
	ErrInvalidCredentials   = errors.New("invalid stripe credentials")
)
