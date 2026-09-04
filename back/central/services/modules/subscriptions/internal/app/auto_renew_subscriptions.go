package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

const autoRenewMonths = 1

func (uc *UseCase) autoRenewIfEnabled(ctx context.Context, business entities.ExpiringBusiness) bool {
	if !business.AutoPaymentEnabled || business.SubscriptionTypeID == 0 {
		return false
	}

	subType, err := uc.repo.GetSubscriptionType(ctx, business.SubscriptionTypeID)
	if err != nil || subType == nil || !subType.Active || !subType.Payable {
		return false
	}
	if subType.Code == freePlanCode {
		return false
	}

	current, err := uc.repo.GetLatestByBusinessID(ctx, business.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("auto renew: failed to load current subscription")
		return false
	}

	overage, err := uc.computeCurrentCycleOverage(ctx, business.BusinessID, current)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("auto renew: failed to compute overage")
		return false
	}

	amount := subType.Price*float64(autoRenewMonths) + overage

	balance, err := uc.wallet.GetBalance(ctx, business.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("auto renew: failed to read wallet balance")
		return false
	}
	if balance < amount {
		return false
	}

	start, endDate, err := uc.computeSubscriptionWindowFrom(current, autoRenewMonths, nil, subType, true)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("auto renew: failed to compute window")
		return false
	}

	systemUserID, err := uc.resolveSystemUserID(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("auto renew: failed to resolve system actor")
		return false
	}

	reference := fmt.Sprintf("SUB-AUTO-%d-%s-%dM", business.BusinessID, subType.Code, autoRenewMonths)
	if err := uc.wallet.Debit(ctx, business.BusinessID, amount, reference, walletConceptSubscription, systemUserID); err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("auto renew: wallet debit failed")
		return false
	}

	sub := &entities.BusinessSubscription{
		BusinessID:           business.BusinessID,
		SubscriptionTypeID:   subType.ID,
		SubscriptionTypeName: subType.Name,
		Months:               autoRenewMonths,
		Amount:               amount,
		StartDate:            start,
		EndDate:              endDate,
		Status:               entities.SubscriptionStatusPaid,
		PaymentMethod:        entities.PaymentMethodWallet,
	}
	if err := uc.repo.CreateSubscriptionAndActivate(ctx, sub, subType.ID, endDate); err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("business_id", business.BusinessID).
			Float64("amount", amount).
			Str("reference", reference).
			Msg("auto renew: wallet debited but subscription could not be activated, requires manual reconciliation")
		return false
	}

	uc.deactivateExpiryAnnouncements(ctx, business.BusinessID)
	uc.recordAudit(ctx, business.BusinessID, systemUserID, entities.AuditActionAutoRenewed,
		fmt.Sprintf("pago automatico de la suscripcion por %.0f (%d mes, plan %s)", amount, autoRenewMonths, subType.Name))

	return true
}
