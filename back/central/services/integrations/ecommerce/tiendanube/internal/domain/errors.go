package domain

import "errors"

var (
	ErrIntegrationNotFound = errors.New("tiendanube: integration not found")
	ErrInvalidCredentials  = errors.New("tiendanube: invalid credentials")
	ErrMissingAccessToken  = errors.New("tiendanube: missing access_token in credentials")
	ErrMissingStoreID      = errors.New("tiendanube: missing store_id in config")
	ErrMissingBaseURL      = errors.New("tiendanube: integration type has no base_url configured")
	ErrMissingBaseURLTest  = errors.New("tiendanube: integration type has no base_url_test configured")
	ErrRateLimited         = errors.New("tiendanube: rate limited")
)
