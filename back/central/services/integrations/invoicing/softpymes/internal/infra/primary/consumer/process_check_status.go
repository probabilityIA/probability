package consumer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
)

// processCheckStatus busca un documento existente en Softpymes para una factura pendiente de DIAN.
// NO crea documentos nuevos — solo consulta ListDocuments y busca por comment "order:<UUID>".
// Si encuentra el documento, retorna success con el número de factura.
// Si no lo encuentra, retorna pending_validation para programar otro check más tarde.
func (c *InvoiceRequestConsumer) processCheckStatus(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	orderID := request.InvoiceData.OrderID
	if orderID == "" {
		return c.createErrorResponse(request, "missing_order_id", "order_id is required for check_status", startTime, nil)
	}

	// 1. Obtener integración y credenciales
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

	// 2. Resolver config y URL
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

	// 3. Buscar documento por order_id en comment field
	branchCode := "001"
	if bc, ok := combinedConfig["branch_code"].(string); ok && bc != "" {
		branchCode = bc
	}

	// Rango: desde hace 30 días hasta hoy (máximo permitido por Softpymes)
	loc, _ := time.LoadLocation("America/Bogota")
	now := time.Now().In(loc)
	dateTo := now.Format("2006-01-02")
	dateFrom := now.AddDate(0, 0, -30).Format("2006-01-02")

	pageSize := "50"
	docs, err := c.softpymesClient.ListDocuments(ctx, apiKey, apiSecret, referer, ports.ListDocumentsParams{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		PageSize: &pageSize,
	}, effectiveURL)

	if err != nil {
		c.log.Warn(ctx).Err(err).
			Uint("invoice_id", request.InvoiceID).
			Str("order_id", orderID).
			Msg("Failed to search documents in Softpymes — keeping as pending")

		// No marcamos como error — es una falla temporal de conectividad
		processingTime := time.Since(startTime).Milliseconds()
		return &queue.InvoiceResponseMessage{
			InvoiceID:      request.InvoiceID,
			Provider:       "softpymes",
			Status:         "pending_validation",
			CorrelationID:  request.CorrelationID,
			Timestamp:      time.Now(),
			ProcessingTime: processingTime,
			Error:          "Check status failed, will retry: " + err.Error(),
		}
	}

	// 4. Buscar documento con comment "order:<UUID>"
	searchComment := "order:" + orderID
	_ = branchCode // branchCode ya está en el filtro implícito de la búsqueda
	for _, doc := range docs {
		if strings.Contains(doc.Comment, searchComment) {
			c.log.Info(ctx).
				Uint("invoice_id", request.InvoiceID).
				Str("order_id", orderID).
				Str("document_number", doc.DocumentNumber).
				Msg("Found existing document in Softpymes for pending invoice")

			// 5. Obtener documento completo
			var fullDocument map[string]interface{}
			if doc.DocumentNumber != "" {
				fullDoc, err := c.softpymesClient.GetDocumentByNumber(ctx, apiKey, apiSecret, referer, doc.DocumentNumber, effectiveURL)
				if err != nil {
					c.log.Warn(ctx).Err(err).Str("document_number", doc.DocumentNumber).Msg("Failed to get full document")
				} else {
					fullDocument = fullDoc
					// Usar documentNumber canónico del full document
					if docNum, ok := fullDocument["documentNumber"].(string); ok && docNum != "" {
						doc.DocumentNumber = docNum
					}
				}
			}

			// 6. Send cash receipt if configured (non-fatal)
			c.sendCashReceiptIfConfigured(ctx, fullDocument, combinedConfig, apiKey, apiSecret, referer, effectiveURL, request.InvoiceID)

			// 7. Parsear fecha
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
	}

	// 7. No encontrado — DIAN sigue validando, mantener pending
	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("order_id", orderID).
		Int("docs_searched", len(docs)).
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
