package app

import (
	"math"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/entities"
)

func (uc *UseCaseScore) scorePurchaseHistory(order *entities.ScoreOrder) (float64, []string) {
	history := order.CustomerHistory
	if history == nil || history.TotalOrders == 0 {
		return 50.0, nil // Neutral for new customers
	}

	var factors []string

	// Sub-signal 1: Order count (25%)
	orderCountScore := tierScore(float64(history.TotalOrders), []tier{
		{0, 0}, {1, 40}, {2, 70}, {5, 85}, {10, 100},
	})

	// Sub-signal 2: Total spent COP (20%)
	totalSpentScore := tierScore(history.TotalSpent, []tier{
		{0, 0}, {1, 30}, {100000, 60}, {500000, 80}, {1000000, 100},
	})

	// Sub-signal 3: Value consistency (10%)
	consistencyScore := 100.0
	if history.AvgOrderValue > 0 && order.TotalAmount > 0 {
		ratio := order.TotalAmount / history.AvgOrderValue
		if ratio > 3.0 {
			consistencyScore = 40.0
			factors = append(factors, "Valor de orden inusualmente alto vs historial")
		} else if ratio > 2.0 {
			consistencyScore = 70.0
		}
	}

	// Sub-signal 4: Customer tenure - days since first order (15%)
	tenureScore := 0.0
	if history.FirstOrderDate != nil {
		days := time.Since(*history.FirstOrderDate).Hours() / 24
		tenureScore = tierScore(days, []tier{
			{0, 0}, {1, 40}, {30, 60}, {90, 80}, {180, 100},
		})
	}

	// Sub-signal 5: Recency - days since last order (15%)
	recencyScore := 50.0
	if history.LastOrderDate != nil {
		days := time.Since(*history.LastOrderDate).Hours() / 24
		recencyScore = tierScoreDesc(days, []tier{
			{7, 100}, {30, 80}, {90, 60}, {180, 40}, {365, 20},
		})
	}

	// Sub-signal 6: Payment failure rate (15%)
	failureRate := 0.0
	if history.TotalOrders > 0 {
		failureRate = float64(history.FailedPayments) / float64(history.TotalOrders) * 100
	}
	failureScore := tierScoreDesc(failureRate, []tier{
		{0, 100}, {10, 70}, {25, 40}, {50, 10},
	})
	if failureRate >= 10 {
		factors = append(factors, "Alta tasa de fallos de pago")
	}

	// Weighted sum
	score := orderCountScore*0.25 + totalSpentScore*0.20 + consistencyScore*0.10 +
		tenureScore*0.15 + recencyScore*0.15 + failureScore*0.15

	return math.Round(score*100) / 100, factors
}

// tier represents a threshold-score pair
type tier struct {
	threshold float64
	score     float64
}

// tierScore returns a score based on ascending tiers (higher value = higher score)
func tierScore(value float64, tiers []tier) float64 {
	result := 0.0
	for _, t := range tiers {
		if value >= t.threshold {
			result = t.score
		}
	}
	return result
}

// tierScoreDesc returns a score based on descending tiers (higher value = lower score)
func tierScoreDesc(value float64, tiers []tier) float64 {
	for _, t := range tiers {
		if value <= t.threshold {
			return t.score
		}
	}
	// Beyond all thresholds
	if len(tiers) > 0 {
		return tiers[len(tiers)-1].score
	}
	return 50.0
}
