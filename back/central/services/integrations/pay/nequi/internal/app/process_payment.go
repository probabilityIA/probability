package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/domain/ports"
)

// PaymentRequestMsg mensaje recibido desde la cola pay.nequi.requests
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

// ProcessPayment procesa una solicitud de pago via Nequi
func (uc *useCase) ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error {
	startTime := time.Now()

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing Nequi payment")

	// Obtener configuración de Nequi desde integration_types
	config, err := uc.integrationRepo.GetNequiConfig(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get Nequi config")
		return uc.publishError(ctx, msg, "config_error", err.Error(), startTime)
	}

	// Generar QR de Nequi
	qrValue, transactionID, err := uc.nequiClient.GenerateQR(ctx, config, msg.Amount, msg.Reference)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Nequi QR generation failed")
		return uc.publishError(ctx, msg, "api_error", err.Error(), startTime)
	}

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Str("nequi_tx_id", transactionID).
		Msg("Nequi QR generated successfully")

	// Publicar respuesta exitosa
	extID := fmt.Sprintf("nequi-%s", transactionID)
	processingTime := time.Since(startTime).Milliseconds()

	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          "nequi",
		Status:               "success",
		ExternalID:           &extID,
		GatewayResponse: map[string]interface{}{
			"qr_value":       qrValue,
			"transaction_id": transactionID,
		},
		CorrelationID:    msg.CorrelationID,
		ProcessingTimeMs: processingTime,
	})
}

func (uc *useCase) publishError(ctx context.Context, msg *PaymentRequestMsg, code, errMsg string, startTime time.Time) error {
	processingTime := time.Since(startTime).Milliseconds()
	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          "nequi",
		Status:               "error",
		Error:                errMsg,
		ErrorCode:            code,
		CorrelationID:        msg.CorrelationID,
		ProcessingTimeMs:     processingTime,
	})
}
