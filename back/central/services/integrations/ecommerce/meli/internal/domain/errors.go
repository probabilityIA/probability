package domain

import "errors"

var (
	ErrIntegrationNotFound  = errors.New("meli: integration not found")
	ErrInvalidCredentials   = errors.New("meli: invalid credentials")
	ErrMissingAccessToken   = errors.New("meli: missing access_token in credentials")
)
