package consumer

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
)

func (c *InvoiceRequestConsumer) processCancelInvoice(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	externalID, _ := request.InvoiceData.Config["external_id"].(string)
	if externalID == "" {
		c.log.Error(ctx).Uint("invoice_id", request.InvoiceID).Msg("external_id missing in cancel config")
		return c.createOperationErrorResponse(request, "cancel", "missing_external_id", "external_id is required for cancellation", startTime, nil)
	}

	ictx, errCode, err := c.resolveIntegration(ctx, request)
	if err != nil {
		c.log.Error(ctx).Err(err).Uint("invoice_id", request.InvoiceID).Msg("Failed to resolve integration for cancel")
		return c.createOperationErrorResponse(request, "cancel", errCode, err.Error(), startTime, nil)
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("siigo_invoice_id", externalID).
		Bool("is_testing", ictx.IsTesting).
		Msg("Annulling invoice in Siigo")

	result, err := c.siigoClient.AnnulInvoice(ctx, ictx.Credentials, externalID)
	if err != nil {
		c.log.Error(ctx).Err(err).Uint("invoice_id", request.InvoiceID).Msg("Siigo annul failed")
		var auditData = resultAudit(result)
		return c.createOperationErrorResponse(request, "cancel", "api_error", err.Error(), startTime, auditData)
	}

	processingTime := time.Since(startTime).Milliseconds()
	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("siigo_invoice_id", externalID).
		Msg("Invoice annulled successfully in Siigo")

	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "siigo",
		Status:         "success",
		Operation:      "cancel",
		ExternalID:     externalID,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}

	if result != nil && result.AuditData != nil {
		resp.AuditRequestURL = result.AuditData.RequestURL
		resp.AuditRequestPayload = toMapPayload(result.AuditData.RequestPayload)
		resp.AuditResponseStatus = result.AuditData.ResponseStatus
		resp.AuditResponseBody = result.AuditData.ResponseBody
	}

	return resp
}
