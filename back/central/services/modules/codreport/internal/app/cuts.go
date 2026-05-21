package app

import (
	"context"
	"sort"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

func (uc *UseCase) ListCuts(ctx context.Context, businessID uint, isAdmin bool) ([]entities.PaymentCut, error) {
	dm := uc.discountMap(ctx, businessID)

	weekAggs, err := uc.repo.WeeklyAggregates(ctx, businessID, 12)
	if err != nil {
		return nil, err
	}
	confirmed, err := uc.repo.ConfirmedCuts(ctx, businessID)
	if err != nil {
		return nil, err
	}

	confirmedByWeek := map[string]entities.PaymentCut{}
	for i := range confirmed {
		confirmedByWeek[confirmed[i].PeriodStart.Format("2006-01-02")] = confirmed[i]
	}

	weekMap := map[string][]entities.CarrierAggregate{}
	weekStart := map[string]time.Time{}
	order := []string{}
	for i := range weekAggs {
		key := weekAggs[i].WeekStart.Format("2006-01-02")
		if _, ok := weekMap[key]; !ok {
			order = append(order, key)
			weekStart[key] = weekAggs[i].WeekStart
		}
		weekMap[key] = append(weekMap[key], entities.CarrierAggregate{
			Carrier:        weekAggs[i].Carrier,
			OrdersCount:    weekAggs[i].Orders,
			TotalCollected: weekAggs[i].Collected,
		})
	}

	cuts := []entities.PaymentCut{}
	seen := map[string]bool{}
	for _, key := range order {
		seen[key] = true
		if c, ok := confirmedByWeek[key]; ok {
			cuts = append(cuts, c)
			continue
		}
		aggs := weekMap[key]
		domain.EnrichCarrierAggregates(aggs, dm)
		oc, tc, td, tn := domain.SumAggregates(aggs)
		ws := weekStart[key]
		cuts = append(cuts, entities.PaymentCut{
			BusinessID:     businessID,
			PeriodStart:    ws,
			PeriodEnd:      ws.AddDate(0, 0, 6),
			Status:         "pending",
			OrdersCount:    oc,
			TotalCollected: tc,
			TotalDiscount:  td,
			TotalNet:       tn,
			ByCarrier:      aggs,
		})
	}
	for i := range confirmed {
		key := confirmed[i].PeriodStart.Format("2006-01-02")
		if !seen[key] {
			cuts = append(cuts, confirmed[i])
		}
	}

	sort.Slice(cuts, func(i, j int) bool {
		return cuts[i].PeriodStart.After(cuts[j].PeriodStart)
	})

	if !isAdmin {
		filtered := []entities.PaymentCut{}
		for i := range cuts {
			if cuts[i].Status == "confirmed" {
				filtered = append(filtered, cuts[i])
			}
		}
		return filtered, nil
	}
	return cuts, nil
}

func (uc *UseCase) ConfirmCut(ctx context.Context, d dtos.ConfirmCutDTO) (*entities.PaymentCut, error) {
	dm := uc.discountMap(ctx, d.BusinessID)

	aggs, err := uc.repo.CutPeriodOrders(ctx, d.BusinessID, d.PeriodStart, d.PeriodEnd)
	if err != nil {
		return nil, err
	}
	domain.EnrichCarrierAggregates(aggs, dm)
	oc, tc, td, tn := domain.SumAggregates(aggs)

	cut := entities.PaymentCut{
		BusinessID:     d.BusinessID,
		PeriodStart:    d.PeriodStart,
		PeriodEnd:      d.PeriodEnd,
		Status:         "confirmed",
		OrdersCount:    oc,
		TotalCollected: tc,
		TotalDiscount:  td,
		TotalNet:       tn,
		ByCarrier:      aggs,
	}

	userName := d.UserName
	if userName == "" {
		userName = uc.repo.UserName(ctx, d.UserID)
	}
	return uc.repo.SaveConfirmedCut(ctx, cut, d.UserID, userName)
}
