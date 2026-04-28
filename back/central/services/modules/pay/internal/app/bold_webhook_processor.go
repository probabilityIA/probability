package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	eventCategoryPay         = "pay"
	eventWalletRechargeOK    = "wallet.recharge.completed"
	eventWalletRechargeFail  = "wallet.recharge.failed"
)

const boldGatewayCode = "bold"
const walletReferencePrefix = "WLT"
const walletSandboxReferencePrefix = "BOLD_SANDBOX_WLT"

func (uc *useCase) ProcessBoldWebhookMessage(ctx context.Context, msg *dtos.BoldWebhookMessage) error {
	if msg == nil || msg.BoldEventID == "" {
		return fmt.Errorf("bold webhook message missing event id")
	}

	rawPayload := msg.RawPayload
	if len(rawPayload) == 0 {
		buf, _ := json.Marshal(msg)
		rawPayload = buf
	}

	event := &dtos.BoldWebhookEvent{
		BoldEventID:    msg.BoldEventID,
		Type:           msg.Type,
		Subject:        msg.Subject,
		Source:         msg.Source,
		OccurredAt:     msg.OccurredAt,
		Payload:        rawPayload,
		SignatureValid: true,
	}

	created, err := uc.repo.RecordBoldWebhookEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("record bold webhook event: %w", err)
	}
	if !created {
		uc.log.Info(ctx).
			Str("bold_event_id", msg.BoldEventID).
			Str("type", msg.Type).
			Msg("bold webhook: duplicate event ignored (idempotent)")
		return nil
	}

	if isWalletRechargeReference(msg.MerchantReference) {
		if procErr := uc.processWalletRechargeWebhook(ctx, event, msg, rawPayload); procErr != nil {
			_ = uc.repo.MarkBoldWebhookProcessed(ctx, event.ID, nil, procErr)
			return procErr
		}
		_ = uc.repo.MarkBoldWebhookProcessed(ctx, event.ID, nil, nil)
		return nil
	}

	tx, lookupErr := uc.findBoldPaymentTransaction(ctx, msg)
	if lookupErr != nil {
		uc.log.Warn(ctx).
			Err(lookupErr).
			Str("bold_event_id", msg.BoldEventID).
			Str("merchant_reference", msg.MerchantReference).
			Str("payment_id", msg.PaymentID).
			Msg("bold webhook: payment_transaction not found")
		_ = uc.repo.MarkBoldWebhookProcessed(ctx, event.ID, nil, lookupErr)
		return nil
	}

	newStatus := mapBoldEventToStatus(msg.Type)
	if newStatus == "" {
		uc.log.Warn(ctx).
			Str("bold_event_id", msg.BoldEventID).
			Str("type", msg.Type).
			Msg("bold webhook: unknown event type")
		_ = uc.repo.MarkBoldWebhookProcessed(ctx, event.ID, &tx.ID, fmt.Errorf("unknown event type %s", msg.Type))
		return nil
	}

	if string(tx.Status) == newStatus {
		uc.log.Info(ctx).
			Uint("transaction_id", tx.ID).
			Str("status", newStatus).
			Msg("bold webhook: status unchanged, skipping update")
		_ = uc.repo.MarkBoldWebhookProcessed(ctx, event.ID, &tx.ID, nil)
		return nil
	}

	tx.Status = entities.PaymentStatus(newStatus)
	if msg.PaymentID != "" {
		ext := msg.PaymentID
		tx.ExternalID = &ext
	}
	if err := uc.repo.UpdatePaymentTransaction(ctx, tx); err != nil {
		_ = uc.repo.MarkBoldWebhookProcessed(ctx, event.ID, &tx.ID, err)
		return fmt.Errorf("update payment_transaction: %w", err)
	}

	switch newStatus {
	case constants.StatusCompleted:
		if uc.ssePublisher != nil {
			_ = uc.ssePublisher.PublishPaymentCompleted(ctx, tx)
		}
	case constants.StatusFailed:
		if uc.ssePublisher != nil {
			_ = uc.ssePublisher.PublishPaymentFailed(ctx, tx, fmt.Sprintf("bold event %s", msg.Type))
		}
	}

	if err := uc.repo.MarkBoldWebhookProcessed(ctx, event.ID, &tx.ID, nil); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("bold webhook: mark processed failed")
	}

	uc.log.Info(ctx).
		Uint("transaction_id", tx.ID).
		Str("type", msg.Type).
		Str("new_status", newStatus).
		Msg("bold webhook processed")

	return nil
}

func (uc *useCase) findBoldPaymentTransaction(ctx context.Context, msg *dtos.BoldWebhookMessage) (*entities.PaymentTransaction, error) {
	if msg.MerchantReference != "" {
		if tx, err := uc.repo.GetPaymentTransactionByReference(ctx, msg.MerchantReference); err == nil && tx != nil {
			return tx, nil
		}
	}
	return nil, fmt.Errorf("payment_transaction not found for reference=%s payment_id=%s", msg.MerchantReference, msg.PaymentID)
}

func isWalletRechargeReference(ref string) bool {
	if ref == "" {
		return false
	}
	return strings.HasPrefix(ref, walletReferencePrefix) || strings.HasPrefix(ref, walletSandboxReferencePrefix)
}

func (uc *useCase) processWalletRechargeWebhook(ctx context.Context, event *dtos.BoldWebhookEvent, msg *dtos.BoldWebhookMessage, rawPayload []byte) error {
	walletTx, err := uc.repo.GetWalletTransactionByReference(ctx, msg.MerchantReference)
	if err != nil {
		return fmt.Errorf("lookup wallet transaction by reference: %w", err)
	}
	if walletTx == nil {
		uc.log.Warn(ctx).
			Str("bold_event_id", msg.BoldEventID).
			Str("merchant_reference", msg.MerchantReference).
			Msg("bold webhook: wallet transaction not found, ignoring")
		return nil
	}

	newStatus := mapBoldEventToWalletStatus(msg.Type)
	if newStatus == "" {
		uc.log.Warn(ctx).
			Str("bold_event_id", msg.BoldEventID).
			Str("type", msg.Type).
			Msg("bold webhook: unknown event type for wallet recharge")
		return nil
	}

	if walletTx.Status != entities.WalletTxStatusPending {
		uc.log.Info(ctx).
			Str("wallet_tx_id", walletTx.ID.String()).
			Str("current_status", walletTx.Status).
			Str("new_status", newStatus).
			Msg("bold webhook: wallet transaction not pending, skipping")
		return nil
	}

	walletTx.Status = newStatus
	if err := uc.repo.UpdateWalletTransaction(ctx, walletTx); err != nil {
		return fmt.Errorf("update wallet transaction: %w", err)
	}

	if len(rawPayload) > 0 {
		if err := uc.repo.SaveWalletTransactionGatewayResponse(ctx, walletTx.ID, rawPayload); err != nil {
			uc.log.Warn(ctx).Err(err).Str("wallet_tx_id", walletTx.ID.String()).Msg("bold webhook: failed to save gateway_response")
		}
	}
	if event != nil && event.ID != uuid.Nil {
		if err := uc.repo.LinkBoldWebhookToWalletTransaction(ctx, event.ID, walletTx.ID); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("bold webhook: failed to link webhook event to wallet tx")
		}
	}

	if newStatus == entities.WalletTxStatusCompleted {
		wallet, err := uc.repo.GetWalletByID(ctx, walletTx.WalletID)
		if err != nil {
			return fmt.Errorf("get wallet: %w", err)
		}
		wallet.Balance += walletTx.Amount
		if err := uc.repo.UpdateWallet(ctx, wallet); err != nil {
			return fmt.Errorf("update wallet balance: %w", err)
		}
		uc.log.Info(ctx).
			Str("wallet_tx_id", walletTx.ID.String()).
			Str("merchant_reference", msg.MerchantReference).
			Float64("amount", walletTx.Amount).
			Float64("new_balance", wallet.Balance).
			Msg("bold webhook: wallet recharge approved and credited")
		uc.publishWalletRechargeEvent(ctx, eventWalletRechargeOK, wallet.BusinessID, walletTx, msg, &wallet.Balance, "")
		return nil
	}

	uc.log.Info(ctx).
		Str("wallet_tx_id", walletTx.ID.String()).
		Str("status", newStatus).
		Msg("bold webhook: wallet recharge marked failed")
	wallet, _ := uc.repo.GetWalletByID(ctx, walletTx.WalletID)
	var businessID uint
	var balancePtr *float64
	if wallet != nil {
		businessID = wallet.BusinessID
		balancePtr = &wallet.Balance
	}
	uc.publishWalletRechargeEvent(ctx, eventWalletRechargeFail, businessID, walletTx, msg, balancePtr, msg.Type)
	return nil
}

func (uc *useCase) publishWalletRechargeEvent(
	ctx context.Context,
	eventType string,
	businessID uint,
	walletTx *entities.WalletTransaction,
	msg *dtos.BoldWebhookMessage,
	newBalance *float64,
	reason string,
) {
	if uc.queue == nil {
		return
	}
	data := map[string]interface{}{
		"order_id":              msg.MerchantReference,
		"wallet_transaction_id": walletTx.ID.String(),
		"amount":                walletTx.Amount,
		"gateway":               boldGatewayCode,
		"bold_event_id":         msg.BoldEventID,
	}
	if newBalance != nil {
		data["new_balance"] = *newBalance
	}
	if reason != "" {
		data["reason"] = reason
	}
	envelope := rabbitmq.EventEnvelope{
		Type:       eventType,
		Category:   eventCategoryPay,
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data:       data,
	}
	if err := rabbitmq.PublishEvent(ctx, uc.queue, envelope); err != nil {
		uc.log.Warn(ctx).Err(err).Str("event_type", eventType).Msg("bold webhook: failed to publish wallet recharge event")
	}
}

func mapBoldEventToWalletStatus(eventType string) string {
	switch strings.ToUpper(eventType) {
	case "SALE_APPROVED":
		return entities.WalletTxStatusCompleted
	case "SALE_REJECTED", "VOID_APPROVED":
		return entities.WalletTxStatusFailed
	default:
		return ""
	}
}

func mapBoldEventToStatus(eventType string) string {
	switch strings.ToUpper(eventType) {
	case "SALE_APPROVED":
		return constants.StatusCompleted
	case "SALE_REJECTED":
		return constants.StatusFailed
	case "VOID_APPROVED":
		return constants.StatusCancelled
	case "VOID_REJECTED":
		return ""
	default:
		return ""
	}
}
