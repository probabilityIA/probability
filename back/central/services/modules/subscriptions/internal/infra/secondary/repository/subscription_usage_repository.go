package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// GetSubscriptionUsage trae el plan vigente de un negocio junto con cuanto
// lleva consumido en el ciclo actual (envios/facturas/ordenes) y el pago
// pronosticado para la renovacion. Reutiliza el mismo conteo y calculo de
// excedente que ya usa el listado admin (ver admin_businesses_repository.go).
func (r *Repository) GetSubscriptionUsage(ctx context.Context, businessID uint) (*entities.SubscriptionUsage, error) {
	var sub models.BusinessSubscription
	err := r.db.Conn(ctx).
		Preload("SubscriptionType").
		Where("business_id = ? AND status = ?", businessID, entities.SubscriptionStatusPaid).
		Order("created_at desc").
		First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if sub.SubscriptionType == nil {
		return nil, nil
	}
	plan := *sub.SubscriptionType

	usage := &entities.SubscriptionUsage{
		PlanName:             plan.Name,
		PlanPrice:            plan.Price,
		BillingPeriod:        plan.BillingPeriod,
		ModuleCodes:          unmarshalModuleCodes(plan.Features),
		MaxEcommerceChannels: plan.MaxEcommerceChannels,
		CycleStartDate:       sub.StartDate,
		CycleEndDate:         sub.EndDate,
		IncludedShipments:    plan.IncludedShipments,
		ShipmentOveragePrice: plan.ShipmentOveragePrice,
		IncludedInvoices:     plan.IncludedInvoices,
		InvoiceOveragePrice:  plan.InvoiceOveragePrice,
		IncludedOrders:       plan.IncludedOrders,
		OrderOveragePrice:    plan.OrderOveragePrice,
	}

	if plan.IncludedShipments != nil {
		count, cerr := r.countShipmentsInRange(ctx, businessID, sub.StartDate, sub.EndDate)
		if cerr != nil {
			return nil, cerr
		}
		usage.ShipmentsUsed = count
	}
	if plan.IncludedInvoices != nil {
		count, cerr := r.countInvoicesInRange(ctx, businessID, sub.StartDate, sub.EndDate)
		if cerr != nil {
			return nil, cerr
		}
		usage.InvoicesUsed = count
	}
	if plan.IncludedOrders != nil {
		count, cerr := r.countOrdersInRange(ctx, businessID, sub.StartDate, sub.EndDate)
		if cerr != nil {
			return nil, cerr
		}
		usage.OrdersUsed = count
	}

	forecast, ferr := r.forecastNextPayment(ctx, businessID, plan, sub.StartDate, sub.EndDate)
	if ferr != nil {
		return nil, ferr
	}
	usage.ForecastedPayment = &forecast

	return usage, nil
}
