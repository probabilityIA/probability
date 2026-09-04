package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

var spanishMonths = [...]string{
	"enero", "febrero", "marzo", "abril", "mayo", "junio",
	"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
}

func formatSpanishDate(t time.Time) string {
	return fmt.Sprintf("%d de %s de %d", t.Day(), spanishMonths[t.Month()-1], t.Year())
}

func (uc *UseCase) notifyPaymentWindowIfNeeded(ctx context.Context, business entities.ExpiringBusiness) {
	if uc.rabbit == nil {
		return
	}

	alreadyNotified, err := uc.repo.HasAuditLogSince(ctx, business.BusinessID, entities.AuditActionPaymentWindowNotified, business.EndDate)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("payment window alert: failed to check idempotency")
		return
	}
	if alreadyNotified {
		return
	}

	contact, err := uc.repo.GetWhatsAppContact(ctx, business.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("payment window alert: failed to resolve whatsapp contact")
		return
	}
	if contact == nil || contact.Phone == "" {
		return
	}

	var cycleAmount float64
	if business.SubscriptionTypeID != 0 {
		subType, err := uc.repo.GetSubscriptionType(ctx, business.SubscriptionTypeID)
		if err == nil && subType != nil {
			current, err := uc.repo.GetLatestByBusinessID(ctx, business.BusinessID)
			if err == nil {
				overage, err := uc.computeCurrentCycleOverage(ctx, business.BusinessID, current)
				if err == nil {
					cycleAmount = subType.Price + overage
				}
			}
		}
	}

	balance, err := uc.wallet.GetBalance(ctx, business.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("payment window alert: failed to read wallet balance")
		return
	}

	systemUserID, err := uc.resolveSystemUserID(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("payment window alert: failed to resolve system actor")
		return
	}

	payload := map[string]any{
		"business_id":    business.BusinessID,
		"business_name":  contact.BusinessName,
		"phone_number":   contact.Phone,
		"due_date":       formatSpanishDate(business.EndDate),
		"cycle_amount":   cycleAmount,
		"wallet_balance": balance,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("payment window alert: failed to serialize event")
		return
	}

	if err := uc.rabbit.Publish(ctx, rabbitmq.QueueSubscriptionPaymentWindowRequested, body); err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("payment window alert: failed to publish event")
		return
	}

	uc.recordAudit(ctx, business.BusinessID, systemUserID, entities.AuditActionPaymentWindowNotified,
		fmt.Sprintf("aviso de whatsapp enviado: periodo de pago abierto, vence %s", formatSpanishDate(business.EndDate)))
}
