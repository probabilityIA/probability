package app

import (
	"context"
	"fmt"
	"sort"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/entities"
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

	grouped := make(map[uint][]entities.Candidate)
	for _, c := range candidates {
		grouped[c.BusinessID] = append(grouped[c.BusinessID], c)
	}

	ids := make([]uint, 0, len(grouped))
	for id := range grouped {
		ids = append(ids, id)
	}

	names := map[uint]string{}
	if uc.businessResolver != nil && len(ids) > 0 {
		if resolved, rerr := uc.businessResolver.ResolveNames(ctx, ids); rerr == nil {
			names = resolved
		} else {
			uc.log.Warn(ctx).Err(rerr).Msg("Failed to resolve business names")
		}
	}

	businesses := make([]dtos.BusinessGroup, 0, len(grouped))
	for id, list := range grouped {
		orders := make([]dtos.OrderCandidateResponse, 0, len(list))
		for _, c := range list {
			orders = append(orders, dtos.OrderCandidateResponse{
				OrderID:        c.OrderID,
				OrderNumber:    c.OrderNumber,
				CustomerPhone:  c.CustomerPhone,
				TrackingNumber: c.TrackingNumber,
				Status:         c.Status,
				Carrier:        c.Carrier,
				CarrierLogoURL: c.CarrierLogoURL,
			})
		}
		businesses = append(businesses, dtos.BusinessGroup{
			BusinessID:   id,
			BusinessName: names[id],
			Count:        len(list),
			Orders:       orders,
		})
	}

	sort.Slice(businesses, func(i, j int) bool {
		if businesses[i].Count != businesses[j].Count {
			return businesses[i].Count > businesses[j].Count
		}
		return businesses[i].BusinessID < businesses[j].BusinessID
	})

	return &dtos.PreviewResponse{
		EventCode:     filter.EventCode,
		TotalEligible: len(candidates),
		Businesses:    businesses,
	}, nil
}
