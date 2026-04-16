package domain

import "errors"

var (
	ErrIntegrationNotFound     = errors.New("meli: integration not found")
	ErrInvalidCredentials      = errors.New("meli: invalid credentials")
	ErrMissingAccessToken      = errors.New("meli: missing access_token in credentials")
	ErrMissingRefreshToken     = errors.New("meli: missing refresh_token in credentials")
	ErrMissingAppID            = errors.New("meli: missing app_id in config")
	ErrMissingClientSecret     = errors.New("meli: missing client_secret in credentials")
	ErrTokenRefreshFailed      = errors.New("meli: failed to refresh access token")
	ErrTokenExpired            = errors.New("meli: access token expired")
	ErrNotificationInvalid     = errors.New("meli: invalid notification payload")
	ErrNotificationUnsupported = errors.New("meli: unsupported notification topic")
	ErrOrderNotFound           = errors.New("meli: order not found")
	ErrSellerIDNotFound        = errors.New("meli: seller_id not found in integration config")
	ErrRateLimited             = errors.New("meli: rate limited by API")
	ErrSignatureInvalid        = errors.New("meli: invalid notification signature")
)
