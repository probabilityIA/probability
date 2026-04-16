package errors

import "errors"

var (
	// Payment Method Errors
	ErrPaymentMethodNotFound         = errors.New("payment method not found")
	ErrPaymentMethodCodeAlreadyExists = errors.New("payment method with this code already exists")
	ErrPaymentMethodHasActiveMappings = errors.New("cannot delete payment method with active mappings")

	// Payment Mapping Errors
	ErrPaymentMappingNotFound       = errors.New("payment mapping not found")
	ErrPaymentMappingAlreadyExists  = errors.New("mapping already exists for this integration type and original method")
	ErrInvalidIntegrationType       = errors.New("invalid integration type")
)
