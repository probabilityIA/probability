package domain

import "errors"

var (
	ErrNotImplemented      = errors.New("mipaquete: not yet implemented")
	ErrInvalidCredentials  = errors.New("mipaquete: invalid credentials")
	ErrShipmentNotFound    = errors.New("mipaquete: shipment not found")
	ErrQuoteFailed         = errors.New("mipaquete: quote failed")
	ErrGenerateFailed      = errors.New("mipaquete: guide generation failed")
	ErrTrackFailed         = errors.New("mipaquete: tracking failed")
	ErrCancelFailed        = errors.New("mipaquete: cancellation failed")
)
