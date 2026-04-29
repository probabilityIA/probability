package errors

import "errors"

var (
	ErrShippingMarginNotFound  = errors.New("shipping margin not found")
	ErrDuplicateCarrier        = errors.New("a shipping margin for this carrier already exists for this business")
	ErrInvalidCarrierCode      = errors.New("invalid carrier_code")
	ErrInvalidMargin           = errors.New("margin_amount and insurance_margin must be >= 0")
)
