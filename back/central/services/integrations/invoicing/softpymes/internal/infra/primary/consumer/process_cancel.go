package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
)

// processCancelInvoice anula una factura emitida en Softpymes
func (c *InvoiceRequestConsumer) processCancelInvoice(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	// 1. Extraer external_id y reason del Config
	documentNumber, _ := request.InvoiceData.Config["external_id"].(string)
	if documentNumber == "" {
		c.log.Error(ctx).Uint("invoice_id", request.InvoiceID).Msg("external_id missing in cancel config")
		return c.createCancelErrorResponse(request, "missing_external_id", "external_id is required for cancellation", startTime)
	}

	reason, _ := request.InvoiceData.Config["cancel_reason"].(string)
	if reason == "" {
		reason = "Anulación de factura"
	}

	// 2. Obtener integración y credenciales
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		c.log.Error(ctx).Msg("integration_id is 0 in cancel request")
		return c.createCancelErrorResponse(request, "missing_integration_id", "integration_id is 0", startTime)
	}

	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get integration for cancel")
		return c.createCancelErrorResponse(request, "integration_not_found", err.Error(), startTime)
	}

	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_key")
		return c.createCancelErrorResponse(request, "decryption_failed", "Failed to decrypt api_key", startTime)
	}

	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_secret")
		return c.createCancelErrorResponse(request, "decryption_failed", "Failed to decrypt api_secret", startTime)
	}

	// 3. Combinar config
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	referer, _ := combinedConfig["referer"].(string)

	// 4. Resolver URL efectiva
	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		c.log.Error(ctx).Uint("integration_id", integrationID).Msg("base_url no configurada para cancelación")
		return c.createCancelErrorResponse(request, "missing_base_url", "base_url no configurada en el tipo de integración Softpymes", startTime)
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("document_number", documentNumber).
		Str("effective_url", effectiveURL).
		Msg("Cancelling invoice in Softpymes")

	// 5. Llamar al cliente Softpymes
	if err := c.softpymesClient.CancelInvoice(ctx, apiKey, apiSecret, referer, documentNumber, reason, effectiveURL); err != nil {
		c.log.Error(ctx).Err(err).Uint("invoice_id", request.InvoiceID).Msg("Softpymes cancel failed")
		return c.createCancelErrorResponse(request, "api_error", err.Error(), startTime)
	}

	processingTime := time.Since(startTime).Milliseconds()
	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("document_number", documentNumber).
		Msg("Invoice cancelled successfully in Softpymes")

	return &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "success",
		Operation:      "cancel",
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}
}

// createCancelErrorResponse crea una respuesta de error para operaciones de cancelación
func (c *InvoiceRequestConsumer) createCancelErrorResponse(
	request *InvoiceRequestMessage,
	errorCode string,
	errorMsg string,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	return &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "error",
		Operation:      "cancel",
		Error:          errorMsg,
		ErrorCode:      errorCode,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: time.Since(startTime).Milliseconds(),
	}
}
