package repository

import "time"

func calculateAvgTicket(totalSpent float64, totalOrders int) float64 {
	if totalOrders <= 0 {
		return 0
	}
	return totalSpent / float64(totalOrders)
}

func coalesceFloat(preferred, fallback float64) float64 {
	if preferred > 0 {
		return preferred
	}
	return fallback
}

func coalesceTime(existing, incoming *time.Time) *time.Time {
	if existing != nil {
		return existing
	}
	return incoming
}

func latestTime(a, b *time.Time) *time.Time {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if b.After(*a) {
		return b
	}
	return a
}

func coalesceString(preferred, fallback string) string {
	if preferred != "" {
		return preferred
	}
	return fallback
}

func paginationOffset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * pageSize
}
