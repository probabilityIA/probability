package request

type CreateShippingMarginRequest struct {
	CarrierCode     string  `json:"carrier_code" binding:"required,min=2,max=50"`
	CarrierName     string  `json:"carrier_name" binding:"required,min=2,max=100"`
	MarginAmount    float64 `json:"margin_amount" binding:"gte=0"`
	InsuranceMargin float64 `json:"insurance_margin" binding:"gte=0"`
	IsActive        *bool   `json:"is_active"`
}

type UpdateShippingMarginRequest struct {
	CarrierName     string  `json:"carrier_name" binding:"required,min=2,max=100"`
	MarginAmount    float64 `json:"margin_amount" binding:"gte=0"`
	InsuranceMargin float64 `json:"insurance_margin" binding:"gte=0"`
	IsActive        *bool   `json:"is_active"`
}
