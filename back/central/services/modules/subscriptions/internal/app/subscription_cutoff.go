package app

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

const defaultCutoffGraceDays = 6

func cutoffReached(business entities.ExpiringBusiness, now time.Time) bool {
	cutoffDate := business.EndDate.AddDate(0, 0, defaultCutoffGraceDays)
	if business.CutoffDay != nil {
		cutoffDate = nextCutoffOnOrAfter(business.EndDate, *business.CutoffDay)
	}
	if business.CourtesyUntil != nil && business.CourtesyUntil.After(cutoffDate) {
		cutoffDate = *business.CourtesyUntil
	}
	return !now.Before(cutoffDate)
}

func nextCutoffOnOrAfter(after time.Time, cutoffDay int) time.Time {
	candidate := clampedCutoffDate(after.Year(), after.Month(), cutoffDay, after.Location())
	if candidate.Before(after) {
		candidate = clampedCutoffDate(after.Year(), after.Month()+1, cutoffDay, after.Location())
	}
	return candidate
}

func clampedCutoffDate(year int, month time.Month, day int, loc *time.Location) time.Time {
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, loc)
	lastDay := firstOfNextMonth.AddDate(0, 0, -1).Day()
	if day > lastDay {
		day = lastDay
	}
	if day < 1 {
		day = 1
	}
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}
