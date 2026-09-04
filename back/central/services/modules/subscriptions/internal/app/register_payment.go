package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
)

func (uc *UseCase) RegisterPayment(ctx context.Context, dto dtos.RegisterPaymentDTO, actorUserID uint) (*entities.BusinessSubscription, error) {
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

	if dto.StartDate == nil {
		current, err := uc.repo.GetLatestByBusinessID(ctx, dto.BusinessID)
		if err != nil {
			return nil, err
		}
		meta, err := uc.repo.GetBusinessSubscriptionMeta(ctx, dto.BusinessID)
		if err != nil {
			return nil, err
		}
		stillWithinCycle := meta != nil && meta.Status == entities.BusinessStatusActive
		if current != nil && stillWithinCycle && current.SubscriptionTypeID == subType.ID && time.Now().Before(current.EndDate) {
			return nil, errs.ErrCannotRenewBeforeCycleEnds
		}
	}

	amount := subType.Price * float64(dto.Months)
	start, endDate, err := uc.computeSubscriptionWindow(ctx, dto.BusinessID, dto.Months, dto.StartDate, subType)
	if err != nil {
		return nil, err
	}

	paymentMethod := entities.PaymentMethodManual
	if dto.PaymentMethod != nil && *dto.PaymentMethod != "" {
		paymentMethod = *dto.PaymentMethod
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
		PaymentMethod:        paymentMethod,
		PaymentReference:     dto.PaymentReference,
		Notes:                dto.Notes,
	}

	if err := uc.repo.CreateSubscriptionAndActivate(ctx, sub, subType.ID, endDate); err != nil {
		return nil, err
	}

	if paymentMethod == entities.PaymentMethodWallet {
		reference := fmt.Sprintf("SUB-%d-%s-%dM", dto.BusinessID, subType.Code, dto.Months)
		if err := uc.wallet.Debit(ctx, dto.BusinessID, amount, reference, walletConceptSubscription, actorUserID); err != nil {
			uc.log.Error(ctx).Err(err).
				Uint("business_id", dto.BusinessID).
				Float64("amount", amount).
				Str("reference", reference).
				Msg("subscription registered but wallet debit failed, requires manual reconciliation")
			return nil, err
		}
	}

	uc.deactivateExpiryAnnouncements(ctx, dto.BusinessID)
	uc.recordAudit(ctx, dto.BusinessID, actorUserID, entities.AuditActionPaymentRegistered,
		fmt.Sprintf("registro un pago de %.0f (%d meses, plan %s)", amount, dto.Months, subType.Name))

	return sub, nil
}
