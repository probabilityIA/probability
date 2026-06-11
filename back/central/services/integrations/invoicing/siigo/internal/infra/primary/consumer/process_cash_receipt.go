package consumer

import (
	"context"
	"time"

	siigoDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
)

func (c *InvoiceRequestConsumer) processCashReceipt(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	invoiceNumber, _ := request.InvoiceData.Config["invoice_number"].(string)
	if invoiceNumber == "" {
		return c.createOperationErrorResponse(request, "cash_receipt", "missing_invoice_number", "invoice_number is required for cash_receipt operation", startTime, nil)
	}

	ictx, errCode, err := c.resolveIntegration(ctx, request)
	if err != nil {
		c.log.Error(ctx).Err(err).Uint("invoice_id", request.InvoiceID).Msg("Failed to resolve integration for cash_receipt")
		return c.createOperationErrorResponse(request, "cash_receipt", errCode, err.Error(), startTime, nil)
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("invoice_number", invoiceNumber).
		Msg("Creating Siigo cash receipt")

	result, err := c.siigoClient.CreateCashReceipt(ctx, &siigoDtos.CreateCashReceiptRequest{
		InvoiceNumber: invoiceNumber,
		Credentials:   ictx.Credentials,
		Config:        ictx.Config,
	})

	if err != nil {
		c.log.Error(ctx).Err(err).
			Uint("invoice_id", request.InvoiceID).
			Str("invoice_number", invoiceNumber).
			Msg("Siigo cash receipt failed")

		resp := c.createOperationErrorResponse(request, "cash_receipt", "cash_receipt_failed", err.Error(), startTime, nil)
		if result != nil && result.AuditData != nil {
			resp.CashReceiptRequestURL = result.AuditData.RequestURL
			resp.CashReceiptRequestPayload = toMapPayload(result.AuditData.RequestPayload)
			resp.CashReceiptResponseStatus = result.AuditData.ResponseStatus
			resp.CashReceiptResponseBody = result.AuditData.ResponseBody
		}
		return resp
	}

	processingTime := time.Since(startTime).Milliseconds()
	resp := &queue.InvoiceResponseMessage{
		InvoiceID:     request.InvoiceID,
		Provider:      "siigo",
		Status:        "success",
		Operation:     "cash_receipt",
		InvoiceNumber: invoiceNumber,
		ExternalID:    result.ReceiptID,
		DocumentJSON: map[string]interface{}{
			"cash_receipt": result.ProviderInfo,
		},
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}

	if result.AuditData != nil {
		resp.CashReceiptRequestURL = result.AuditData.RequestURL
		resp.CashReceiptRequestPayload = toMapPayload(result.AuditData.RequestPayload)
		resp.CashReceiptResponseStatus = result.AuditData.ResponseStatus
		resp.CashReceiptResponseBody = result.AuditData.ResponseBody
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("invoice_number", invoiceNumber).
		Str("receipt_name", result.ReceiptName).
		Msg("Siigo cash receipt generated successfully")

	return resp
}
