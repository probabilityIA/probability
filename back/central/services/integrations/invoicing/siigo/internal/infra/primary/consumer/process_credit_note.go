package consumer

import (
	"context"
	"time"

	siigoDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
)

func (c *InvoiceRequestConsumer) processCreditNote(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	externalID, _ := request.InvoiceData.Config["external_id"].(string)
	if externalID == "" {
		return c.createOperationErrorResponse(request, "credit_note", "missing_external_id", "external_id is required for credit note", startTime, nil)
	}

	invoiceNumber, _ := request.InvoiceData.Config["invoice_number"].(string)
	reason, _ := request.InvoiceData.Config["reason"].(string)
	noteType, _ := request.InvoiceData.Config["note_type"].(string)
	customerDNI, _ := request.InvoiceData.Config["customer_dni"].(string)
	amount := 0.0
	switch v := request.InvoiceData.Config["amount"].(type) {
	case float64:
		amount = v
	}

	ictx, errCode, err := c.resolveIntegration(ctx, request)
	if err != nil {
		c.log.Error(ctx).Err(err).Uint("invoice_id", request.InvoiceID).Msg("Failed to resolve integration for credit_note")
		return c.createOperationErrorResponse(request, "credit_note", errCode, err.Error(), startTime, nil)
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("invoice_external_id", externalID).
		Float64("amount", amount).
		Msg("Creating Siigo credit note")

	result, err := c.siigoClient.CreateCreditNote(ctx, &siigoDtos.CreateCreditNoteRequest{
		InvoiceExternalID: externalID,
		InvoiceNumber:     invoiceNumber,
		Amount:            amount,
		Reason:            reason,
		NoteType:          noteType,
		CustomerDNI:       customerDNI,
		Credentials:       ictx.Credentials,
		Config:            ictx.Config,
	})

	if err != nil {
		c.log.Error(ctx).Err(err).Uint("invoice_id", request.InvoiceID).Msg("Siigo credit note failed")
		var auditData *siigoDtos.AuditData
		if result != nil {
			auditData = result.AuditData
		}
		return c.createOperationErrorResponse(request, "credit_note", "credit_note_failed", err.Error(), startTime, auditData)
	}

	processingTime := time.Since(startTime).Milliseconds()
	resp := &queue.InvoiceResponseMessage{
		InvoiceID:     request.InvoiceID,
		Provider:      "siigo",
		Status:        "success",
		Operation:     "credit_note",
		InvoiceNumber: result.CreditNoteNumber,
		ExternalID:    result.CreditNoteID,
		CUFE:          result.CUFE,
		DocumentJSON: map[string]interface{}{
			"credit_note": result.ProviderInfo,
		},
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

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("credit_note_number", result.CreditNoteNumber).
		Msg("Siigo credit note generated successfully")

	return resp
}
