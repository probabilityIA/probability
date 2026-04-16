package response

import "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"

// PriceCalculationResponse respuesta del cálculo de precio
type PriceCalculationResponse struct {
	BasePrice               float64 `json:"base_price"`
	AdjustedPrice           float64 `json:"adjusted_price"`
	FinalPrice              float64 `json:"final_price"`
	ClientAdjustmentApplied bool    `json:"client_adjustment_applied"`
	ClientAdjustmentType    string  `json:"client_adjustment_type,omitempty"`
	ClientAdjustmentValue   float64 `json:"client_adjustment_value,omitempty"`
	QuantityDiscountApplied bool    `json:"quantity_discount_applied"`
	QuantityDiscountPercent float64 `json:"quantity_discount_percent,omitempty"`
}

// FromPriceResult convierte un PriceResult a response
func FromPriceResult(r *entities.PriceResult) PriceCalculationResponse {
	return PriceCalculationResponse{
		BasePrice:               r.BasePrice,
		AdjustedPrice:           r.AdjustedPrice,
		FinalPrice:              r.FinalPrice,
		ClientAdjustmentApplied: r.ClientAdjustmentApplied,
		ClientAdjustmentType:    r.ClientAdjustmentType,
		ClientAdjustmentValue:   r.ClientAdjustmentValue,
		QuantityDiscountApplied: r.QuantityDiscountApplied,
		QuantityDiscountPercent: r.QuantityDiscountPercent,
	}
}
