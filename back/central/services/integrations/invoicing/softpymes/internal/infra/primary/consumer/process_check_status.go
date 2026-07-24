package consumer

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
)

const (
	checkStatusMaxDaysBack    = 7
	checkStatusPageSize       = 50
	checkStatusMaxPagesPerDay = 60
)

func (c *InvoiceRequestConsumer) processCheckStatus(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	orderID := request.InvoiceData.OrderID
	if orderID == "" {
		return c.createErrorResponse(request, "missing_order_id", "order_id is required for check_status", startTime, nil)
	}

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

	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	referer, _ := combinedConfig["referer"].(string)

	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		return c.createErrorResponse(request, "missing_base_url", "base_url no configurada", startTime, nil)
	}

	loc, _ := time.LoadLocation("America/Bogota")
	now := time.Now().In(loc)

	searchComment := "order:" + orderID
	pageSize := strconv.Itoa(checkStatusPageSize)

	var foundDoc *ports.ListedDocument
	docsSearched := 0

	for daysBack := 0; daysBack <= checkStatusMaxDaysBack && foundDoc == nil; daysBack++ {
		day := now.AddDate(0, 0, -daysBack).Format("2006-01-02")

		for page := 1; page <= checkStatusMaxPagesPerDay; page++ {
			pageStr := strconv.Itoa(page)
			docs, listErr := c.softpymesClient.ListDocuments(ctx, apiKey, apiSecret, referer, ports.ListDocumentsParams{
				DateFrom: day,
				DateTo:   day,
				PageSize: &pageSize,
				Page:     &pageStr,
			}, effectiveURL)

			if listErr != nil {
				c.log.Warn(ctx).Err(listErr).
					Uint("invoice_id", request.InvoiceID).
					Str("order_id", orderID).
					Str("day", day).
					Int("page", page).
					Msg("Failed to search documents in Softpymes — keeping as pending")

				processingTime := time.Since(startTime).Milliseconds()
				return &queue.InvoiceResponseMessage{
					InvoiceID:      request.InvoiceID,
					Provider:       "softpymes",
					Status:         "pending_validation",
					CorrelationID:  request.CorrelationID,
					Timestamp:      time.Now(),
					ProcessingTime: processingTime,
					Error:          "Check status failed, will retry: " + listErr.Error(),
				}
			}

			docsSearched += len(docs)

			for i := range docs {
				if strings.Contains(docs[i].Comment, searchComment) {
					foundDoc = &docs[i]
					break
				}
			}

			if foundDoc != nil || len(docs) < checkStatusPageSize {
				break
			}
		}
	}

	if foundDoc != nil {
		doc := *foundDoc
		c.log.Info(ctx).
			Uint("invoice_id", request.InvoiceID).
			Str("order_id", orderID).
			Str("document_number", doc.DocumentNumber).
			Msg("Found existing document in Softpymes for pending invoice")

		var fullDocument map[string]interface{}
		if doc.DocumentNumber != "" {
			fullDoc, docErr := c.softpymesClient.GetDocumentByNumber(ctx, apiKey, apiSecret, referer, doc.DocumentNumber, effectiveURL)
			if docErr != nil {
				c.log.Warn(ctx).Err(docErr).Str("document_number", doc.DocumentNumber).Msg("Failed to get full document")
			} else {
				fullDocument = fullDoc
				if docNum, ok := fullDocument["documentNumber"].(string); ok && docNum != "" {
					doc.DocumentNumber = docNum
				}
			}
		}

		c.sendCashReceiptIfConfigured(ctx, fullDocument, combinedConfig, apiKey, apiSecret, referer, effectiveURL, request.InvoiceID)

		var issuedAt *time.Time
		if doc.DocumentDate != "" {
			if parsed, parseErr := time.Parse("2006-01-02", doc.DocumentDate); parseErr == nil {
				issuedAt = &parsed
			}
		}

		processingTime := time.Since(startTime).Milliseconds()
		return &queue.InvoiceResponseMessage{
			InvoiceID:      request.InvoiceID,
			Provider:       "softpymes",
			Status:         "success",
			InvoiceNumber:  doc.DocumentNumber,
			ExternalID:     doc.DocumentNumber,
			IssuedAt:       issuedAt,
			DocumentJSON:   fullDocument,
			CorrelationID:  request.CorrelationID,
			Timestamp:      time.Now(),
			ProcessingTime: processingTime,
		}
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("order_id", orderID).
		Int("docs_searched", docsSearched).
		Int("days_searched", checkStatusMaxDaysBack+1).
		Msg("Document not found in Softpymes — DIAN still validating, keeping pending")

	processingTime := time.Since(startTime).Milliseconds()
	return &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "pending_validation",
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
		Error:          "Document not found yet — DIAN still validating",
	}
}
