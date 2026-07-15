package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

func (uc *UseCase) ListOrders(ctx context.Context, f dtos.OrdersFilter) ([]entities.CodOrder, int64, error) {
	orders, total, err := uc.repo.ListCodOrders(ctx, f)
	if err != nil {
		return nil, 0, err
	}

	dm := uc.discountMap(ctx, f.BusinessID)

	confirmed, err := uc.repo.ConfirmedCuts(ctx, f.BusinessID)
	if err != nil {
		return nil, 0, err
	}
	confirmedWeeks := map[string]bool{}
	for i := range confirmed {
		confirmedWeeks[confirmed[i].PeriodStart.Format("2006-01-02")] = true
	}

	for i := range orders {
		pct := dm[orders[i].Carrier]
		d, n := domain.ApplyDiscount(orders[i].CodTotal, pct)
		orders[i].DiscountPct = pct
		orders[i].Discount = d
		orders[i].Net = n
		orders[i].CodState = domain.PaymentState(orders[i].Status, orders[i].Paid)
		orders[i].Collected = orders[i].Paid
		if orders[i].Paid {
			orders[i].CutStatus = "confirmed"
		} else if orders[i].DeliveredAt != nil {
			ws, _ := domain.WeekBounds(*orders[i].DeliveredAt)
			if confirmedWeeks[ws.Format("2006-01-02")] {
				orders[i].CutStatus = "pending"
			}
		}
	}
	return orders, total, nil
}
