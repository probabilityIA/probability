package domain

import "errors"

var (
	ErrIntegrationNotFound = errors.New("falabella: integration not found")
	ErrInvalidCredentials  = errors.New("falabella: invalid credentials")
	ErrMissingAPIKey       = errors.New("falabella: missing api_key in credentials")
	ErrMissingUserID       = errors.New("falabella: missing user_id in config")
)
