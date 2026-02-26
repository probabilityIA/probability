package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/melipago/internal/domain/ports"
)

// PaymentRequestMsg mensaje recibido desde la cola pay.melipago.requests
type PaymentRequestMsg struct {
	PaymentTransactionID uint                   `json:"payment_transaction_id"`
	BusinessID           uint                   `json:"business_id"`
	GatewayCode          string                 `json:"gateway_code"`
	Amount               float64                `json:"amount"`
	Currency             string                 `json:"currency"`
	Reference            string                 `json:"reference"`
	PaymentMethod        string                 `json:"payment_method"`
	Description          string                 `json:"description"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	CorrelationID        string                 `json:"correlation_id"`
	Timestamp            time.Time              `json:"timestamp"`
}

// ProcessPayment procesa una solicitud de pago via MercadoPago
func (uc *useCase) ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error {
	startTime := time.Now()

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing MercadoPago payment")

	config, err := uc.integrationRepo.GetMeliPagoConfig(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get MercadoPago config")
		return uc.publishError(ctx, msg, "config_error", err.Error(), startTime)
	}

	currency := msg.Currency
	if currency == "" {
		currency = "COP"
	}

	preferenceID, checkoutURL, err := uc.meliPagoClient.CreatePreference(ctx, config, msg.Amount, currency, msg.Reference, msg.Description)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("MercadoPago preference creation failed")
		return uc.publishError(ctx, msg, "api_error", err.Error(), startTime)
	}

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Str("preference_id", preferenceID).
		Msg("MercadoPago preference created successfully")

	processingTime := time.Since(startTime).Milliseconds()

	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          "melipago",
		Status:               "success",
		ExternalID:           &preferenceID,
		GatewayResponse: map[string]interface{}{
			"preference_id": preferenceID,
			"checkout_url":  checkoutURL,
		},
		CorrelationID:    msg.CorrelationID,
		ProcessingTimeMs: processingTime,
	})
}

func (uc *useCase) publishError(ctx context.Context, msg *PaymentRequestMsg, code, errMsg string, startTime time.Time) error {
	processingTime := time.Since(startTime).Milliseconds()
	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          "melipago",
		Status:               "error",
		Error:                errMsg,
		ErrorCode:            code,
		CorrelationID:        msg.CorrelationID,
		ProcessingTimeMs:     processingTime,
	})
}
