package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
)

// processCashReceipt genera un recibo de caja para una factura ya emitida.
// Obtiene el documento completo de Softpymes por invoice_number y envía el recibo de caja.
func (c *InvoiceRequestConsumer) processCashReceipt(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	// 1. Obtener invoice_number del config
	invoiceNumber, _ := request.InvoiceData.Config["invoice_number"].(string)
	if invoiceNumber == "" {
		return c.createErrorResponse(request, "missing_invoice_number", "invoice_number is required for cash_receipt operation", startTime, nil)
	}

	// 2. Obtener integración y credenciales
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		return c.createErrorResponse(request, "missing_integration_id", "integration_id is 0", startTime, nil)
	}

	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		return c.createErrorResponse(request, "integration_not_found", err.Error(), startTime, nil)
	}

	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt api_key", startTime, nil)
	}

	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt api_secret", startTime, nil)
	}

	// 3. Combinar config
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	// 4. Resolver URL
	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		return c.createErrorResponse(request, "missing_base_url", "base_url no configurada", startTime, nil)
	}

	referer, _ := combinedConfig["referer"].(string)

	// 5. Obtener documento completo por número de factura
	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("invoice_number", invoiceNumber).
		Msg("Fetching full document from Softpymes for cash receipt")

	fullDocument, err := c.softpymesClient.GetDocumentByNumber(ctx, apiKey, apiSecret, referer, invoiceNumber, effectiveURL)
	if err != nil {
		c.log.Error(ctx).Err(err).
			Str("invoice_number", invoiceNumber).
			Msg("Failed to get document for cash receipt")
		return c.createErrorResponse(request, "document_not_found", "Failed to get document: "+err.Error(), startTime, nil)
	}

	// 6. Enviar recibo de caja
	cashReceiptAudit := c.sendCashReceiptIfConfigured(ctx, fullDocument, combinedConfig, apiKey, apiSecret, referer, effectiveURL, request.InvoiceID)

	if cashReceiptAudit == nil {
		// Cash receipt no se envió (config disabled o documento nil — no debería pasar aquí)
		return c.createErrorResponse(request, "cash_receipt_not_sent", "Cash receipt was not sent — check configuration", startTime, nil)
	}

	// 7. Verificar si hubo error en el recibo de caja
	if cashReceiptStatus, ok := fullDocument["cash_receipt"].(map[string]interface{}); ok {
		if status, ok := cashReceiptStatus["status"].(string); ok && status == "failed" {
			errorMsg := "Cash receipt failed"
			if errDetail, ok := cashReceiptStatus["error"].(string); ok {
				errorMsg = errDetail
			}
			processingTime := time.Since(startTime).Milliseconds()
			resp := &queue.InvoiceResponseMessage{
				InvoiceID:      request.InvoiceID,
				Provider:       "softpymes",
				Status:         "error",
				Operation:      "cash_receipt",
				Error:          errorMsg,
				ErrorCode:      "cash_receipt_failed",
				CorrelationID:  request.CorrelationID,
				Timestamp:      time.Now(),
				ProcessingTime: processingTime,
			}
			resp.CashReceiptRequestURL = cashReceiptAudit.RequestURL
			resp.CashReceiptRequestPayload = cashReceiptAudit.RequestPayload
			resp.CashReceiptResponseStatus = cashReceiptAudit.ResponseStatus
			resp.CashReceiptResponseBody = cashReceiptAudit.ResponseBody
			return resp
		}
	}

	// 8. Éxito — construir response
	processingTime := time.Since(startTime).Milliseconds()
	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "success",
		Operation:      "cash_receipt",
		DocumentJSON:   fullDocument,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}

	resp.CashReceiptRequestURL = cashReceiptAudit.RequestURL
	resp.CashReceiptRequestPayload = cashReceiptAudit.RequestPayload
	resp.CashReceiptResponseStatus = cashReceiptAudit.ResponseStatus
	resp.CashReceiptResponseBody = cashReceiptAudit.ResponseBody

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("invoice_number", invoiceNumber).
		Msg("Cash receipt generated successfully")

	return resp
}
