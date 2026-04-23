package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
)

func (uc *useCase) Preview(ctx context.Context, filter dtos.BackfillFilter) (*dtos.PreviewResponse, error) {
	selector, ok := uc.registry.Get(filter.EventCode)
	if !ok {
		return nil, fmt.Errorf("event_code no soportado: %s", filter.EventCode)
	}

	candidates, err := selector.Preview(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("preview failed: %w", err)
	}

	breakdown := make(map[uint]int)
	for _, c := range candidates {
		breakdown[c.BusinessID]++
	}

	sampleSize := len(candidates)
	if sampleSize > 10 {
		sampleSize = 10
	}
	sample := make([]dtos.CandidateResponse, 0, sampleSize)
	for i := 0; i < sampleSize; i++ {
		c := candidates[i]
		sample = append(sample, dtos.CandidateResponse{
			OrderID:        c.OrderID,
			OrderNumber:    c.OrderNumber,
			BusinessID:     c.BusinessID,
			CustomerPhone:  c.CustomerPhone,
			TrackingNumber: c.TrackingNumber,
			Status:         c.Status,
		})
	}

	return &dtos.PreviewResponse{
		EventCode:      filter.EventCode,
		TotalEligible:  len(candidates),
		BreakdownByBiz: breakdown,
		Sample:         sample,
	}, nil
}
