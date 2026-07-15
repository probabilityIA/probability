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

func (uc *UseCase) SelectableOrders(ctx context.Context, f dtos.SelectableOrdersFilter) ([]entities.CodOrder, error) {
	orders, err := uc.repo.SelectableCutOrders(ctx, f)
	if err != nil {
		return nil, err
	}
	dm := uc.discountMap(ctx, f.BusinessID)
	for i := range orders {
		pct := dm[orders[i].Carrier]
		d, n := domain.ApplyDiscount(orders[i].CodTotal, pct)
		orders[i].DiscountPct = pct
		orders[i].Discount = d
		orders[i].Net = n
		orders[i].CodState = domain.CodStatePendingPayment
	}
	return orders, nil
}

func (uc *UseCase) ConfirmCut(ctx context.Context, d dtos.ConfirmCutDTO) (*entities.PaymentCut, error) {
	dm := uc.discountMap(ctx, d.BusinessID)

	payouts, err := uc.repo.PayoutOrders(ctx, d.BusinessID, d.OrderIDs)
	if err != nil {
		return nil, err
	}

	userName := d.UserName
	if userName == "" {
		userName = uc.repo.UserName(ctx, d.UserID)
	}

	shell := entities.PaymentCut{
		BusinessID:  d.BusinessID,
		PeriodStart: d.PeriodStart,
		PeriodEnd:   d.PeriodEnd,
	}
	cutID, err := uc.repo.UpsertCutOrders(ctx, shell, payouts, d.UserID, userName)
	if err != nil {
		return nil, err
	}

	aggs, err := uc.repo.PaidAggregatesForCut(ctx, cutID)
	if err != nil {
		return nil, err
	}
	domain.EnrichCarrierAggregates(aggs, dm)
	oc, tc, td, tn := domain.SumAggregates(aggs)

	now := time.Now().UTC()
	cut := entities.PaymentCut{
		ID:              cutID,
		BusinessID:      d.BusinessID,
		PeriodStart:     d.PeriodStart,
		PeriodEnd:       d.PeriodEnd,
		Status:          "confirmed",
		OrdersCount:     oc,
		TotalCollected:  tc,
		TotalDiscount:   td,
		TotalNet:        tn,
		ByCarrier:       aggs,
		ConfirmedBy:     d.UserID,
		ConfirmedByName: userName,
		ConfirmedAt:     &now,
	}
	if err := uc.repo.UpdateCutTotals(ctx, cut); err != nil {
		return nil, err
	}
	return &cut, nil
}
