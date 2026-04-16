package domain

import "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"

// CalculatePrice calcula el precio final aplicando regla de cliente y descuento por cantidad.
// Es una función pura sin dependencias externas.
func CalculatePrice(basePrice float64, clientRule *entities.ClientPricingRule, qtyDiscount *entities.QuantityDiscount) entities.PriceResult {
	result := entities.PriceResult{
		BasePrice:     basePrice,
		AdjustedPrice: basePrice,
		FinalPrice:    basePrice,
	}

	// Paso 1: Aplicar regla de cliente
	if clientRule != nil && clientRule.IsActive {
		result.ClientAdjustmentApplied = true
		result.ClientAdjustmentType = clientRule.AdjustmentType
		result.ClientAdjustmentValue = clientRule.AdjustmentValue

		switch clientRule.AdjustmentType {
		case "percentage":
			// value negativo = descuento (ej: -10 -> price * 0.90)
			result.AdjustedPrice = basePrice * (1 + clientRule.AdjustmentValue/100)
		case "fixed":
			// value negativo = descuento (ej: -500 -> price - 500)
			result.AdjustedPrice = basePrice + clientRule.AdjustmentValue
		}

		if result.AdjustedPrice < 0 {
			result.AdjustedPrice = 0
		}
	}

	result.FinalPrice = result.AdjustedPrice

	// Paso 2: Aplicar descuento por cantidad SOBRE el precio ajustado
	if qtyDiscount != nil && qtyDiscount.IsActive {
		result.QuantityDiscountApplied = true
		result.QuantityDiscountPercent = qtyDiscount.DiscountPercent
		result.FinalPrice = result.AdjustedPrice * (1 - qtyDiscount.DiscountPercent/100)

		if result.FinalPrice < 0 {
			result.FinalPrice = 0
		}
	}

	return result
}
