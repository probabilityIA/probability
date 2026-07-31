package domain

import "errors"

var (
	ErrIntegrationNotFound = errors.New("tiktok: integration not found")
	ErrMissingClientKey    = errors.New("tiktok: missing client_key in credentials")
	ErrMissingClientSecret = errors.New("tiktok: missing client_secret in credentials")
	ErrMissingShopName     = errors.New("tiktok: missing shop_name in config")
	ErrMissingBaseURL      = errors.New("tiktok: integration type has no base_url configured")
	ErrMissingBaseURLTest  = errors.New("tiktok: integration type has no base_url_test configured")
	ErrShopNotFound        = errors.New("tiktok: shop not found")
)
