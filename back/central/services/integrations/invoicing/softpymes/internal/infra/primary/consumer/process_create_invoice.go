package consumer

import (
	"context"
	"fmt"
	"time"

	spDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
)

var documentFetchDelays = []time.Duration{
	500 * time.Millisecond,
	time.Second,
	1500 * time.Millisecond,
	2 * time.Second,
}

func (c *InvoiceRequestConsumer) processCreateInvoice(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		c.log.Error(ctx).Msg("integration_id is 0 in invoice_data")
		return c.createErrorResponse(request, "missing_integration_id", "integration_id is 0", startTime, nil)
	}

	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get integration")
		return c.createErrorResponse(request, "integration_not_found", err.Error(), startTime, nil)
	}

	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_key")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt api_key", startTime, nil)
	}

	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_secret")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt api_secret", startTime, nil)
	}

	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}

	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		c.log.Error(ctx).
			Uint("integration_id", integrationID).
			Msg("base_url no configurada en el tipo de integracion Softpymes")
		return c.createErrorResponse(request, "missing_base_url",
			"base_url no configurada en el tipo de integracion Softpymes (integration_types.base_url)",
			startTime, nil)
	}

	c.log.Info(ctx).
		Bool("is_testing", integration.IsTesting).
		Str("effective_url", effectiveURL).
		Msg("Resolved effective Softpymes URL")

	invoiceReq := &spDtos.CreateInvoiceRequest{
		Customer: spDtos.CustomerData{
			Name:    request.InvoiceData.Customer.Name,
			Email:   request.InvoiceData.Customer.Email,
			Phone:   request.InvoiceData.Customer.Phone,
			DNI:     request.InvoiceData.Customer.DNI,
			Address: request.InvoiceData.Customer.Address,
		},
		Items:            mapItemsToClientDTOs(request.InvoiceData.Items),
		Total:            request.InvoiceData.Total,
		Subtotal:         request.InvoiceData.Subtotal,
		Tax:              request.InvoiceData.Tax,
		Discount:         request.InvoiceData.Discount,
		ShippingCost:     request.InvoiceData.ShippingCost,
		ShippingDiscount: request.InvoiceData.ShippingDiscount,
		ShippingCostBase: request.InvoiceData.ShippingCostBase,
		Currency:         request.InvoiceData.Currency,
		OrderID:          request.InvoiceData.OrderID,
		OrderNumber:      request.InvoiceData.OrderNumber,
		Credentials: spDtos.Credentials{
			APIKey:    apiKey,
			APISecret: apiSecret,
		},
		Config:  combinedConfig,
		IsRetry: request.Operation == "retry",
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("effective_url", effectiveURL).
		Msg("Calling Softpymes API")

	result, err := c.softpymesClient.CreateInvoice(ctx, invoiceReq, effectiveURL)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", request.InvoiceID).
			Msg("Softpymes API call failed")

		var auditData *spDtos.AuditData
		if result != nil {
			auditData = result.AuditData
		}
		return c.createErrorResponse(request, "api_error", err.Error(), startTime, auditData)
	}

	if result.PendingValidation {
		c.log.Info(ctx).
			Uint("invoice_id", request.InvoiceID).
			Str("message", result.ProviderMessage).
			Msg("Invoice pending DIAN validation")

		processingTime := time.Since(startTime).Milliseconds()
		resp := &queue.InvoiceResponseMessage{
			InvoiceID:      request.InvoiceID,
			Provider:       "softpymes",
			Status:         "pending_validation",
			CorrelationID:  request.CorrelationID,
			Timestamp:      time.Now(),
			ProcessingTime: processingTime,
			Error:          result.ProviderMessage,
		}
		if result.AuditData != nil {
			resp.AuditRequestURL = result.AuditData.RequestURL
			resp.AuditRequestPayload = toMapPayload(result.AuditData.RequestPayload)
			resp.AuditResponseStatus = result.AuditData.ResponseStatus
			resp.AuditResponseBody = result.AuditData.ResponseBody
		}
		return resp
	}

	var fullDocument map[string]interface{}
	if result.InvoiceNumber != "" {
		referer, _ := combinedConfig["referer"].(string)

		for attempt, delay := range documentFetchDelays {
			time.Sleep(delay)

			doc, docErr := c.softpymesClient.GetDocumentByNumber(ctx, apiKey, apiSecret, referer, result.InvoiceNumber, effectiveURL)
			if docErr == nil && doc != nil {
				fullDocument = doc
				c.log.Info(ctx).
					Str("invoice_number", result.InvoiceNumber).
					Int("attempt", attempt+1).
					Msg("Full document retrieved")
				break
			}

			c.log.Info(ctx).
				Err(docErr).
				Str("invoice_number", result.InvoiceNumber).
				Int("attempt", attempt+1).
				Msg("Document not yet available in Softpymes, retrying")
		}

		if fullDocument == nil {
			c.log.Warn(ctx).
				Str("invoice_number", result.InvoiceNumber).
				Msg("Failed to fetch full document - continuing without it")
		} else if docNum, ok := fullDocument["documentNumber"].(string); ok && docNum != "" {
			c.log.Info(ctx).
				Str("old_invoice_number", result.InvoiceNumber).
				Str("new_invoice_number", docNum).
				Msg("Overriding invoice number with canonical padded format from full document")
			result.InvoiceNumber = docNum
			result.ExternalID = docNum
		}
	}

	referer, _ := combinedConfig["referer"].(string)
	cashReceiptAudit := c.sendCashReceiptIfConfigured(ctx, fullDocument, combinedConfig, apiKey, apiSecret, referer, effectiveURL, request.InvoiceID)

	var issuedAt *time.Time
	if result.IssuedAt != "" {
		if parsed, parseErr := time.Parse(time.RFC3339, result.IssuedAt); parseErr == nil {
			issuedAt = &parsed
		}
	}

	processingTime := time.Since(startTime).Milliseconds()

	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "success",
		InvoiceNumber:  result.InvoiceNumber,
		ExternalID:     result.ExternalID,
		IssuedAt:       issuedAt,
		DocumentJSON:   fullDocument,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}

	if result.AuditData != nil {
		resp.AuditRequestURL = result.AuditData.RequestURL
		resp.AuditRequestPayload = toMapPayload(result.AuditData.RequestPayload)
		resp.AuditResponseStatus = result.AuditData.ResponseStatus
		resp.AuditResponseBody = result.AuditData.ResponseBody
	}

	if cashReceiptAudit != nil {
		resp.CashReceiptRequestURL = cashReceiptAudit.RequestURL
		resp.CashReceiptRequestPayload = cashReceiptAudit.RequestPayload
		resp.CashReceiptResponseStatus = cashReceiptAudit.ResponseStatus
		resp.CashReceiptResponseBody = cashReceiptAudit.ResponseBody
	}

	return resp
}
