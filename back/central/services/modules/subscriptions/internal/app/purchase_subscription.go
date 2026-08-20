package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
)

const walletConceptSubscription = "SUBSCRIPTION"

func (uc *UseCase) PurchaseSubscription(ctx context.Context, dto dtos.PurchaseSubscriptionDTO) (*entities.BusinessSubscription, error) {
	if dto.Months <= 0 {
		return nil, errs.ErrInvalidMonths
	}

	subType, err := uc.repo.GetSubscriptionType(ctx, dto.SubscriptionTypeID)
	if err != nil {
		return nil, err
	}
	if subType == nil {
		return nil, errs.ErrSubscriptionTypeNotFound
	}
	if !subType.Active {
		return nil, errs.ErrSubscriptionTypeInactive
	}
	if !subType.Payable {
		return nil, errs.ErrSubscriptionTypeNotPayable
	}
	if subType.BusinessID != nil && *subType.BusinessID != dto.BusinessID {
		return nil, errs.ErrSubscriptionTypeNotFound
	}

	amount := subType.Price * float64(dto.Months)

	balance, err := uc.wallet.GetBalance(ctx, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if balance < amount {
		return nil, errs.ErrInsufficientBalance
	}

	start, endDate, err := uc.computeSubscriptionWindow(ctx, dto.BusinessID, dto.Months, nil, subType)
	if err != nil {
		return nil, err
	}

	reference := fmt.Sprintf("SUB-%d-%s-%dM", dto.BusinessID, subType.Code, dto.Months)
	if err := uc.wallet.Debit(ctx, dto.BusinessID, amount, reference, walletConceptSubscription, dto.UserID); err != nil {
		return nil, err
	}

	sub := &entities.BusinessSubscription{
		BusinessID:           dto.BusinessID,
		SubscriptionTypeID:   subType.ID,
		SubscriptionTypeName: subType.Name,
		Months:               dto.Months,
		Amount:               amount,
		StartDate:            start,
		EndDate:              endDate,
		Status:               entities.SubscriptionStatusPaid,
	}

	if err := uc.repo.CreateSubscriptionAndActivate(ctx, sub, subType.ID, endDate); err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("business_id", dto.BusinessID).
			Float64("amount", amount).
			Str("reference", reference).
			Msg("wallet debited for subscription but subscription could not be activated, requires manual reconciliation")
		return nil, err
	}

	uc.deactivateExpiryAnnouncements(ctx, dto.BusinessID)

	return sub, nil
}

// computeSubscriptionWindow calcula el inicio/fin de un periodo de suscripcion.
// Los planes con TrialDurationDays configurado (ej. el plan "trial") se miden
// en dias, no en meses: asignarlos manualmente con "Months" produce fechas
// incorrectas (1 mes calendario no son los dias de duracion del trial).
func (uc *UseCase) computeSubscriptionWindow(ctx context.Context, businessID uint, months int, startOverride *time.Time, subType *entities.SubscriptionType) (time.Time, time.Time, error) {
	addPeriod := func(start time.Time) time.Time {
		if subType != nil && subType.TrialDurationDays != nil && *subType.TrialDurationDays > 0 {
			return start.AddDate(0, 0, *subType.TrialDurationDays)
		}
		return start.AddDate(0, months, 0)
	}

	if startOverride != nil {
		start := *startOverride
		return start, addPeriod(start), nil
	}

	now := time.Now()
	start := now

	current, err := uc.repo.GetLatestByBusinessID(ctx, businessID)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if current != nil && current.EndDate.After(now) {
		start = current.EndDate
	}

	return start, addPeriod(start), nil
}
