package domain

import "errors"

var (
	ErrIntegrationNotFound = errors.New("vtex: integration not found")
	ErrInvalidCredentials  = errors.New("vtex: invalid credentials")
	ErrMissingAPIKey       = errors.New("vtex: missing api_key in credentials")
	ErrMissingAPIToken     = errors.New("vtex: missing api_token in credentials")
	ErrMissingStoreURL     = errors.New("vtex: missing store_url in config")
	ErrOrderNotFound       = errors.New("vtex: order not found")
	ErrRateLimited         = errors.New("vtex: rate limited, too many requests")
)
