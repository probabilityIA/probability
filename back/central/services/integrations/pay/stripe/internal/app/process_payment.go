package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/stripe/internal/domain/ports"
)

// PaymentRequestMsg mensaje recibido desde la cola pay.stripe.requests
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

// ProcessPayment procesa una solicitud de pago via Stripe
func (uc *useCase) ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error {
	startTime := time.Now()

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing Stripe payment")

	config, err := uc.integrationRepo.GetStripeConfig(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get Stripe config")
		return uc.publishError(ctx, msg, "config_error", err.Error(), startTime)
	}

	currency := msg.Currency
	if currency == "" {
		currency = "usd"
	}

	paymentIntentID, clientSecret, err := uc.stripeClient.CreatePaymentIntent(ctx, config, msg.Amount, currency, msg.Reference, msg.Description)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Stripe payment intent creation failed")
		return uc.publishError(ctx, msg, "api_error", err.Error(), startTime)
	}

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Str("payment_intent_id", paymentIntentID).
		Msg("Stripe payment intent created successfully")

	processingTime := time.Since(startTime).Milliseconds()

	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          "stripe",
		Status:               "success",
		ExternalID:           &paymentIntentID,
		GatewayResponse: map[string]interface{}{
			"payment_intent_id": paymentIntentID,
			"client_secret":     clientSecret,
		},
		CorrelationID:    msg.CorrelationID,
		ProcessingTimeMs: processingTime,
	})
}

func (uc *useCase) publishError(ctx context.Context, msg *PaymentRequestMsg, code, errMsg string, startTime time.Time) error {
	processingTime := time.Since(startTime).Milliseconds()
	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          "stripe",
		Status:               "error",
		Error:                errMsg,
		ErrorCode:            code,
		CorrelationID:        msg.CorrelationID,
		ProcessingTimeMs:     processingTime,
	})
}
