package consumer

import (
	"context"
	"encoding/json"

	siigoDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type providerStockSyncMessage struct {
	CorrelationID string                  `json:"correlation_id"`
	BusinessID    uint                    `json:"business_id"`
	IntegrationID uint                    `json:"integration_id"`
	Provider      string                  `json:"provider"`
	Items         []providerStockSyncItem `json:"items"`
}

type providerStockSyncItem struct {
	SKU         string `json:"sku"`
	WarehouseID uint   `json:"warehouse_id"`
	Quantity    int    `json:"quantity"`
}

type warehouseMapping struct {
	warehouseID uint
	siigoWarehouseID    int
}

func (c *InvoiceRequestConsumer) processInventorySyncRequest(
	ctx context.Context,
	request *InvoiceRequestMessage,
) error {
	businessID := businessIDFromConfig(request.InvoiceData.Config)

	c.log.Info(ctx).
		Uint("business_id", businessID).
		Uint("integration_id", request.InvoiceData.IntegrationID).
		Str("correlation_id", request.CorrelationID).
		Msg("Starting Siigo inventory sync request")

	ictx, _, err := c.resolveIntegration(ctx, request)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to resolve integration for inventory sync")
		return c.publishInventorySync(ctx, &providerStockSyncMessage{
			CorrelationID: request.CorrelationID,
			BusinessID:    businessID,
			IntegrationID: request.InvoiceData.IntegrationID,
			Provider:      "siigo",
			Items:         []providerStockSyncItem{},
		})
	}

	if enabled, ok := ictx.Config["inventory_sync_enabled"].(bool); ok && !enabled {
		c.log.Warn(ctx).Msg("inventory_sync_enabled is false; proceeding with manual trigger anyway")
	}

	mode, _ := ictx.Config["inventory_warehouse_mode"].(string)
	if mode == "" {
		mode = "single"
	}
	singleWarehouseID := uintFromConfig(ictx.Config, "inventory_single_warehouse_id")
	mappings := parseWarehouseMappings(ictx.Config["inventory_warehouse_mappings"])

	products := make([]siigoDtos.ProductItem, 0)
	pageSize := 100
	for page := 1; ; page++ {
		batch, err := c.siigoClient.ListProducts(ctx, ictx.Credentials, page, pageSize)
		if err != nil {
			c.log.Error(ctx).Err(err).Int("page", page).Msg("Failed to list Siigo products for inventory sync")
			break
		}
		products = append(products, batch...)
		if len(batch) < pageSize {
			break
		}
	}

	if enabled, _ := ictx.Config["product_sync_enabled"].(bool); enabled {
		c.publishProductUpserts(ctx, businessID, request.InvoiceData.IntegrationID, products)
	}

	items := buildInventorySyncItems(products, mode, singleWarehouseID, mappings)

	c.log.Info(ctx).
		Int("products", len(products)).
		Int("items", len(items)).
		Str("mode", mode).
		Msg("Siigo inventory fetched, publishing to inventory sync queue")

	return c.publishInventorySync(ctx, &providerStockSyncMessage{
		CorrelationID: request.CorrelationID,
		BusinessID:    businessID,
		IntegrationID: request.InvoiceData.IntegrationID,
		Provider:      "siigo",
		Items:         items,
	})
}

func buildInventorySyncItems(
	products []siigoDtos.ProductItem,
	mode string,
	singleWarehouseID uint,
	mappings []warehouseMapping,
) []providerStockSyncItem {
	items := make([]providerStockSyncItem, 0, len(products))
	for _, p := range products {
		if !p.StockControl || p.Code == "" {
			continue
		}
		if mode == "mapped" {
			for _, m := range mappings {
				qty := 0.0
				for _, w := range p.Warehouses {
					if w.ID == m.siigoWarehouseID {
						qty = w.Quantity
						break
					}
				}
				items = append(items, providerStockSyncItem{
					SKU:         p.Code,
					WarehouseID: m.warehouseID,
					Quantity:    int(qty),
				})
			}
			continue
		}
		items = append(items, providerStockSyncItem{
			SKU:         p.Code,
			WarehouseID: singleWarehouseID,
			Quantity:    int(p.AvailableQuantity),
		})
	}
	return items
}

type productUpsertMessage struct {
	BusinessID     uint    `json:"business_id"`
	IntegrationID  uint    `json:"integration_id"`
	SKU            string  `json:"sku"`
	Name           string  `json:"name"`
	TrackInventory bool    `json:"track_inventory"`
	Price          float64 `json:"price"`
	ExternalID     string  `json:"external_id"`
}

func (c *InvoiceRequestConsumer) publishProductUpserts(ctx context.Context, businessID, integrationID uint, products []siigoDtos.ProductItem) {
	if err := c.rabbit.DeclareQueue(rabbitmq.QueueProductsProviderUpsert, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare products upsert queue")
		return
	}

	published := 0
	for _, p := range products {
		if p.Code == "" {
			continue
		}
		data, err := json.Marshal(productUpsertMessage{
			BusinessID:     businessID,
			IntegrationID:  integrationID,
			SKU:            p.Code,
			Name:           p.Name,
			TrackInventory: p.StockControl,
			Price:          p.Price,
			ExternalID:     p.ID,
		})
		if err != nil {
			continue
		}
		if err := c.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, data); err != nil {
			c.log.Error(ctx).Err(err).Str("sku", p.Code).Msg("Failed to publish product upsert")
			continue
		}
		published++
	}

	c.log.Info(ctx).Int("published", published).Uint("business_id", businessID).Msg("Siigo product upserts published")
}

func (c *InvoiceRequestConsumer) publishInventorySync(ctx context.Context, msg *providerStockSyncMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to marshal inventory sync message")
		return err
	}
	if err := c.rabbit.DeclareQueue(rabbitmq.QueueInventoryProviderSync, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare inventory sync queue")
		return err
	}
	if err := c.rabbit.Publish(ctx, rabbitmq.QueueInventoryProviderSync, data); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to publish inventory sync message")
		return err
	}
	return nil
}

func uintFromConfig(config map[string]interface{}, key string) uint {
	switch v := config[key].(type) {
	case float64:
		if v > 0 {
			return uint(v)
		}
	case int:
		if v > 0 {
			return uint(v)
		}
	}
	return 0
}

func parseWarehouseMappings(raw interface{}) []warehouseMapping {
	list, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	mappings := make([]warehouseMapping, 0, len(list))
	for _, entry := range list {
		m, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		warehouseIDValue := uintFromConfig(m, "warehouse_id")
		if warehouseIDValue == 0 {
			warehouseIDValue = uintFromConfig(m, "velocity_warehouse_id")
		}
		siigoID := 0
		switch v := m["siigo_warehouse_id"].(type) {
		case float64:
			siigoID = int(v)
		case int:
			siigoID = v
		}
		if warehouseIDValue == 0 {
			continue
		}
		mappings = append(mappings, warehouseMapping{warehouseID: warehouseIDValue, siigoWarehouseID: siigoID})
	}
	return mappings
}
