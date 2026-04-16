package domain

import "errors"

var (
	ErrIntegrationNotFound = errors.New("tiendanube: integration not found")
	ErrInvalidCredentials  = errors.New("tiendanube: invalid credentials")
	ErrMissingAccessToken  = errors.New("tiendanube: missing access_token in credentials")
	ErrMissingStoreURL     = errors.New("tiendanube: missing store_url in config")
)
