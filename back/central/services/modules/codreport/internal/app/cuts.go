package app

import (
	"context"
	"fmt"
	"sort"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

func (uc *UseCase) ListCuts(ctx context.Context, businessID uint, isAdmin bool) ([]entities.PaymentCut, error) {
	cuts, err := uc.repo.ConfirmedCuts(ctx, businessID)
	if err != nil {
		return nil, err
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

func (uc *UseCase) DeleteCut(ctx context.Context, businessID uint, cutID uint) error {
	return uc.repo.DeleteCut(ctx, businessID, cutID)
}

func (uc *UseCase) CutOrders(ctx context.Context, businessID uint, cutID uint) ([]entities.CodOrder, error) {
	orders, err := uc.repo.CutOrders(ctx, businessID, cutID)
	if err != nil {
		return nil, err
	}
	dm := uc.discountMap(ctx, businessID)
	for i := range orders {
		pct := dm[orders[i].Carrier]
		d, n := domain.ApplyDiscount(orders[i].CodTotal, pct)
		orders[i].DiscountPct = pct
		orders[i].Discount = d
		orders[i].Net = n
		orders[i].CodState = domain.CodStateCollected
		orders[i].CutStatus = "confirmed"
	}
	return orders, nil
}

func (uc *UseCase) CreateDraft(ctx context.Context, d dtos.ConfirmCutDTO) (*entities.PaymentCut, error) {
	dm := uc.discountMap(ctx, d.BusinessID)

	payouts, err := uc.repo.PayoutOrders(ctx, d.BusinessID, d.OrderIDs)
	if err != nil {
		return nil, err
	}
	if len(payouts) == 0 {
		return nil, fmt.Errorf("no hay ordenes validas para el corte (entregadas y aun sin consignar)")
	}

	userName := d.UserName
	if userName == "" {
		userName = uc.repo.UserName(ctx, d.UserID)
	}

	shell := entities.PaymentCut{
		BusinessID:  d.BusinessID,
		PeriodStart: d.PeriodStart,
		PeriodEnd:   d.PeriodEnd,
		Status:      "draft",
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

	cut := entities.PaymentCut{
		ID:             cutID,
		BusinessID:     d.BusinessID,
		PeriodStart:    d.PeriodStart,
		PeriodEnd:      d.PeriodEnd,
		Status:         "draft",
		OrdersCount:    oc,
		TotalCollected: tc,
		TotalDiscount:  td,
		TotalNet:       tn,
		ByCarrier:      aggs,
	}
	if err := uc.repo.UpdateCutTotals(ctx, cut); err != nil {
		return nil, err
	}
	return &cut, nil
}

func (uc *UseCase) ConfirmCut(ctx context.Context, businessID uint, cutID uint, userID uint, userName string) error {
	if userName == "" {
		userName = uc.repo.UserName(ctx, userID)
	}
	return uc.repo.ConfirmDraftCut(ctx, businessID, cutID, userID, userName)
}
