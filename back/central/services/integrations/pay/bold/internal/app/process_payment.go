package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/ports"
)

// PaymentRequestMsg mensaje recibido desde la cola pay.bold.requests
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

// ProcessPayment procesa una solicitud de pago via Bold
func (uc *useCase) ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error {
	startTime := time.Now()

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing Bold payment")

	// Obtener configuración de Bold desde integration_types
	config, err := uc.integrationRepo.GetBoldConfig(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get Bold config")
		return uc.publishError(ctx, msg, "config_error", err.Error(), startTime)
	}

	currency := msg.Currency
	if currency == "" {
		currency = "COP"
	}

	// Crear link de pago en Bold
	linkID, checkoutURL, err := uc.boldClient.CreatePaymentLink(ctx, config, msg.Amount, currency, msg.Reference, msg.Description)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Bold payment link creation failed")
		return uc.publishError(ctx, msg, "api_error", err.Error(), startTime)
	}

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Str("link_id", linkID).
		Msg("Bold payment link created successfully")

	// Publicar respuesta exitosa
	processingTime := time.Since(startTime).Milliseconds()

	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          "bold",
		Status:               "success",
		ExternalID:           &linkID,
		GatewayResponse: map[string]interface{}{
			"payment_link_id": linkID,
			"checkout_url":    checkoutURL,
			"status":          "ACTIVE",
		},
		CorrelationID:    msg.CorrelationID,
		ProcessingTimeMs: processingTime,
	})
}

func (uc *useCase) publishError(ctx context.Context, msg *PaymentRequestMsg, code, errMsg string, startTime time.Time) error {
	processingTime := time.Since(startTime).Milliseconds()
	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          "bold",
		Status:               "error",
		Error:                errMsg,
		ErrorCode:            code,
		CorrelationID:        msg.CorrelationID,
		ProcessingTimeMs:     processingTime,
	})
}
