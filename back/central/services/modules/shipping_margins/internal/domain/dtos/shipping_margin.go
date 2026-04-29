package dtos

type ListShippingMarginsParams struct {
	BusinessID  uint
	CarrierCode string
	Page        int
	PageSize    int
}

func (p ListShippingMarginsParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type CreateShippingMarginDTO struct {
	BusinessID      uint
	CarrierCode     string
	CarrierName     string
	MarginAmount    float64
	InsuranceMargin float64
	IsActive        bool
}

type UpdateShippingMarginDTO struct {
	ID              uint
	BusinessID      uint
	CarrierName     string
	MarginAmount    float64
	InsuranceMargin float64
	IsActive        bool
}
