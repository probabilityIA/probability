package consumer

import (
	"context"
	"fmt"
	"time"

	siigoDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
)

// processCreateJournal procesa la creación de un comprobante contable en Siigo
func (c *InvoiceRequestConsumer) processCreateJournal(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	// 1. Obtener integration_id del DTO
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		c.log.Error(ctx).Msg("integration_id is 0 in invoice_data")
		return c.createErrorResponse(request, "missing_integration_id", "integration_id is 0", startTime, nil)
	}

	// 2. Obtener integración desde IntegrationCore
	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get integration")
		return c.createErrorResponse(request, "integration_not_found", err.Error(), startTime, nil)
	}

	// 3. Desencriptar credenciales de Siigo
	username, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "username")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt username")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt username", startTime, nil)
	}

	accessKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "access_key")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt access_key")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt access_key", startTime, nil)
	}

	accountID, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "account_id")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt account_id")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt account_id", startTime, nil)
	}

	partnerID, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "partner_id")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt partner_id")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt partner_id", startTime, nil)
	}

	// api_url es opcional
	apiURL, _ := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_url")

	// Resolver URL efectiva: si is_testing, usar base_url_test del integration_type
	effectiveURL := apiURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}

	c.log.Info(ctx).
		Bool("is_testing", integration.IsTesting).
		Str("effective_url", effectiveURL).
		Msg("Resolved effective Siigo URL for journal")

	// 4. Combinar config de integración con config del request
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	// 5. Mapear items del mensaje a JournalItemData
	journalItems := make([]siigoDtos.JournalItemData, 0, len(request.InvoiceData.Items))
	for _, item := range request.InvoiceData.Items {
		journalItems = append(journalItems, siigoDtos.JournalItemData{
			SKU:        item.SKU,
			Name:       item.Name,
			Quantity:   item.Quantity,
			TotalPrice: item.TotalPrice,
			CustomerDNI: request.InvoiceData.Customer.DNI,
		})
	}

	// 6. Construir request tipado para el cliente Siigo
	journalReq := &siigoDtos.CreateJournalRequest{
		Items:       journalItems,
		Currency:    request.InvoiceData.Currency,
		OrderID:     request.InvoiceData.OrderID,
		Credentials: siigoDtos.Credentials{
			Username:  username,
			AccessKey: accessKey,
			AccountID: accountID,
			PartnerID: partnerID,
			BaseURL:   effectiveURL,
		},
		Config: combinedConfig,
	}

	// 7. Llamar al cliente HTTP de Siigo
	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("order_id", request.InvoiceData.OrderID).
		Int("items_count", len(journalItems)).
		Msg("Calling Siigo API for journal creation")

	result, err := c.siigoClient.CreateJournal(ctx, journalReq)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", request.InvoiceID).
			Msg("Siigo journal API call failed")

		var auditData *siigoDtos.AuditData
		if result != nil {
			auditData = result.AuditData
		}
		return c.createErrorResponse(request, "api_error", err.Error(), startTime, auditData)
	}

	// 8. Construir response exitosa
	// Reutiliza InvoiceResponseMessage: JournalName -> InvoiceNumber, JournalID -> ExternalID
	processingTime := time.Since(startTime).Milliseconds()

	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "siigo",
		Status:         "success",
		Operation:      "create_journal",
		InvoiceNumber:  result.JournalName,
		ExternalID:     result.JournalID,
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

	return resp
}
