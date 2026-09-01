package domain

import "context"

type ShippingMargin struct {
	MarginAmount     float64
	InsuranceMargin  float64
	CODMarginPercent float64
	CODMarginAmount  float64
}

func (m ShippingMargin) CODMargin(codCarrierFee float64) float64 {
	if codCarrierFee <= 0 {
		return 0
	}
	if m.CODMarginAmount > 0 {
		return m.CODMarginAmount
	}
	if m.CODMarginPercent > 0 {
		return codCarrierFee * m.CODMarginPercent / 100.0
	}
	return 0
}

type IShippingMarginReader interface {
	Get(ctx context.Context, businessID uint, carrierCode string) (ShippingMargin, error)
}
