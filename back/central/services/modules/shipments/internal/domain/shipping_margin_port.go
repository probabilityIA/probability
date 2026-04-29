package domain

import "context"

type ShippingMargin struct {
	MarginAmount    float64
	InsuranceMargin float64
}

type IShippingMarginReader interface {
	Get(ctx context.Context, businessID uint, carrierCode string) (ShippingMargin, error)
}
