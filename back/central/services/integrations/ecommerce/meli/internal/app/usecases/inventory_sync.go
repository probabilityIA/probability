package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func toUint(v interface{}) uint {
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
	cfg.SingleWarehouseID = toUint(config["inventory_single_warehouse_id"])
	if raw, ok := config["inventory_warehouse_ids"].([]interface{}); ok {
		for _, item := range raw {
			if id := toUint(item); id > 0 {
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

func (uc *meliUseCase) SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
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

	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return err
	}

	mapped, err := uc.inventoryRepo.ListMappedItems(ctx, uint(integIDUint))
	if err != nil {
		return fmt.Errorf("listing mapped items: %w", err)
	}

	productIDs := make([]string, 0, len(mapped))
	for _, m := range mapped {
		productIDs = append(productIDs, m.ProductID)
	}
	stock, err := uc.inventoryRepo.GetStockForProducts(ctx, productIDs, warehouseIDs)
	if err != nil {
		return fmt.Errorf("getting stock: %w", err)
	}

	total := len(mapped)
	uc.emitInventoryEvent(ctx, businessID, uint(integIDUint), "meli.inventory.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
	})

	updated, unchanged, skipped, failed := 0, 0, 0, 0
	for i, m := range mapped {
		qty := stock[m.ProductID]
		if uerr := uc.client.UpdateStock(ctx, accessToken, m.ExternalItemID, qty); uerr != nil {
			if uerr == domain.ErrTokenExpired {
				newToken, rerr := uc.EnsureValidToken(ctx, integrationID)
				if rerr == nil {
					accessToken = newToken
					if retry := uc.client.UpdateStock(ctx, accessToken, m.ExternalItemID, qty); retry == nil {
						updated++
						uc.maybeInventoryProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, updated, unchanged, skipped, failed)
						continue
					}
				}
			}
			uc.logger.Error(ctx).Err(uerr).Str("sku", m.SKU).Str("item_id", m.ExternalItemID).Msg("Error al actualizar stock en MercadoLibre")
			failed++
		} else {
			updated++
		}
		uc.maybeInventoryProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, updated, unchanged, skipped, failed)
	}

	uc.emitInventoryEvent(ctx, businessID, uint(integIDUint), "meli.inventory.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
		"updated":        updated,
		"unchanged":      unchanged,
		"skipped":        skipped,
		"failed":         failed,
	})
	return nil
}

func (uc *meliUseCase) maybeInventoryProgress(ctx context.Context, businessID, integrationID uint, correlationID string, processed, total, updated, unchanged, skipped, failed int) {
	if processed%productSyncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitInventoryEvent(ctx, businessID, integrationID, "meli.inventory.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"updated":        updated,
		"unchanged":      unchanged,
		"skipped":        skipped,
		"failed":         failed,
	})
}

func (uc *meliUseCase) emitInventoryEvent(ctx context.Context, businessID, integrationID uint, eventType string, data map[string]interface{}) {
	uc.emitSyncEvent(ctx, businessID, integrationID, eventType, data)
}
