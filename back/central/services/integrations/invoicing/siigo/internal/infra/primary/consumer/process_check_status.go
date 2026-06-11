package consumer

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
)

func (c *InvoiceRequestConsumer) processCheckStatus(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	externalID, _ := request.InvoiceData.Config["external_id"].(string)
	if externalID == "" {
		c.log.Warn(ctx).
			Uint("invoice_id", request.InvoiceID).
			Msg("external_id missing in check_status config - cannot query Siigo, keeping pending")
		return c.pendingValidationResponse(request, startTime, "external_id no disponible para consultar estado en Siigo")
	}

	ictx, errCode, err := c.resolveIntegration(ctx, request)
	if err != nil {
		c.log.Error(ctx).Err(err).Uint("invoice_id", request.InvoiceID).Msg("Failed to resolve integration for check_status")
		return c.createOperationErrorResponse(request, "check_status", errCode, err.Error(), startTime, nil)
	}

	detail, err := c.siigoClient.GetInvoiceByID(ctx, ictx.Credentials, externalID)
	if err != nil {
		c.log.Warn(ctx).Err(err).
			Uint("invoice_id", request.InvoiceID).
			Str("siigo_invoice_id", externalID).
			Msg("Failed to get invoice from Siigo - keeping as pending")
		return c.pendingValidationResponse(request, startTime, "Check status failed, will retry: "+err.Error())
	}

	stampStatus := strings.ToLower(detail.StampStatus)

	switch stampStatus {
	case "rejected", "error", "failed":
		errorMsg := "DIAN rechazo la factura"
		if stampErrors, stampErr := c.siigoClient.GetStampErrors(ctx, ictx.Credentials, externalID); stampErr == nil && len(stampErrors) > 0 {
			messages := make([]string, 0, len(stampErrors))
			for _, e := range stampErrors {
				messages = append(messages, e.Message)
			}
			errorMsg = "DIAN rechazo la factura: " + strings.Join(messages, "; ")
		}
		c.log.Error(ctx).
			Uint("invoice_id", request.InvoiceID).
			Str("siigo_invoice_id", externalID).
			Str("stamp_status", detail.StampStatus).
			Msg("Siigo invoice rejected by DIAN")
		return c.createOperationErrorResponse(request, "check_status", "dian_rejected", errorMsg, startTime, nil)

	case "pending", "processing", "sent", "waiting", "inprocess":
		c.log.Info(ctx).
			Uint("invoice_id", request.InvoiceID).
			Str("siigo_invoice_id", externalID).
			Str("stamp_status", detail.StampStatus).
			Msg("Siigo invoice still validating in DIAN - keeping pending")
		return c.pendingValidationResponse(request, startTime, "DIAN still validating (stamp status: "+detail.StampStatus+")")
	}

	var issuedAt *time.Time
	if detail.Date != "" {
		if parsed, parseErr := time.Parse("2006-01-02", detail.Date); parseErr == nil {
			issuedAt = &parsed
		}
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("siigo_invoice_id", externalID).
		Str("invoice_number", detail.Name).
		Str("stamp_status", detail.StampStatus).
		Msg("Siigo invoice confirmed")

	processingTime := time.Since(startTime).Milliseconds()
	return &queue.InvoiceResponseMessage{
		InvoiceID:     request.InvoiceID,
		Provider:      "siigo",
		Status:        "success",
		InvoiceNumber: detail.Name,
		ExternalID:    detail.ID,
		CUFE:          detail.CUFE,
		InvoiceURL:    detail.PublicURL,
		IssuedAt:      issuedAt,
		DocumentJSON: map[string]interface{}{
			"siigo_id":     detail.ID,
			"invoice_name": detail.Name,
			"total":        detail.Total,
			"balance":      detail.Balance,
			"status":       detail.Status,
			"stamp_status": detail.StampStatus,
			"public_url":   detail.PublicURL,
		},
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}
}

func (c *InvoiceRequestConsumer) pendingValidationResponse(
	request *InvoiceRequestMessage,
	startTime time.Time,
	message string,
) *queue.InvoiceResponseMessage {
	return &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "siigo",
		Status:         "pending_validation",
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: time.Since(startTime).Milliseconds(),
		Error:          message,
	}
}
