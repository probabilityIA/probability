package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const syncProgressBatch = 25

func (uc *jumpsellerUseCase) emitSyncEvent(ctx context.Context, businessID, integrationID uint, eventType string, data map[string]interface{}) {
	if uc.rabbit == nil {
		return
	}
	_ = rabbitmq.PublishEvent(ctx, uc.rabbit, rabbitmq.EventEnvelope{
		Type:          eventType,
		Category:      "integration",
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Data:          data,
	})
}

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
	cfg := domain.InventoryConfig{Mode: domain.InventoryModeSingle}
	if v, ok := config["inventory_warehouse_mode"].(string); ok && v == domain.InventoryModeMapped {
		cfg.Mode = domain.InventoryModeMapped
	}
	cfg.SingleWarehouseID = invToUint(config["inventory_single_warehouse_id"])
	if v, ok := config["inventory_sync_enabled"].(bool); ok {
		cfg.Enabled = v
	}
	cfg.DefaultLocationID = int64(invToUint(config["jumpseller_default_location_id"]))
	if raw, ok := config["jumpseller_location_mappings"].([]interface{}); ok {
		for _, item := range raw {
			entry, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			warehouseID := invToUint(entry["internal_warehouse_id"])
			locationID := int64(invToUint(entry["jumpseller_location_id"]))
			if warehouseID > 0 && locationID > 0 {
				cfg.LocationMappings = append(cfg.LocationMappings, domain.LocationMapping{
					InternalWarehouseID:  warehouseID,
					JumpsellerLocationID: locationID,
				})
			}
		}
	}
	return cfg
}

func resolveWarehouseIDs(cfg domain.InventoryConfig) []uint {
	if cfg.Mode == domain.InventoryModeMapped {
		ids := make([]uint, 0, len(cfg.LocationMappings))
		for _, m := range cfg.LocationMappings {
			ids = append(ids, m.InternalWarehouseID)
		}
		return ids
	}
	if cfg.SingleWarehouseID > 0 {
		return []uint{cfg.SingleWarehouseID}
	}
	return nil
}

func (uc *jumpsellerUseCase) SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	integration, cred, err := uc.resolveIntegrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	cfg := parseInventoryConfig(integration.Config)

	if groups := cfg.LocationGroups(); cfg.Mode == domain.InventoryModeMapped && len(groups) > 1 {
		uc.logger.Error(ctx).
			Str("integration_id", integrationID).
			Int("locations", len(groups)).
			Msg("El emparejamiento apunta a varias bodegas de Jumpseller y todavia no podemos escribir stock por bodega")
		return domain.ErrPerLocationStockNotSupported
	}

	warehouseIDs := resolveWarehouseIDs(cfg)

	mapped, err := uc.productRepo.ListMappedItems(ctx, uint(integIDUint))
	if err != nil {
		return fmt.Errorf("listing mapped items: %w", err)
	}

	productIDs := make([]string, 0, len(mapped))
	for _, item := range mapped {
		productIDs = append(productIDs, item.ProductID)
	}

	stock, err := uc.productRepo.GetStockForProducts(ctx, productIDs, warehouseIDs)
	if err != nil {
		return fmt.Errorf("getting stock: %w", err)
	}

	total := len(mapped)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.inventory.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
	})

	fails := &failedSKUs{}
	updated := 0
	for i, item := range mapped {
		qty := stock[item.ProductID]
		action := "updated"

		productID, variantID, perr := parseExternalProductID(item.ExternalItemID)
		if perr != nil {
			fails.add(item.SKU)
			action = "failed"
		} else {
			var uerr error
			if variantID > 0 {
				uerr = uc.client.SetVariantStock(ctx, cred, productID, variantID, qty)
			} else {
				uerr = uc.client.SetProductStock(ctx, cred, productID, qty)
			}
			if uerr != nil {
				uc.logger.Error(ctx).Err(uerr).
					Str("sku", item.SKU).
					Str("external_product_id", item.ExternalItemID).
					Msg("Error al actualizar stock en Jumpseller")
				fails.add(item.SKU)
				action = "failed"
			} else {
				updated++
			}
		}

		uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.inventory.sync.item", map[string]interface{}{
			"correlation_id": correlationID,
			"sku":            item.SKU,
			"quantity":       qty,
			"action":         action,
		})
		uc.maybeInventoryProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.inventory.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
		"updated":        updated,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *jumpsellerUseCase) maybeInventoryProgress(ctx context.Context, businessID, integrationID uint, correlationID string, processed, total, updated, failed int) {
	if processed%syncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "jumpseller.inventory.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"updated":        updated,
		"failed":         failed,
	})
}
