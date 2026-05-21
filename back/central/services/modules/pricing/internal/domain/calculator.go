package domain

import "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"

func ResolveEffectivePrice(productID string, basePrice float64, groupPrice, clientPrice *float64, groupID *uint) entities.EffectivePrice {
	result := entities.EffectivePrice{
		ProductID:  productID,
		BasePrice:  basePrice,
		FinalPrice: basePrice,
		Source:     "base",
		GroupID:    groupID,
	}

	if groupPrice != nil {
		result.FinalPrice = *groupPrice
		result.Source = "group"
	}

	if clientPrice != nil {
		result.FinalPrice = *clientPrice
		result.Source = "client"
	}

	if result.FinalPrice < 0 {
		result.FinalPrice = 0
	}

	return result
}
