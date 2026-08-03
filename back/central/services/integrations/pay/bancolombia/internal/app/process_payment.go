package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/ports"
)

const gatewayCode = "bancolombia"

type PaymentRequestMsg struct {
	PaymentTransactionID uint                   `json:"payment_transaction_id"`
	BusinessID           uint                   `json:"business_id"`
	GatewayCode          string                 `json:"gateway_code"`
	Amount               float64                `json:"amount"`
	Currency             string                 `json:"currency"`
	Reference            string                 `json:"reference"`
	PaymentMethod        string                 `json:"payment_method"`
	Description          string                 `json:"description"`
	IsTest               bool                   `json:"is_test"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	CorrelationID        string                 `json:"correlation_id"`
	Timestamp            time.Time              `json:"timestamp"`
}

func (uc *useCase) ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error {
	startTime := time.Now()

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Bool("is_test", msg.IsTest).
		Msg("Processing Bancolombia QR payment")

	config, err := uc.integrationRepo.GetConfigForMode(ctx, msg.IsTest)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get Bancolombia config")
		return uc.publishError(ctx, msg, "config_error", err.Error(), startTime)
	}

	currency := msg.Currency
	if currency == "" {
		currency = "COP"
	}

	result, err := uc.client.GenerateQR(ctx, config, msg.Amount, currency, msg.Reference, msg.Description)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Bancolombia QR generation failed")
		return uc.publishError(ctx, msg, "api_error", err.Error(), startTime)
	}

	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Str("external_id", result.ExternalID).
		Msg("Bancolombia QR generated successfully")

	gatewayResponse := map[string]interface{}{
		"qr_value":        result.QRValue,
		"qr_image_base64": result.QRImageBase64,
		"environment":     config.Environment,
		"status":          "PENDING",
	}
	if result.RawResponse != nil {
		gatewayResponse["raw"] = result.RawResponse
	}

	externalID := result.ExternalID
	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          gatewayCode,
		Status:               "success",
		ExternalID:           &externalID,
		GatewayResponse:      gatewayResponse,
		CorrelationID:        msg.CorrelationID,
		ProcessingTimeMs:     time.Since(startTime).Milliseconds(),
	})
}

func (uc *useCase) publishError(ctx context.Context, msg *PaymentRequestMsg, code, errMsg string, startTime time.Time) error {
	return uc.responsePublisher.PublishPaymentResponse(ctx, &ports.PaymentResponseMsg{
		PaymentTransactionID: msg.PaymentTransactionID,
		GatewayCode:          gatewayCode,
		Status:               "error",
		Error:                errMsg,
		ErrorCode:            code,
		CorrelationID:        msg.CorrelationID,
		ProcessingTimeMs:     time.Since(startTime).Milliseconds(),
	})
}
