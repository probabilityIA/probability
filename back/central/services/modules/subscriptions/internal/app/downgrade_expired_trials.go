package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
)

const freePlanCode = "free"

func (uc *UseCase) DowngradeExpiredTrials(ctx context.Context) error {
	now := time.Now()

	businessIDs, err := uc.repo.ListBusinessesWithExpiredTrial(ctx, now)
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
		if err := uc.downgradeTrialToFree(ctx, businessID, freePlan, now, systemUserID); err != nil {
			uc.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("failed to downgrade expired trial to free plan")
		}
	}

	return nil
}

func (uc *UseCase) downgradeTrialToFree(ctx context.Context, businessID uint, freePlan *entities.SubscriptionType, now time.Time, systemUserID uint) error {
	endDate := now.AddDate(0, 1, 0)
	notes := "degradacion automatica: el plan de prueba vencio"

	sub := &entities.BusinessSubscription{
		BusinessID:           businessID,
		SubscriptionTypeID:   freePlan.ID,
		SubscriptionTypeName: freePlan.Name,
		Months:               1,
		Amount:               0,
		StartDate:            now,
		EndDate:              endDate,
		Status:               entities.SubscriptionStatusPaid,
		PaymentMethod:        entities.PaymentMethodCourtesy,
		Notes:                &notes,
	}

	if err := uc.repo.CreateSubscriptionAndActivate(ctx, sub, freePlan.ID, endDate); err != nil {
		return err
	}

	uc.recordAudit(ctx, businessID, systemUserID, entities.AuditActionTrialDowngraded,
		fmt.Sprintf("plan de prueba vencido, degradado automaticamente a %s", freePlan.Name))

	return nil
}
