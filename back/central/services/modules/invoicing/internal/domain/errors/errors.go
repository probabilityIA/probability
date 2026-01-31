package errors

import "errors"

// Errores de validación
var (
	ErrInvalidInvoiceData     = errors.New("invalid invoice data")
	ErrInvalidProviderConfig  = errors.New("invalid provider configuration")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrInvalidAmount          = errors.New("invalid amount")
	ErrInvalidCurrency        = errors.New("invalid currency")
	ErrInvalidCustomerData    = errors.New("invalid customer data")
	ErrInvalidOrderData       = errors.New("invalid order data")
	ErrMissingRequiredField   = errors.New("missing required field")
)

// Errores de estado
var (
	ErrInvoiceAlreadyExists   = errors.New("invoice already exists for this order")
	ErrInvoiceNotFound        = errors.New("invoice not found")
	ErrInvoiceAlreadyIssued   = errors.New("invoice already issued")
	ErrInvoiceAlreadyCancelled = errors.New("invoice already cancelled")
	ErrInvoiceCannotBeCancelled = errors.New("invoice cannot be cancelled")
	ErrOrderNotInvoiceable    = errors.New("order is not invoiceable")
	ErrOrderAlreadyInvoiced   = errors.New("order already has an invoice")
)

// Errores de proveedor
var (
	ErrProviderNotFound       = errors.New("invoicing provider not found")
	ErrProviderNotActive      = errors.New("invoicing provider is not active")
	ErrProviderTypeNotFound   = errors.New("invoicing provider type not found")
	ErrProviderNotConfigured  = errors.New("invoicing provider not configured for this integration")
	ErrProviderAPIError       = errors.New("provider API error")
	ErrProviderTimeout        = errors.New("provider timeout")
	ErrProviderUnauthorized   = errors.New("provider unauthorized")
	ErrProviderRateLimitExceeded = errors.New("provider rate limit exceeded")
)

// Errores de configuración
var (
	ErrConfigNotFound         = errors.New("invoicing config not found")
	ErrConfigNotEnabled       = errors.New("invoicing config is not enabled")
	ErrConfigAlreadyExists    = errors.New("invoicing config already exists for this integration")
	ErrAutoInvoiceNotEnabled  = errors.New("auto invoice is not enabled")
)

// Errores de filtros
var (
	ErrOrderBelowMinAmount    = errors.New("order amount is below minimum threshold")
	ErrOrderNotPaid           = errors.New("order is not paid")
	ErrPaymentMethodNotAllowed = errors.New("payment method is not allowed")
	ErrOrderTypeNotAllowed    = errors.New("order type is not allowed")
	ErrOrderStatusExcluded    = errors.New("order status is excluded from invoicing")
)

// Errores de sincronización
var (
	ErrSyncFailed             = errors.New("synchronization failed")
	ErrMaxRetriesExceeded     = errors.New("maximum retries exceeded")
	ErrRetryNotAllowed        = errors.New("retry not allowed")
	ErrSyncInProgress         = errors.New("synchronization already in progress")
)

// Errores de notas de crédito
var (
	ErrCreditNoteNotFound     = errors.New("credit note not found")
	ErrCreditNoteAlreadyIssued = errors.New("credit note already issued")
	ErrCreditNoteAmountExceeds = errors.New("credit note amount exceeds invoice total")
	ErrInvoiceNotIssued       = errors.New("invoice must be issued before creating credit note")
)

// Errores de encriptación
var (
	ErrEncryptionFailed       = errors.New("encryption failed")
	ErrDecryptionFailed       = errors.New("decryption failed")
	ErrInvalidEncryptionKey   = errors.New("invalid encryption key")
)

// Errores de autenticación con proveedor
var (
	ErrAuthenticationFailed   = errors.New("authentication failed with provider")
	ErrTokenExpired           = errors.New("provider token expired")
	ErrTokenRefreshFailed     = errors.New("token refresh failed")
)
