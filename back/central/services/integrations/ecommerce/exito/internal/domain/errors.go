package domain

import "errors"

var (
	ErrIntegrationNotFound = errors.New("exito: integration not found")
	ErrInvalidCredentials  = errors.New("exito: invalid credentials")
	ErrMissingAPIKey       = errors.New("exito: missing api_key in credentials")
	ErrMissingSellerID     = errors.New("exito: missing seller_id in config")
)
