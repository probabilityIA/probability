package app

import (
	"context"
	"time"
)

func (uc *UseCase) GetPaymentWindow(ctx context.Context, businessID uint) (*time.Time, *time.Time, error) {
	meta, err := uc.repo.GetBusinessSubscriptionMeta(ctx, businessID)
	if err != nil {
		return nil, nil, err
	}
	if meta == nil || meta.EndDate == nil {
		return nil, nil, nil
	}

	start := *meta.EndDate
	end := computeCutoffDate(*meta.EndDate, meta.CutoffDay, meta.CourtesyUntil)
	return &start, &end, nil
}
