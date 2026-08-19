package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const syncProgressBatch = 25

func (uc *tiendanubeUseCase) emitSyncEvent(ctx context.Context, businessID, integrationID uint, eventType string, data map[string]interface{}) {
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
	cfg := domain.InventoryConfig{}
	if v, ok := config["inventory_sync_enabled"].(bool); ok {
		cfg.Enabled = v
	}
	cfg.SingleWarehouseID = invToUint(config["inventory_single_warehouse_id"])
	return cfg
}

func resolveWarehouseIDs(cfg domain.InventoryConfig) []uint {
	if cfg.SingleWarehouseID > 0 {
		return []uint{cfg.SingleWarehouseID}
	}
	return nil
}

func (uc *tiendanubeUseCase) SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	integration, cred, err := uc.resolveIntegrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	cfg := parseInventoryConfig(integration.Config)
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
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "tiendanube.inventory.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
	})

	fails := &failedSKUs{}
	updated := 0
	for i, item := range mapped {
		qty := stock[item.ProductID]
		action := "updated"

		productID, variantID, perr := parseExternalProductID(item.ExternalItemID)
		if perr == nil && variantID == 0 {
			variantID, _ = strconv.ParseInt(item.ExternalVariantID, 10, 64)
		}

		switch {
		case perr != nil || productID == 0:
			fails.add(item.SKU)
			action = "failed"
		case variantID == 0:
			target, terr := uc.client.ResolveStockTarget(ctx, cred, item.SKU)
			if terr != nil || target == nil || !target.Found {
				uc.logger.Warn(ctx).
					Str("sku", item.SKU).
					Str("external_product_id", item.ExternalItemID).
					Msg("No se pudo resolver la variante de Tiendanube para escribir stock")
				fails.add(item.SKU)
				action = "failed"
				break
			}
			if uerr := uc.client.SetVariantStock(ctx, cred, target.ProductID, target.VariantID, qty); uerr != nil {
				uc.logger.Error(ctx).Err(uerr).Str("sku", item.SKU).Msg("Error al actualizar stock en Tiendanube")
				fails.add(item.SKU)
				action = "failed"
			} else {
				updated++
			}
		default:
			if uerr := uc.client.SetVariantStock(ctx, cred, productID, variantID, qty); uerr != nil {
				uc.logger.Error(ctx).Err(uerr).
					Str("sku", item.SKU).
					Str("external_product_id", item.ExternalItemID).
					Msg("Error al actualizar stock en Tiendanube")
				fails.add(item.SKU)
				action = "failed"
			} else {
				updated++
			}
		}

		uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "tiendanube.inventory.sync.item", map[string]interface{}{
			"correlation_id": correlationID,
			"sku":            item.SKU,
			"quantity":       qty,
			"action":         action,
		})
		uc.maybeInventoryProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "tiendanube.inventory.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
		"updated":        updated,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *tiendanubeUseCase) maybeInventoryProgress(ctx context.Context, businessID, integrationID uint, correlationID string, processed, total, updated, failed int) {
	if processed%syncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "tiendanube.inventory.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"updated":        updated,
		"failed":         failed,
	})
}
