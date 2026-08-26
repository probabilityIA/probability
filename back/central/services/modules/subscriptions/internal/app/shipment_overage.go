package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
)

const (
	walletConceptOverage = "SHIPMENT_OVERAGE"

	OverageReasonNotAccepted = "overage_not_accepted"
	OverageReasonPaymentDue  = "overage_payment_due"
)

func (uc *UseCase) CheckShipmentOverage(ctx context.Context, businessID uint) (blocked bool, reason string, fee float64, err error) {
	usage, err := uc.repo.GetSubscriptionUsage(ctx, businessID)
	if err != nil {
		return false, "", 0, err
	}
	if usage == nil {
		return false, "", 0, nil
	}
	if usage.PlanCode != freePlanCode {
		return false, "", 0, nil
	}

	if usage.OverageAmountDue != nil && usage.OverageAmountPaidAt == nil {
		return true, OverageReasonPaymentDue, *usage.OverageAmountDue, nil
	}

	if usage.IncludedShipments == nil || usage.ShipmentOveragePrice == nil {
		return false, "", 0, nil
	}
	if usage.ShipmentsUsed < int64(*usage.IncludedShipments) {
		return false, "", 0, nil
	}
	if usage.OverageAccepted {
		return false, "", 0, nil
	}
	return true, OverageReasonNotAccepted, *usage.ShipmentOveragePrice, nil
}

func (uc *UseCase) AcceptShipmentOverage(ctx context.Context, businessID uint) error {
	return uc.repo.AcceptOverage(ctx, businessID, time.Now())
}

func (uc *UseCase) PayShipmentOverage(ctx context.Context, businessID uint, actorUserID uint) error {
	sub, err := uc.repo.GetLatestByBusinessID(ctx, businessID)
	if err != nil {
		return err
	}
	if sub == nil || sub.OverageAmountDue == nil || sub.OverageAmountPaidAt != nil {
		return errs.ErrNoOverageDue
	}

	amount := *sub.OverageAmountDue
	reference := fmt.Sprintf("OVERAGE-%d-%d", businessID, sub.ID)
	if err := uc.wallet.Debit(ctx, businessID, amount, reference, walletConceptOverage, actorUserID); err != nil {
		return err
	}

	now := time.Now()
	if err := uc.repo.MarkOverageAmountPaid(ctx, sub.ID, now); err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("business_id", businessID).
			Float64("amount", amount).
			Str("reference", reference).
			Msg("wallet debited for overage but could not be marked as paid, requires manual reconciliation")
		return err
	}

	freePlan, err := uc.repo.GetSubscriptionTypeByCode(ctx, freePlanCode)
	if err != nil {
		return err
	}
	if freePlan == nil {
		return errs.ErrSubscriptionTypeNotFound
	}

	notes := "renovacion tras pago del excedente de envios"
	endDate := now.AddDate(0, 1, 0)
	if err := uc.activateNextFreeCycle(ctx, businessID, freePlan, now, endDate, &notes); err != nil {
		return err
	}

	uc.recordAudit(ctx, businessID, actorUserID, entities.AuditActionOverageDuePaid,
		fmt.Sprintf("pago el cargo de excedente de %.0f, ciclo renovado", amount))

	return nil
}
