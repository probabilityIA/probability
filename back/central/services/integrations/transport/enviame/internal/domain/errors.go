package domain

import "errors"

var (
	ErrNotImplemented      = errors.New("enviame: not yet implemented")
	ErrInvalidCredentials  = errors.New("enviame: invalid credentials")
	ErrShipmentNotFound    = errors.New("enviame: shipment not found")
	ErrQuoteFailed         = errors.New("enviame: quote failed")
	ErrGenerateFailed      = errors.New("enviame: guide generation failed")
	ErrTrackFailed         = errors.New("enviame: tracking failed")
	ErrCancelFailed        = errors.New("enviame: cancellation failed")
)
