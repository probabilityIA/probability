package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
)

func (uc *UseCase) RegisterPayment(ctx context.Context, dto dtos.RegisterPaymentDTO) (*entities.BusinessSubscription, error) {
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

	amount := subType.Price * float64(dto.Months)
	start, endDate := uc.computeSubscriptionWindow(ctx, dto.BusinessID, dto.Months)

	sub := &entities.BusinessSubscription{
		BusinessID:           dto.BusinessID,
		SubscriptionTypeID:   subType.ID,
		SubscriptionTypeName: subType.Name,
		Months:               dto.Months,
		Amount:               amount,
		StartDate:            start,
		EndDate:              endDate,
		Status:               entities.SubscriptionStatusPaid,
		PaymentReference:     dto.PaymentReference,
		Notes:                dto.Notes,
	}

	if err := uc.repo.CreateBusinessSubscription(ctx, sub); err != nil {
		return nil, err
	}

	if err := uc.repo.UpdateBusinessCurrentSubscriptionType(ctx, dto.BusinessID, subType.ID, entities.BusinessStatusActive, endDate); err != nil {
		return nil, err
	}

	uc.deactivateExpiryAnnouncements(ctx, dto.BusinessID)

	return sub, nil
}
