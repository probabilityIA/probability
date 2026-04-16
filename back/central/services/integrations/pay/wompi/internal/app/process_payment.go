package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/wompi/internal/domain/ports"
)

// PaymentRequestMsg mensaje recibido desde la cola pay.wompi.requests
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

// ProcessPayment procesa una solicitud de pago via Wompi
func (uc *useCase) ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error {
	startTime := time.Now()

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing Wompi payment")

	config, err := uc.integrationRepo.GetWompiConfig(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get Wompi config")
		return uc.publishError(ctx, msg, "config_error", err.Error(), startTime)
	}

	currency := msg.Currency
	if currency == "" {
		currency = "COP"
	}

	transactionID, redirectURL, err := uc.wompiClient.CreateTransaction(ctx, config, msg.Amount, currency, msg.Reference, msg.Description)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Wompi transaction creation failed")
		return uc.publishError(ctx, msg, "api_error", err.Error(), startTime)
	}

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Str("wompi_tx_id", transactionID).
		Msg("Wompi transaction created successfully")

	processingTime := time.Since(startTime).Milliseconds()

	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          "wompi",
		Status:               "success",
		ExternalID:           &transactionID,
		GatewayResponse: map[string]interface{}{
			"transaction_id": transactionID,
			"redirect_url":   redirectURL,
		},
		CorrelationID:    msg.CorrelationID,
		ProcessingTimeMs: processingTime,
	})
}

func (uc *useCase) publishError(ctx context.Context, msg *PaymentRequestMsg, code, errMsg string, startTime time.Time) error {
	processingTime := time.Since(startTime).Milliseconds()
	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          "wompi",
		Status:               "error",
		Error:                errMsg,
		ErrorCode:            code,
		CorrelationID:        msg.CorrelationID,
		ProcessingTimeMs:     processingTime,
	})
}
