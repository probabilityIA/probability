package app

import "github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/entities"

// scoreDataQuality calculates Category 1: Data Quality score (0-100)
// Uses the existing 6 factors, normalized to a 0-100 scale
func (uc *UseCaseScore) scoreDataQuality(order *entities.ScoreOrder) (float64, []string) {
	factors := uc.GetStaticNegativeFactors(order)
	pointsPerFactor := 100.0 / 6.0
	score := 100.0 - (float64(len(factors)) * pointsPerFactor)

	switch order.ShippingGeoConfidence {
	case "low":
		score -= 30
		factors = append(factors, "Direccion no verificada (dudosa)")
	case "medium":
		score -= 10
		factors = append(factors, "Direccion aproximada, verificar")
	}

	if score < 0 {
		score = 0
	}
	return score, factors
}
