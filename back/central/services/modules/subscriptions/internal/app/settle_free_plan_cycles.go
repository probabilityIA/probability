package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
)

func (uc *UseCase) SettleFreePlanCycles(ctx context.Context) error {
	now := time.Now()

	businessIDs, err := uc.repo.ListBusinessesWithEndedFreeCycle(ctx, now)
	if err != nil {
		return err
	}
	if len(businessIDs) == 0 {
		return nil
	}

	freePlan, err := uc.repo.GetSubscriptionTypeByCode(ctx, freePlanCode)
	if err != nil {
		return err
	}
	if freePlan == nil {
		return errs.ErrSubscriptionTypeNotFound
	}

	systemUserID, err := uc.resolveSystemUserID(ctx)
	if err != nil {
		return err
	}

	for _, businessID := range businessIDs {
		if err := uc.settleFreePlanCycle(ctx, businessID, freePlan, now, systemUserID); err != nil {
			uc.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("failed to settle free plan cycle")
		}
	}

	return nil
}

func (uc *UseCase) settleFreePlanCycle(ctx context.Context, businessID uint, freePlan *entities.SubscriptionType, now time.Time, systemUserID uint) error {
	usage, err := uc.repo.GetSubscriptionUsage(ctx, businessID)
	if err != nil {
		return err
	}
	if usage == nil {
		return nil
	}

	current, err := uc.repo.GetLatestByBusinessID(ctx, businessID)
	if err != nil {
		return err
	}
	if current == nil {
		return nil
	}

	var overageQty int64
	if usage.IncludedShipments != nil {
		overageQty = usage.ShipmentsUsed - int64(*usage.IncludedShipments)
	}

	if overageQty > 0 && usage.ShipmentOveragePrice != nil {
		amountDue := float64(overageQty) * *usage.ShipmentOveragePrice
		if err := uc.repo.SetOverageAmountDue(ctx, current.ID, amountDue); err != nil {
			return err
		}
		uc.recordAudit(ctx, businessID, systemUserID, entities.AuditActionOverageSettled,
			fmt.Sprintf("cierre de ciclo con %d envios de excedente, cargo de %.0f pendiente de pago", overageQty, amountDue))
		return nil
	}

	notes := "renovacion automatica del ciclo del plan gratuito"
	endDate := now.AddDate(0, 1, 0)
	return uc.activateNextFreeCycle(ctx, businessID, freePlan, now, endDate, &notes)
}

func (uc *UseCase) activateNextFreeCycle(ctx context.Context, businessID uint, freePlan *entities.SubscriptionType, start, end time.Time, notes *string) error {
	sub := &entities.BusinessSubscription{
		BusinessID:           businessID,
		SubscriptionTypeID:   freePlan.ID,
		SubscriptionTypeName: freePlan.Name,
		Months:               1,
		Amount:               0,
		StartDate:            start,
		EndDate:              end,
		Status:               entities.SubscriptionStatusPaid,
		PaymentMethod:        entities.PaymentMethodCourtesy,
		Notes:                notes,
	}
	return uc.repo.CreateSubscriptionAndActivate(ctx, sub, freePlan.ID, end)
}
