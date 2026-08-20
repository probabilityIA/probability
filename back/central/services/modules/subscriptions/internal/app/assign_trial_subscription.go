package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

const trialPlanCode = "trial"
const defaultTrialDurationDays = 15

func (uc *UseCase) AssignTrialSubscription(ctx context.Context, businessID uint) {
	trialPlan, err := uc.repo.GetSubscriptionTypeByCode(ctx, trialPlanCode)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("failed to load trial plan for new business")
		return
	}
	if trialPlan == nil {
		uc.log.Warn(ctx).Uint("business_id", businessID).Msg("trial plan not configured, skipping trial assignment")
		return
	}

	days := defaultTrialDurationDays
	if trialPlan.TrialDurationDays != nil && *trialPlan.TrialDurationDays > 0 {
		days = *trialPlan.TrialDurationDays
	}

	now := time.Now()
	endDate := now.AddDate(0, 0, days)
	notes := "trial asignado automaticamente al registrarse"

	sub := &entities.BusinessSubscription{
		BusinessID:           businessID,
		SubscriptionTypeID:   trialPlan.ID,
		SubscriptionTypeName: trialPlan.Name,
		Amount:               0,
		StartDate:            now,
		EndDate:              endDate,
		Status:               entities.SubscriptionStatusPaid,
		PaymentMethod:        entities.PaymentMethodCourtesy,
		Notes:                &notes,
	}

	if err := uc.repo.CreateSubscriptionAndActivate(ctx, sub, trialPlan.ID, endDate); err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("failed to assign trial subscription to new business")
	}
}
