package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (r *Repository) ComputeOverageAmount(ctx context.Context, businessID uint, plan *entities.SubscriptionType, cycleStart, cycleEnd time.Time) (float64, error) {
	if plan == nil {
		return 0, nil
	}

	var overage float64

	if plan.IncludedShipments != nil && plan.ShipmentOveragePrice != nil {
		count, err := r.countShipmentsInRange(ctx, businessID, cycleStart, cycleEnd)
		if err != nil {
			return 0, err
		}
		if extra := count - int64(*plan.IncludedShipments); extra > 0 {
			overage += float64(extra) * *plan.ShipmentOveragePrice
		}
	}

	if plan.IncludedInvoices != nil && plan.InvoiceOveragePrice != nil {
		count, err := r.countInvoicesInRange(ctx, businessID, cycleStart, cycleEnd)
		if err != nil {
			return 0, err
		}
		if extra := count - int64(*plan.IncludedInvoices); extra > 0 {
			overage += float64(extra) * *plan.InvoiceOveragePrice
		}
	}

	if plan.IncludedOrders != nil && plan.OrderOveragePrice != nil {
		count, err := r.countOrdersInRange(ctx, businessID, cycleStart, cycleEnd)
		if err != nil {
			return 0, err
		}
		if extra := count - int64(*plan.IncludedOrders); extra > 0 {
			overage += float64(extra) * *plan.OrderOveragePrice
		}
	}

	return overage, nil
}
