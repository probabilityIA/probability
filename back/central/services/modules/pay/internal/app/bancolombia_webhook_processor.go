package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

func (uc *useCase) ProcessBancolombiaWebhookMessage(ctx context.Context, msg *dtos.BancolombiaWebhookMessage) error {
	if msg == nil || msg.EventID == "" {
		return fmt.Errorf("bancolombia webhook message missing event id")
	}

	rawPayload := msg.RawPayload
	if len(rawPayload) == 0 {
		buf, _ := json.Marshal(msg)
		rawPayload = buf
	}

	event := &dtos.PaymentWebhookEvent{
		Gateway:        dtos.GatewayBancolombia,
		EventID:        msg.EventID,
		Type:           msg.Type,
		ExternalID:     msg.TransferCode,
		ExternalState:  msg.TransferState,
		Reference:      msg.Reference,
		OccurredAt:     msg.OccurredAt,
		Payload:        rawPayload,
		SignatureValid: msg.SignatureValid,
	}

	created, err := uc.repo.RecordPaymentWebhookEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("record bancolombia webhook event: %w", err)
	}
	if !created {
		uc.log.Info(ctx).
			Str("event_id", msg.EventID).
			Str("type", msg.Type).
			Msg("bancolombia webhook: duplicate event ignored (idempotent)")
		return nil
	}

	if msg.Reference == "" {
		lookupErr := fmt.Errorf("bancolombia webhook without reference")
		uc.log.Warn(ctx).
			Str("event_id", msg.EventID).
			Msg("bancolombia webhook: missing reference, cannot match payment")
		_ = uc.repo.MarkPaymentWebhookProcessed(ctx, event.ID, nil, lookupErr)
		return nil
	}

	tx, lookupErr := uc.repo.GetPaymentTransactionByReference(ctx, msg.Reference)
	if lookupErr != nil || tx == nil {
		if lookupErr == nil {
			lookupErr = fmt.Errorf("payment_transaction not found for reference=%s", msg.Reference)
		}
		uc.log.Warn(ctx).
			Err(lookupErr).
			Str("event_id", msg.EventID).
			Str("reference", msg.Reference).
			Msg("bancolombia webhook: payment_transaction not found")
		_ = uc.repo.MarkPaymentWebhookProcessed(ctx, event.ID, nil, lookupErr)
		return nil
	}

	newStatus := mapBancolombiaStateToStatus(msg.TransferState, msg.Type)
	if newStatus == "" {
		uc.log.Warn(ctx).
			Str("event_id", msg.EventID).
			Str("type", msg.Type).
			Str("state", msg.TransferState).
			Msg("bancolombia webhook: unknown transfer state")
		_ = uc.repo.MarkPaymentWebhookProcessed(ctx, event.ID, &tx.ID, fmt.Errorf("unknown transfer state %s type %s", msg.TransferState, msg.Type))
		return nil
	}

	if string(tx.Status) == newStatus {
		uc.log.Info(ctx).
			Uint("transaction_id", tx.ID).
			Str("status", newStatus).
			Msg("bancolombia webhook: status unchanged, skipping update")
		_ = uc.repo.MarkPaymentWebhookProcessed(ctx, event.ID, &tx.ID, nil)
		return nil
	}

	tx.Status = entities.PaymentStatus(newStatus)
	if msg.TransferCode != "" {
		ext := msg.TransferCode
		tx.ExternalID = &ext
	}
	if err := uc.repo.UpdatePaymentTransaction(ctx, tx); err != nil {
		_ = uc.repo.MarkPaymentWebhookProcessed(ctx, event.ID, &tx.ID, err)
		return fmt.Errorf("update payment_transaction: %w", err)
	}

	switch newStatus {
	case constants.StatusCompleted:
		if uc.ssePublisher != nil {
			_ = uc.ssePublisher.PublishPaymentCompleted(ctx, tx)
		}
	case constants.StatusFailed:
		if uc.ssePublisher != nil {
			_ = uc.ssePublisher.PublishPaymentFailed(ctx, tx, fmt.Sprintf("bancolombia state %s", msg.TransferState))
		}
	}

	if err := uc.repo.MarkPaymentWebhookProcessed(ctx, event.ID, &tx.ID, nil); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("bancolombia webhook: mark processed failed")
	}

	uc.log.Info(ctx).
		Uint("transaction_id", tx.ID).
		Str("state", msg.TransferState).
		Str("new_status", newStatus).
		Msg("bancolombia webhook processed")

	return nil
}

func mapBancolombiaStateToStatus(state, eventType string) string {
	normalized := strings.ToUpper(strings.TrimSpace(state))
	if normalized == "" {
		normalized = strings.ToUpper(strings.TrimSpace(eventType))
	}
	switch {
	case containsAny(normalized, "APPROVED", "SUCCESS", "COMPLETED", "SETTLED", "OK"):
		return constants.StatusCompleted
	case containsAny(normalized, "REJECTED", "FAILED", "DECLINED", "ERROR"):
		return constants.StatusFailed
	case containsAny(normalized, "CANCELLED", "CANCELED", "VOIDED", "EXPIRED"):
		return constants.StatusCancelled
	default:
		return ""
	}
}

func containsAny(s string, fragments ...string) bool {
	for _, f := range fragments {
		if strings.Contains(s, f) {
			return true
		}
	}
	return false
}
