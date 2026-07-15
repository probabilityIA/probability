package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

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

func (uc *wooCommerceUseCase) SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return domain.ErrIntegrationNotFound
	}

	cfg := parseInventoryConfig(integration.Config)
	warehouseIDs := resolveWarehouseIDs(cfg)

	storeURL, err := extractString(integration.Config, "store_url")
	if err != nil {
		return domain.ErrMissingStoreURL
	}
	storeURL = resolveEffectiveStoreURL(integration, storeURL)

	consumerKey, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_key")
	if err != nil {
		return fmt.Errorf("decrypting consumer_key: %w", err)
	}
	consumerSecret, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_secret")
	if err != nil {
		return fmt.Errorf("decrypting consumer_secret: %w", err)
	}

	mapped, err := uc.productRepo.ListMappedItems(ctx, uint(integIDUint))
	if err != nil {
		return fmt.Errorf("listing mapped items: %w", err)
	}

	productIDs := make([]string, 0, len(mapped))
	for _, m := range mapped {
		productIDs = append(productIDs, m.ProductID)
	}
	stock, err := uc.productRepo.GetStockForProducts(ctx, productIDs, warehouseIDs)
	if err != nil {
		return fmt.Errorf("getting stock: %w", err)
	}

	total := len(mapped)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woo.inventory.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
	})

	updated, unchanged, skipped, failed := 0, 0, 0, 0
	for i, m := range mapped {
		qty := stock[m.ProductID]
		action := "updated"
		if uerr := uc.client.UpdateProductStock(ctx, storeURL, consumerKey, consumerSecret, m.ExternalItemID, qty); uerr != nil {
			uc.logger.Error(ctx).Err(uerr).Str("sku", m.SKU).Str("external_product_id", m.ExternalItemID).Msg("Error al actualizar stock en WooCommerce")
			failed++
			action = "failed"
		} else {
			updated++
		}
		uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woo.inventory.sync.item", map[string]interface{}{
			"correlation_id": correlationID,
			"sku":            m.SKU,
			"quantity":       qty,
			"action":         action,
		})
		uc.maybeInventoryProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, updated, unchanged, skipped, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woo.inventory.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
		"updated":        updated,
		"unchanged":      unchanged,
		"skipped":        skipped,
		"failed":         failed,
	})
	return nil
}

func (uc *wooCommerceUseCase) maybeInventoryProgress(ctx context.Context, businessID, integrationID uint, correlationID string, processed, total, updated, unchanged, skipped, failed int) {
	if processed%productSyncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "woo.inventory.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"updated":        updated,
		"unchanged":      unchanged,
		"skipped":        skipped,
		"failed":         failed,
	})
}
