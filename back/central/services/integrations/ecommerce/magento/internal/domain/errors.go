package domain

import "errors"

var (
	ErrIntegrationNotFound = errors.New("magento: integration not found")
	ErrInvalidCredentials  = errors.New("magento: invalid credentials")
	ErrMissingAccessToken  = errors.New("magento: missing access_token in credentials")
	ErrMissingStoreURL     = errors.New("magento: missing store_url in config")
)
