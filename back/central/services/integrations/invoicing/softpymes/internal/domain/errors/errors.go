package errors

import "errors"

// ═══════════════════════════════════════════════════════════════
// ERRORES DE DOMINIO - Softpymes Providers
// ═══════════════════════════════════════════════════════════════

var (
	// Provider errors
	ErrProviderNotFound       = errors.New("provider not found")
	ErrProviderAlreadyExists  = errors.New("provider already exists for this business and type")
	ErrProviderTypeNotFound   = errors.New("provider type not found")
	ErrProviderTypeInactive   = errors.New("provider type is not active")
	ErrInvalidProviderConfig  = errors.New("invalid provider configuration")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrDefaultProviderExists  = errors.New("a default provider already exists for this business")
	ErrCannotDeleteDefault    = errors.New("cannot delete default provider")
	ErrProviderInUse          = errors.New("provider is in use and cannot be deleted")

	// Validation errors
	ErrBusinessIDRequired     = errors.New("business ID is required")
	ErrProviderTypeRequired   = errors.New("provider type is required")
	ErrProviderNameRequired   = errors.New("provider name is required")
	ErrCredentialsRequired    = errors.New("credentials are required")
	ErrConfigRequired         = errors.New("configuration is required")

	// Connection errors
	ErrConnectionFailed       = errors.New("failed to connect to provider API")
	ErrAuthenticationFailed   = errors.New("authentication with provider failed")
	ErrAPIKeyRequired         = errors.New("API key is required in credentials")
	ErrAPISecretRequired      = errors.New("API secret is required in credentials")
)
