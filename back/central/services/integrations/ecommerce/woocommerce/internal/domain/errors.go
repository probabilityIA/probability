package domain

import "errors"

var (
	ErrIntegrationNotFound  = errors.New("woocommerce: integration not found")
	ErrInvalidCredentials   = errors.New("woocommerce: invalid credentials")
	ErrMissingConsumerKey   = errors.New("woocommerce: missing consumer_key in credentials")
	ErrMissingConsumerSecret = errors.New("woocommerce: missing consumer_secret in credentials")
	ErrMissingStoreURL      = errors.New("woocommerce: missing store_url in config")
)
