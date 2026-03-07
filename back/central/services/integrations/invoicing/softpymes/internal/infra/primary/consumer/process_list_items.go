package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
)

// processListItemsRequest obtiene ítems del catálogo del proveedor con paginación completa
// y publica un ListItemsResponseMessage con todos los ítems encontrados.
func (c *InvoiceRequestConsumer) processListItemsRequest(
	ctx context.Context,
	request *InvoiceRequestMessage,
) error {
	// 1. Extraer parámetros del Config
	businessID := uint(0)
	if bid, ok := request.InvoiceData.Config["business_id"].(float64); ok {
		businessID = uint(bid)
	}

	c.log.Info(ctx).
		Uint("business_id", businessID).
		Str("correlation_id", request.CorrelationID).
		Msg("Starting list_items request")

	// Helper para publicar error en el canal de list_items
	publishErr := func(errMsg string) error {
		return c.responsePublisher.PublishListItemsResponse(ctx, &queue.ListItemsResponseMessage{
			Operation:     "list_items",
			CorrelationID: request.CorrelationID,
			BusinessID:    businessID,
			Error:         errMsg,
			Timestamp:     time.Now(),
		})
	}

	// 2. Obtener integración y credenciales
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		c.log.Error(ctx).Msg("integration_id is 0 in list_items request")
		return publishErr("integration_id is 0")
	}

	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get integration for list_items")
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

	// 4. Resolver URL efectiva desde integration_type
	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		c.log.Error(ctx).
			Uint("integration_id", integrationID).
			Msg("base_url no configurada en el tipo de integración Softpymes")
		return publishErr("base_url no configurada en el tipo de integración Softpymes")
	}

	c.log.Info(ctx).
		Bool("is_testing", integration.IsTesting).
		Str("effective_url", effectiveURL).
		Msg("Resolved effective Softpymes URL for list_items")

	// 5. Paginación: obtener todos los ítems del proveedor
	allItems := make([]queue.ListItemsItem, 0)
	pageSize := 20

	for page := 1; ; page++ {
		c.log.Info(ctx).
			Int("page", page).
			Msg("Fetching items page from Softpymes")

		items, err := c.softpymesClient.ListItems(ctx, apiKey, apiSecret, referer, effectiveURL, page, pageSize)
		if err != nil {
			c.log.Error(ctx).Err(err).Int("page", page).Msg("Failed to list items")
			return publishErr(fmt.Sprintf("failed to list items (page %d): %s", page, err.Error()))
		}

		for _, item := range items {
			allItems = append(allItems, queue.ListItemsItem{
				ItemCode:      item.ItemCode,
				ItemName:      item.ItemName,
				ItemPrice:     item.ItemPrice,
				UnitCost:      item.UnitCost,
				Description:   item.Description,
				MinimumStock:  item.MinimumStock,
				OrderQuantity: item.OrderQuantity,
			})
		}

		c.log.Info(ctx).
			Int("page", page).
			Int("page_count", len(items)).
			Int("total_accumulated", len(allItems)).
			Msg("Items page fetched")

		// Última página cuando se devuelven menos registros que el tamaño de página
		if len(items) < pageSize {
			break
		}
	}

	c.log.Info(ctx).
		Int("total_items", len(allItems)).
		Str("correlation_id", request.CorrelationID).
		Msg("All provider items fetched, publishing list_items response")

	// 6. Publicar resultado
	return c.responsePublisher.PublishListItemsResponse(ctx, &queue.ListItemsResponseMessage{
		Operation:     "list_items",
		CorrelationID: request.CorrelationID,
		BusinessID:    businessID,
		Items:         allItems,
		Timestamp:     time.Now(),
	})
}
