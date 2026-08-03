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
	eventCategoryPay        = "pay"
	eventWalletRechargeOK   = "wallet.recharge.completed"
	eventWalletRechargeFail = "wallet.recharge.failed"
)

const boldGatewayCode = "bold"
const walletReferencePrefix = "WLT"
const walletSandboxReferencePrefix = "BOLD_SANDBOX_WLT"
const storefrontCheckoutReferencePrefix = "SFO"

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

	if isStorefrontCheckoutReference(msg.MerchantReference) {
		if procErr := uc.processStorefrontCheckoutWebhook(ctx, msg); procErr != nil {
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

func isStorefrontCheckoutReference(ref string) bool {
	return strings.HasPrefix(ref, storefrontCheckoutReferencePrefix)
}

// processStorefrontCheckoutWebhook confirma (o rechaza) un checkout de la tienda publica
// y, si Bold aprueba el pago, publica la orden canonica a probability.orders.canonical
// — el mismo pipeline que ya usan Shopify/WooCommerce.
func (uc *useCase) processStorefrontCheckoutWebhook(ctx context.Context, msg *dtos.BoldWebhookMessage) error {
	outcome := mapBoldEventToOutcome(msg.Type)
	if outcome == "" {
		uc.log.Warn(ctx).
			Str("bold_event_id", msg.BoldEventID).
			Str("type", msg.Type).
			Msg("bold webhook: unknown event type for storefront checkout")
		return nil
	}

	checkout, err := uc.repo.GetStorefrontCheckoutByReference(ctx, msg.MerchantReference)
	if err != nil {
		return fmt.Errorf("lookup storefront checkout: %w", err)
	}
	if checkout == nil {
		uc.log.Warn(ctx).Str("reference", msg.MerchantReference).Msg("bold webhook: storefront checkout not found, ignoring")
		return nil
	}
	if checkout.Status != "pending" {
		uc.log.Info(ctx).
			Str("reference", checkout.Reference).
			Str("current_status", checkout.Status).
			Msg("bold webhook: storefront checkout not pending, skipping")
		return nil
	}

	if outcome == dtos.WalletRechargeOutcomeRejected {
		if uerr := uc.repo.MarkStorefrontCheckoutStatus(ctx, checkout.Reference, "failed"); uerr != nil {
			return fmt.Errorf("mark storefront checkout failed: %w", uerr)
		}
		uc.log.Info(ctx).Str("reference", checkout.Reference).Msg("storefront checkout: payment rejected")
		return nil
	}

	if err := uc.repo.MarkStorefrontCheckoutStatus(ctx, checkout.Reference, "paid"); err != nil {
		return fmt.Errorf("mark storefront checkout paid: %w", err)
	}

	if err := uc.publishStorefrontOrder(ctx, checkout, msg); err != nil {
		uc.log.Error(ctx).Err(err).Str("reference", checkout.Reference).Msg("failed to publish storefront canonical order")
		return err
	}

	uc.log.Info(ctx).
		Str("reference", checkout.Reference).
		Uint("business_id", checkout.BusinessID).
		Float64("amount", checkout.Amount).
		Msg("storefront checkout paid, order published")
	return nil
}

func (uc *useCase) PublishAgreedStorefrontOrder(ctx context.Context, reference string) error {
	checkout, err := uc.repo.GetStorefrontCheckoutByReference(ctx, reference)
	if err != nil {
		return fmt.Errorf("lookup storefront checkout: %w", err)
	}
	if checkout == nil {
		return fmt.Errorf("checkout no encontrado: %s", reference)
	}
	return uc.publishStorefrontOrder(ctx, checkout, nil)
}

func (uc *useCase) publishStorefrontOrder(ctx context.Context, checkout *entities.StorefrontCheckoutSnapshot, msg *dtos.BoldWebhookMessage) error {
	if uc.queue == nil {
		return fmt.Errorf("cola rabbitmq no disponible")
	}

	var items []entities.StorefrontCheckoutItem
	_ = json.Unmarshal(checkout.ItemsJSON, &items)

	orderItems := make([]map[string]interface{}, 0, len(items))
	for _, it := range items {
		orderItems = append(orderItems, map[string]interface{}{
			"product_id":   it.ProductID,
			"product_sku":  it.SKU,
			"product_name": it.Name,
			"quantity":     it.Quantity,
			"unit_price":   it.UnitPrice,
			"total_price":  it.UnitPrice * float64(it.Quantity),
			"currency":     checkout.Currency,
			"image_url":    it.ImageURL,
		})
	}

	var addresses []map[string]interface{}
	if len(checkout.ShippingAddressJSON) > 0 {
		var addr entities.StorefrontCheckoutAddress
		if err := json.Unmarshal(checkout.ShippingAddressJSON, &addr); err == nil {
			addresses = append(addresses, map[string]interface{}{
				"type":         "shipping",
				"first_name":   checkout.CustomerName,
				"phone":        checkout.CustomerPhone,
				"street":       addr.Street,
				"city":         addr.City,
				"state":        addr.State,
				"country":      addr.Country,
				"postal_code":  addr.PostalCode,
				"instructions": addr.Instructions,
			})
		}
	}

	now := time.Now()
	payments := []map[string]interface{}{}
	status := "pending"
	originalStatus := "agreed"
	if msg != nil {
		paidAt := now.Format(time.RFC3339)
		payments = append(payments, map[string]interface{}{
			"payment_method_id": 1,
			"amount":            checkout.Amount,
			"currency":          checkout.Currency,
			"status":            "completed",
			"paid_at":           paidAt,
			"transaction_id":    msg.PaymentID,
			"payment_reference": checkout.Reference,
			"gateway":           boldGatewayCode,
		})
		status = "processing"
		originalStatus = "paid"
	}

	canonicalOrder := map[string]interface{}{
		"business_id":      checkout.BusinessID,
		"integration_id":   checkout.IntegrationID,
		"integration_type": "platform",
		"platform":         "tienda_web",
		"external_id":      "sfo-" + checkout.Reference,
		"order_number":     checkout.Reference,
		"subtotal":         checkout.Amount,
		"tax":              0,
		"discount":         0,
		"shipping_cost":    0,
		"total_amount":     checkout.Amount,
		"currency":         checkout.Currency,
		"is_cod":           false,
		"customer_name":    checkout.CustomerName,
		"customer_email":   checkout.CustomerEmail,
		"customer_phone":   checkout.CustomerPhone,
		"customer_dni":     checkout.CustomerDni,
		"status":           status,
		"original_status":  originalStatus,
		"invoiceable":      false,
		"occurred_at":      now.Format(time.RFC3339),
		"imported_at":      now.Format(time.RFC3339),
		"order_items":      orderItems,
		"addresses":        addresses,
		"payments":         payments,
		"shipments":        []interface{}{},
	}

	orderJSON, err := json.Marshal(canonicalOrder)
	if err != nil {
		return fmt.Errorf("marshal storefront canonical order: %w", err)
	}

	if err := uc.queue.Publish(ctx, rabbitmq.QueueOrdersCanonical, orderJSON); err != nil {
		return fmt.Errorf("publish storefront canonical order: %w", err)
	}
	return nil
}

func (uc *useCase) processWalletRechargeWebhook(ctx context.Context, event *dtos.BoldWebhookEvent, msg *dtos.BoldWebhookMessage, rawPayload []byte) error {
	outcome := mapBoldEventToOutcome(msg.Type)
	if outcome == "" {
		uc.log.Warn(ctx).
			Str("bold_event_id", msg.BoldEventID).
			Str("type", msg.Type).
			Msg("bold webhook: unknown event type for wallet recharge")
		return nil
	}

	in := &dtos.WalletRechargeStatusInput{
		OrderID:         msg.MerchantReference,
		Outcome:         outcome,
		Source:          "webhook",
		BoldEventID:     msg.BoldEventID,
		GatewayResponse: rawPayload,
		Reason:          msg.Type,
	}
	if err := uc.ApplyWalletRechargeStatus(ctx, in); err != nil {
		return err
	}
	if event != nil && event.ID != uuid.Nil {
		walletTx, lookupErr := uc.repo.GetWalletTransactionByReference(ctx, msg.MerchantReference)
		if lookupErr == nil && walletTx != nil {
			if linkErr := uc.repo.LinkBoldWebhookToWalletTransaction(ctx, event.ID, walletTx.ID); linkErr != nil {
				uc.log.Warn(ctx).Err(linkErr).Msg("bold webhook: failed to link webhook event to wallet tx")
			}
		}
	}
	return nil
}

func (uc *useCase) ApplyWalletRechargeStatus(ctx context.Context, in *dtos.WalletRechargeStatusInput) error {
	if in == nil || in.OrderID == "" {
		return fmt.Errorf("apply wallet recharge: missing order id")
	}

	walletTx, err := uc.repo.GetWalletTransactionByReference(ctx, in.OrderID)
	if err != nil {
		return fmt.Errorf("lookup wallet transaction by reference: %w", err)
	}
	if walletTx == nil {
		uc.log.Warn(ctx).
			Str("source", in.Source).
			Str("order_id", in.OrderID).
			Msg("wallet recharge: transaction not found, ignoring")
		return nil
	}

	var newStatus string
	switch in.Outcome {
	case dtos.WalletRechargeOutcomeApproved:
		newStatus = entities.WalletTxStatusCompleted
	case dtos.WalletRechargeOutcomeRejected:
		newStatus = entities.WalletTxStatusFailed
	default:
		uc.log.Warn(ctx).
			Str("source", in.Source).
			Str("outcome", in.Outcome).
			Msg("wallet recharge: unknown outcome, ignoring")
		return nil
	}

	if walletTx.Status == entities.WalletTxStatusCompleted {
		uc.log.Info(ctx).
			Str("source", in.Source).
			Str("wallet_tx_id", walletTx.ID.String()).
			Str("current_status", walletTx.Status).
			Str("target_status", newStatus).
			Msg("wallet recharge: already completed, skipping")
		return nil
	}

	walletTx.Status = newStatus
	if err := uc.repo.UpdateWalletTransaction(ctx, walletTx); err != nil {
		return fmt.Errorf("update wallet transaction: %w", err)
	}

	if len(in.GatewayResponse) > 0 {
		if saveErr := uc.repo.SaveWalletTransactionGatewayResponse(ctx, walletTx.ID, in.GatewayResponse); saveErr != nil {
			uc.log.Warn(ctx).Err(saveErr).Str("wallet_tx_id", walletTx.ID.String()).Msg("wallet recharge: failed to save gateway_response")
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
			Str("source", in.Source).
			Str("wallet_tx_id", walletTx.ID.String()).
			Str("order_id", in.OrderID).
			Float64("amount", walletTx.Amount).
			Float64("new_balance", wallet.Balance).
			Msg("wallet recharge approved and credited")
		uc.publishWalletRechargeEventRaw(ctx, eventWalletRechargeOK, wallet.BusinessID, walletTx, in.OrderID, in.BoldEventID, &wallet.Balance, "")
		return nil
	}

	uc.log.Info(ctx).
		Str("source", in.Source).
		Str("wallet_tx_id", walletTx.ID.String()).
		Str("status", newStatus).
		Msg("wallet recharge marked failed")
	wallet, _ := uc.repo.GetWalletByID(ctx, walletTx.WalletID)
	var businessID uint
	var balancePtr *float64
	if wallet != nil {
		businessID = wallet.BusinessID
		balancePtr = &wallet.Balance
	}
	uc.publishWalletRechargeEventRaw(ctx, eventWalletRechargeFail, businessID, walletTx, in.OrderID, in.BoldEventID, balancePtr, in.Reason)
	return nil
}

func mapBoldEventToOutcome(eventType string) string {
	switch strings.ToUpper(eventType) {
	case "SALE_APPROVED":
		return dtos.WalletRechargeOutcomeApproved
	case "SALE_REJECTED", "VOID_APPROVED":
		return dtos.WalletRechargeOutcomeRejected
	default:
		return ""
	}
}

func (uc *useCase) publishWalletRechargeEventRaw(
	ctx context.Context,
	eventType string,
	businessID uint,
	walletTx *entities.WalletTransaction,
	orderID string,
	boldEventID string,
	newBalance *float64,
	reason string,
) {
	if uc.queue == nil {
		return
	}
	data := map[string]interface{}{
		"order_id":              orderID,
		"wallet_transaction_id": walletTx.ID.String(),
		"amount":                walletTx.Amount,
		"gateway":               boldGatewayCode,
	}
	if boldEventID != "" {
		data["bold_event_id"] = boldEventID
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
		uc.log.Warn(ctx).Err(err).Str("event_type", eventType).Msg("wallet recharge: failed to publish event")
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
