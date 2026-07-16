package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/utils"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

const inventoryProgressBatch = 10

func invToUint(v interface{}) uint {
	switch val := v.(type) {
	case float64:
		return uint(val)
	case int:
		return uint(val)
	case int64:
		return uint(val)
	case string:
		if n, err := strconv.ParseUint(val, 10, 64); err == nil {
			return uint(n)
		}
	}
	return 0
}

func invToInt64(v interface{}) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int64:
		return val
	case int:
		return int64(val)
	case string:
		if n, err := strconv.ParseInt(val, 10, 64); err == nil {
			return n
		}
	}
	return 0
}

func parseLocationMappings(config map[string]interface{}) []domain.WarehouseMapping {
	raw, ok := config["shopify_location_mappings"].([]interface{})
	if !ok {
		return nil
	}
	var mappings []domain.WarehouseMapping
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		internalID := invToUint(m["internal_warehouse_id"])
		locationID := invToInt64(m["shopify_location_id"])
		if internalID == 0 || locationID == 0 {
			continue
		}
		mappings = append(mappings, domain.WarehouseMapping{
			InternalWarehouseID: internalID,
			ShopifyLocationID:   locationID,
		})
	}
	return mappings
}

func parseInventoryConfig(config map[string]interface{}) domain.InventoryConfig {
	cfg := domain.InventoryConfig{Mode: "sum"}
	if v, ok := config["inventory_warehouse_mode"].(string); ok && v != "" {
		cfg.Mode = v
	}
	cfg.SingleWarehouseID = invToUint(config["inventory_single_warehouse_id"])
	if raw, ok := config["inventory_warehouse_ids"].([]interface{}); ok {
		for _, item := range raw {
			if id := invToUint(item); id > 0 {
				cfg.WarehouseIDs = append(cfg.WarehouseIDs, id)
			}
		}
	}
	if v, ok := config["inventory_sync_enabled"].(bool); ok {
		cfg.Enabled = v
	}
	cfg.DefaultLocationID = invToInt64(config["shopify_default_location_id"])
	cfg.LocationMappings = parseLocationMappings(config)
	return cfg
}

func resolveWarehouseIDs(cfg domain.InventoryConfig) []uint {
	if cfg.Mode == "single" {
		if cfg.SingleWarehouseID > 0 {
			return []uint{cfg.SingleWarehouseID}
		}
		return nil
	}
	return cfg.WarehouseIDs
}

func resolveInventoryItemID(product *domain.ShopifyProduct, sku string) int64 {
	if product == nil {
		return 0
	}
	if len(product.Variants) > 1 && sku != "" {
		for _, v := range product.Variants {
			if v.SKU == sku && v.InventoryItemID != 0 {
				return v.InventoryItemID
			}
		}
	}
	for _, v := range product.Variants {
		if v.InventoryItemID != 0 {
			return v.InventoryItemID
		}
	}
	return 0
}

func (uc *SyncOrdersUseCase) resolveStoreAndToken(ctx context.Context, integration *domain.Integration, integrationID string) (string, string, error) {
	config, err := utils.NormalizeConfig(integration.Config, integration.Name)
	if err != nil {
		return "", "", err
	}
	storeDomain, err := utils.ExtractStoreName(config, integration.Name)
	if err != nil {
		return "", "", err
	}
	storeDomain = utils.ResolveEffectiveStoreDomain(integration, storeDomain)
	accessToken, err := utils.GetAccessToken(ctx, uc.integrationService, integrationID)
	if err != nil {
		return "", "", err
	}
	return storeDomain, accessToken, nil
}

func (uc *SyncOrdersUseCase) emitInventoryEvent(ctx context.Context, integrationID uint, businessID *uint, eventType string, data map[string]interface{}) {
	if uc.syncEventPublisher == nil {
		return
	}
	uc.syncEventPublisher.PublishSyncEvent(ctx, integrationID, businessID, eventType, data)
}

func (uc *SyncOrdersUseCase) resolveDefaultLocation(ctx context.Context, storeDomain, accessToken string, cfg domain.InventoryConfig) (int64, error) {
	if cfg.DefaultLocationID != 0 {
		return cfg.DefaultLocationID, nil
	}
	locations, err := uc.shopifyClient.GetLocations(ctx, storeDomain, accessToken)
	if err != nil {
		return 0, err
	}
	if len(locations) == 0 {
		return 0, fmt.Errorf("no locations found in Shopify")
	}
	return locations[0].ID, nil
}

func (uc *SyncOrdersUseCase) SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return fmt.Errorf("integration not found")
	}

	storeDomain, accessToken, err := uc.resolveStoreAndToken(ctx, integration, integrationID)
	if err != nil {
		return err
	}

	cfg := parseInventoryConfig(integration.Config)
	mapped, err := uc.inventoryRepo.ListMappedItems(ctx, uint(integIDUint))
	if err != nil {
		return fmt.Errorf("listing mapped items: %w", err)
	}

	productIDs := make([]string, 0, len(mapped))
	for _, m := range mapped {
		productIDs = append(productIDs, m.ProductID)
	}

	total := len(mapped)
	uc.emitInventoryEvent(ctx, uint(integIDUint), integration.BusinessID, "shopify.inventory.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
	})

	updated, unchanged, skipped, failed := 0, 0, 0, 0

	if len(cfg.LocationMappings) > 0 {
		warehouseIDs := make([]uint, 0, len(cfg.LocationMappings))
		for _, wm := range cfg.LocationMappings {
			warehouseIDs = append(warehouseIDs, wm.InternalWarehouseID)
		}
		byWarehouse, berr := uc.inventoryRepo.GetInventoryByWarehouses(ctx, productIDs, warehouseIDs)
		if berr != nil {
			return berr
		}
		for i, m := range mapped {
			invItemID, ok := uc.resolveItem(ctx, storeDomain, accessToken, m)
			if !ok {
				skipped++
				uc.progress(ctx, uint(integIDUint), integration.BusinessID, correlationID, i+1, total, updated, unchanged, skipped, failed)
				continue
			}
			stockByWh := byWarehouse[m.ProductID]
			ok2 := true
			for _, wm := range cfg.LocationMappings {
				qty := stockByWh[wm.InternalWarehouseID]
				if serr := uc.shopifyClient.SetInventoryLevel(ctx, storeDomain, accessToken, wm.ShopifyLocationID, invItemID, qty); serr != nil {
					uc.log.Error(ctx).Err(serr).Str("sku", m.SKU).Int64("location_id", wm.ShopifyLocationID).Msg("Error al fijar inventario en Shopify")
					ok2 = false
				}
			}
			if ok2 {
				updated++
			} else {
				failed++
			}
			uc.progress(ctx, uint(integIDUint), integration.BusinessID, correlationID, i+1, total, updated, unchanged, skipped, failed)
		}
	} else {
		warehouseIDs := resolveWarehouseIDs(cfg)
		stock, serr := uc.inventoryRepo.GetStockForProducts(ctx, productIDs, warehouseIDs)
		if serr != nil {
			return serr
		}
		locationID, lerr := uc.resolveDefaultLocation(ctx, storeDomain, accessToken, cfg)
		if lerr != nil {
			return lerr
		}
		for i, m := range mapped {
			invItemID, ok := uc.resolveItem(ctx, storeDomain, accessToken, m)
			if !ok {
				skipped++
				uc.progress(ctx, uint(integIDUint), integration.BusinessID, correlationID, i+1, total, updated, unchanged, skipped, failed)
				continue
			}
			qty := stock[m.ProductID]
			if serr := uc.shopifyClient.SetInventoryLevel(ctx, storeDomain, accessToken, locationID, invItemID, qty); serr != nil {
				uc.log.Error(ctx).Err(serr).Str("sku", m.SKU).Msg("Error al fijar inventario en Shopify")
				failed++
			} else {
				updated++
			}
			uc.progress(ctx, uint(integIDUint), integration.BusinessID, correlationID, i+1, total, updated, unchanged, skipped, failed)
		}
	}

	uc.emitInventoryEvent(ctx, uint(integIDUint), integration.BusinessID, "shopify.inventory.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
		"updated":        updated,
		"unchanged":      unchanged,
		"skipped":        skipped,
		"failed":         failed,
	})
	return nil
}

func (uc *SyncOrdersUseCase) resolveItem(ctx context.Context, storeDomain, accessToken string, m domain.MappedItem) (int64, bool) {
	product, err := uc.shopifyClient.GetProduct(ctx, storeDomain, accessToken, m.ExternalItemID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("external_product_id", m.ExternalItemID).Msg("Error al obtener producto de Shopify")
		return 0, false
	}
	invItemID := resolveInventoryItemID(product, m.SKU)
	if invItemID == 0 {
		uc.log.Info(ctx).Str("sku", m.SKU).Msg("No inventory_item_id para el producto de Shopify")
		return 0, false
	}
	return invItemID, true
}

func (uc *SyncOrdersUseCase) progress(ctx context.Context, integrationID uint, businessID *uint, correlationID string, processed, total, updated, unchanged, skipped, failed int) {
	if processed%inventoryProgressBatch != 0 && processed != total {
		return
	}
	uc.emitInventoryEvent(ctx, integrationID, businessID, "shopify.inventory.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"updated":        updated,
		"unchanged":      unchanged,
		"skipped":        skipped,
		"failed":         failed,
	})
}

func (uc *SyncOrdersUseCase) UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error {
	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return fmt.Errorf("integration not found")
	}
	if enabled, _ := integration.Config["inventory_sync_enabled"].(bool); !enabled {
		uc.log.Info(ctx).Str("integration_id", integrationID).Msg("Sync de inventario desactivado para Shopify, push omitido")
		return nil
	}
	storeDomain, accessToken, err := uc.resolveStoreAndToken(ctx, integration, integrationID)
	if err != nil {
		return err
	}
	cfg := parseInventoryConfig(integration.Config)
	locationID, err := uc.resolveDefaultLocation(ctx, storeDomain, accessToken, cfg)
	if err != nil {
		return err
	}
	product, err := uc.shopifyClient.GetProduct(ctx, storeDomain, accessToken, productExternalID)
	if err != nil {
		return err
	}
	invItemID := resolveInventoryItemID(product, "")
	if invItemID == 0 {
		return fmt.Errorf("no inventory_item_id for product %s", productExternalID)
	}
	return uc.shopifyClient.SetInventoryLevel(ctx, storeDomain, accessToken, locationID, invItemID, quantity)
}
