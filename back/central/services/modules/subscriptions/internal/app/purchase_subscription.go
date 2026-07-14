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

	amount := subType.Price * float64(dto.Months)

	balance, err := uc.wallet.GetBalance(ctx, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if balance < amount {
		return nil, errs.ErrInsufficientBalance
	}

	start, endDate := uc.computeSubscriptionWindow(ctx, dto.BusinessID, dto.Months)

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

	if err := uc.repo.CreateBusinessSubscription(ctx, sub); err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("business_id", dto.BusinessID).
			Float64("amount", amount).
			Str("reference", reference).
			Msg("wallet debited for subscription but subscription record could not be created, requires manual reconciliation")
		return nil, err
	}

	if err := uc.repo.UpdateBusinessCurrentSubscriptionType(ctx, dto.BusinessID, subType.ID, entities.BusinessStatusActive, endDate); err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("business_id", dto.BusinessID).
			Uint("subscription_type_id", subType.ID).
			Msg("subscription created but business current subscription type could not be updated, requires manual reconciliation")
		return nil, err
	}

	uc.deactivateExpiryAnnouncements(ctx, dto.BusinessID)

	return sub, nil
}

func (uc *UseCase) computeSubscriptionWindow(ctx context.Context, businessID uint, months int) (time.Time, time.Time) {
	now := time.Now()
	start := now

	current, err := uc.repo.GetLatestByBusinessID(ctx, businessID)
	if err == nil && current != nil && current.EndDate.After(now) {
		start = current.EndDate
	}

	return start, start.AddDate(0, months, 0)
}
