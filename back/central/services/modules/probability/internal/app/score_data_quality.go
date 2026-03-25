package app

import "github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/entities"

// scoreDataQuality calculates Category 1: Data Quality score (0-100)
// Uses the existing 7 factors, normalized to a 0-100 scale
func (uc *UseCaseScore) scoreDataQuality(order *entities.ScoreOrder) (float64, []string) {
	factors := uc.GetStaticNegativeFactors(order)
	pointsPerFactor := 100.0 / 7.0
	score := 100.0 - (float64(len(factors)) * pointsPerFactor)
	if score < 0 {
		score = 0
	}
	return score, factors
}
