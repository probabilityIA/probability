package consumer

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
)

func (c *InvoiceRequestConsumer) processListSiigoWarehousesRequest(
	ctx context.Context,
	request *InvoiceRequestMessage,
) error {
	businessID := businessIDFromConfig(request.InvoiceData.Config)

	c.log.Info(ctx).
		Uint("business_id", businessID).
		Str("correlation_id", request.CorrelationID).
		Msg("Starting Siigo list_siigo_warehouses request")

	publishErr := func(errMsg string) error {
		return c.responsePublisher.PublishListSiigoWarehousesResponse(ctx, &queue.ListSiigoWarehousesResponseMessage{
			Operation:     "list_siigo_warehouses",
			CorrelationID: request.CorrelationID,
			BusinessID:    businessID,
			Error:         errMsg,
			Timestamp:     time.Now(),
		})
	}

	ictx, _, err := c.resolveIntegration(ctx, request)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to resolve integration for list_siigo_warehouses")
		return publishErr("failed to resolve integration: " + err.Error())
	}

	warehouses, err := c.siigoClient.ListWarehouses(ctx, ictx.Credentials)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to list Siigo warehouses")
		return publishErr("failed to list warehouses: " + err.Error())
	}

	items := make([]queue.SiigoWarehouseItem, 0, len(warehouses))
	for _, w := range warehouses {
		items = append(items, queue.SiigoWarehouseItem{ID: w.ID, Name: w.Name})
	}

	c.log.Info(ctx).
		Int("total_warehouses", len(items)).
		Str("correlation_id", request.CorrelationID).
		Msg("Siigo warehouses fetched, publishing list_siigo_warehouses response")

	return c.responsePublisher.PublishListSiigoWarehousesResponse(ctx, &queue.ListSiigoWarehousesResponseMessage{
		Operation:     "list_siigo_warehouses",
		CorrelationID: request.CorrelationID,
		BusinessID:    businessID,
		Items:         items,
		Timestamp:     time.Now(),
	})
}
