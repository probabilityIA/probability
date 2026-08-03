package errors

import stderrs "errors"

var (
	ErrConfigNotFound     = stderrs.New("bancolombia integration type configuration not found")
	ErrAPIError           = stderrs.New("bancolombia API error")
	ErrInvalidCredentials = stderrs.New("invalid bancolombia credentials")
	ErrInvalidSignature   = stderrs.New("invalid bancolombia webhook signature")
	ErrAuthFailed         = stderrs.New("bancolombia oauth token request failed")
)
