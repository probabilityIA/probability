package domain

import "errors"

var (
	ErrIntegrationNotFound = errors.New("amazon: integration not found")
	ErrInvalidCredentials  = errors.New("amazon: invalid credentials")
	ErrMissingSellerID     = errors.New("amazon: missing seller_id in config")
	ErrMissingRefreshToken = errors.New("amazon: missing refresh_token in credentials")
	ErrMissingClientID     = errors.New("amazon: missing client_id in credentials")
	ErrMissingClientSecret = errors.New("amazon: missing client_secret in credentials")
)
