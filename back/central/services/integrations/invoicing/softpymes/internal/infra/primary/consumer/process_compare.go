package consumer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
)

// processCompareRequest obtiene documentos del proveedor en el rango de fechas
// y publica un CompareResponseMessage con todos los documentos encontrados.
func (c *InvoiceRequestConsumer) processCompareRequest(
	ctx context.Context,
	request *InvoiceRequestMessage,
) error {
	// 1. Extraer parámetros del Config
	dateFrom, _ := request.InvoiceData.Config["date_from"].(string)
	dateTo, _ := request.InvoiceData.Config["date_to"].(string)
	businessID := uint(0)
	if bid, ok := request.InvoiceData.Config["business_id"].(float64); ok {
		businessID = uint(bid)
	}

	c.log.Info(ctx).
		Str("date_from", dateFrom).
		Str("date_to", dateTo).
		Uint("business_id", businessID).
		Str("correlation_id", request.CorrelationID).
		Msg("Starting compare request")

	// Helper para publicar error en el canal de comparación
	publishErr := func(errMsg string) error {
		return c.responsePublisher.PublishCompareResponse(ctx, &queue.CompareResponseMessage{
			Operation:     "compare",
			CorrelationID: request.CorrelationID,
			BusinessID:    businessID,
			DateFrom:      dateFrom,
			DateTo:        dateTo,
			Error:         errMsg,
			Timestamp:     time.Now(),
		})
	}

	if dateFrom == "" || dateTo == "" {
		c.log.Error(ctx).Msg("date_from or date_to missing in compare config")
		return publishErr("date_from and date_to are required in compare config")
	}

	// 2. Obtener integración y credenciales
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		c.log.Error(ctx).Msg("integration_id is 0 in compare request")
		return publishErr("integration_id is 0")
	}

	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get integration for compare")
		return publishErr("failed to get integration: " + err.Error())
	}

	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_key")
		return publishErr("failed to decrypt api_key")
	}

	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_secret")
		return publishErr("failed to decrypt api_secret")
	}

	// 3. Combinar config de integración con config del mensaje
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	referer, _ := combinedConfig["referer"].(string)

	// 4. Resolver URL efectiva desde integration_type (base_url / base_url_test)
	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		c.log.Error(ctx).
			Uint("integration_id", integrationID).
			Msg("base_url no configurada en el tipo de integración Softpymes")
		return publishErr("base_url no configurada en el tipo de integración Softpymes (integration_types.base_url)")
	}

	c.log.Info(ctx).
		Bool("is_testing", integration.IsTesting).
		Str("effective_url", effectiveURL).
		Msg("Resolved effective Softpymes URL for compare")

	// 5. Paginación: obtener todos los documentos del proveedor
	allDocs := make([]queue.CompareDocument, 0)
	pageSize := 20
	pageSizeStr := strconv.Itoa(pageSize)

	for page := 1; ; page++ {
		pageStr := strconv.Itoa(page)

		c.log.Info(ctx).
			Int("page", page).
			Str("date_from", dateFrom).
			Str("date_to", dateTo).
			Msg("Fetching documents page from Softpymes")

		docs, err := c.softpymesClient.ListDocuments(ctx, apiKey, apiSecret, referer, ports.ListDocumentsParams{
			DateFrom: dateFrom,
			DateTo:   dateTo,
			Page:     &pageStr,
			PageSize: &pageSizeStr,
		}, effectiveURL)
		if err != nil {
			c.log.Error(ctx).Err(err).Int("page", page).Msg("Failed to list documents")
			return publishErr(fmt.Sprintf("failed to list documents (page %d): %s", page, err.Error()))
		}

		for _, doc := range docs {
			details := make([]queue.CompareDocumentDetail, 0, len(doc.Details))
			for _, d := range doc.Details {
				details = append(details, queue.CompareDocumentDetail{
					ItemCode: d.ItemCode,
					ItemName: d.ItemName,
					Quantity: d.Quantity,
					Value:    d.Value,
					IVA:      d.IVA,
				})
			}
			allDocs = append(allDocs, queue.CompareDocument{
				DocumentNumber: doc.DocumentNumber,
				DocumentDate:   doc.DocumentDate,
				Total:          doc.Total,
				CustomerNit:    doc.CustomerNit,
				CustomerName:   doc.CustomerName,
				Comment:        doc.Comment,
				Prefix:         doc.Prefix,
				Details:        details,
			})
		}

		c.log.Info(ctx).
			Int("page", page).
			Int("page_count", len(docs)).
			Int("total_accumulated", len(allDocs)).
			Msg("Documents page fetched")

		// Última página cuando se devuelven menos registros que el tamaño de página
		if len(docs) < pageSize {
			break
		}
	}

	c.log.Info(ctx).
		Int("total_documents", len(allDocs)).
		Str("correlation_id", request.CorrelationID).
		Msg("All provider documents fetched, publishing compare response")

	// 6. Publicar resultado
	return c.responsePublisher.PublishCompareResponse(ctx, &queue.CompareResponseMessage{
		Operation:         "compare",
		CorrelationID:     request.CorrelationID,
		BusinessID:        businessID,
		DateFrom:          dateFrom,
		DateTo:            dateTo,
		ProviderDocuments: allDocs,
		Timestamp:         time.Now(),
	})
}
