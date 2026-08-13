package usecases

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
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

func toStr(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case int:
		return strconv.Itoa(val)
	}
	return ""
}

func parseWarehouseMappings(config map[string]interface{}) []domain.WarehouseMapping {
	raw, ok := config["warehouse_mappings"].([]interface{})
	if !ok {
		return nil
	}
	var mappings []domain.WarehouseMapping
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		internalID := toUint(m["internal_warehouse_id"])
		storeID := toStr(m["ml_store_id"])
		if internalID == 0 || storeID == "" {
			continue
		}
		mappings = append(mappings, domain.WarehouseMapping{
			InternalWarehouseID: internalID,
			MLStoreID:           storeID,
			MLNetworkNodeID:     toStr(m["ml_network_node_id"]),
		})
	}
	return mappings
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
	cfg.WarehouseMappings = parseWarehouseMappings(config)
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

func filtrarPorSKU(mapped []domain.MappedItem, skus []string) []domain.MappedItem {
	if len(skus) == 0 {
		return mapped
	}
	buscados := make(map[string]bool, len(skus))
	for _, sku := range skus {
		limpio := strings.ToUpper(strings.TrimSpace(sku))
		if limpio != "" {
			buscados[limpio] = true
		}
	}
	if len(buscados) == 0 {
		return mapped
	}
	filtrados := make([]domain.MappedItem, 0, len(buscados))
	for _, m := range mapped {
		if buscados[strings.ToUpper(strings.TrimSpace(m.SKU))] {
			filtrados = append(filtrados, m)
		}
	}
	return filtrados
}

func (uc *meliUseCase) SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return domain.ErrIntegrationNotFound
	}

	cfg := parseInventoryConfig(integration.Config)

	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return err
	}

	mapped, err := uc.inventoryRepo.ListMappedItems(ctx, uint(integIDUint))
	if err != nil {
		return fmt.Errorf("listing mapped items: %w", err)
	}
	mapped = filtrarPorSKU(mapped, skus)

	cli := uc.clientFor(ctx, integration)
	rules := productmatch.Sanitize(integration.ProductMatchRules)

	if len(cfg.WarehouseMappings) > 0 {
		return uc.syncInventoryMultiWarehouse(ctx, cli, businessID, uint(integIDUint), integrationID, accessToken, cfg, mapped, rules, correlationID)
	}

	return uc.syncInventorySingle(ctx, cli, businessID, uint(integIDUint), integrationID, accessToken, cfg, mapped, correlationID, len(skus) > 0)
}

func (uc *meliUseCase) syncInventorySingle(ctx context.Context, cli domain.IMeliClient, businessID, integrationID uint, integrationIDStr, accessToken string, cfg domain.InventoryConfig, mapped []domain.MappedItem, correlationID string, seleccionManual bool) error {
	warehouseIDs := resolveWarehouseIDs(cfg)

	productIDs := make([]string, 0, len(mapped))
	for _, m := range mapped {
		productIDs = append(productIDs, m.ProductID)
	}
	stock, err := uc.inventoryRepo.GetStockForProducts(ctx, productIDs, warehouseIDs)
	if err != nil {
		return fmt.Errorf("getting stock: %w", err)
	}

	total := len(mapped)
	uc.emitInventoryEvent(ctx, businessID, integrationID, "meli.inventory.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
	})

	updated, unchanged, skipped, failed := 0, 0, 0, 0
	for i, m := range mapped {
		if m.LogisticType == domain.LogisticFulfillment {
			skipped++
			uc.emitInventoryEvent(ctx, businessID, integrationID, "meli.inventory.sync.item", map[string]interface{}{
				"correlation_id": correlationID,
				"sku":            m.SKU,
				"action":         "skipped",
				"error":          "publicacion de fulfillment, el stock lo administra MercadoLibre",
			})
			uc.maybeInventoryProgress(ctx, businessID, integrationID, correlationID, i+1, total, updated, unchanged, skipped, failed)
			continue
		}
		qty, hasLevel := stock[m.ProductID]
		if hasLevel && m.LastPushedQty != nil && *m.LastPushedQty == qty {
			unchanged++
			uc.maybeInventoryProgress(ctx, businessID, integrationID, correlationID, i+1, total, updated, unchanged, skipped, failed)
			continue
		}
		if !hasLevel {
			if !seleccionManual {
				skipped++
				uc.emitInventoryEvent(ctx, businessID, integrationID, "meli.inventory.sync.item", map[string]interface{}{
					"correlation_id": correlationID,
					"sku":            m.SKU,
					"action":         "skipped",
					"error":          "el producto no tiene registro de inventario, se conserva el stock actual en MercadoLibre",
				})
				uc.maybeInventoryProgress(ctx, businessID, integrationID, correlationID, i+1, total, updated, unchanged, skipped, failed)
				continue
			}
			qty = 0
		}
		if uerr := cli.UpdateStock(ctx, accessToken, m.ExternalItemID, m.ExternalVariantID, qty); uerr != nil {
			if uerr == domain.ErrTokenExpired {
				newToken, rerr := uc.EnsureValidToken(ctx, integrationIDStr)
				if rerr == nil {
					accessToken = newToken
					if retry := cli.UpdateStock(ctx, accessToken, m.ExternalItemID, m.ExternalVariantID, qty); retry == nil {
						updated++
						if merr := uc.inventoryRepo.MarkPushedQty(ctx, integrationID, m.ProductID, m.ExternalVariantID, qty); merr != nil {
							uc.logger.Warn(ctx).Err(merr).Str("sku", m.SKU).Msg("No se pudo registrar la cantidad empujada")
						}
						uc.maybeInventoryProgress(ctx, businessID, integrationID, correlationID, i+1, total, updated, unchanged, skipped, failed)
						continue
					}
				}
			}
			uc.logger.Error(ctx).Err(uerr).Str("sku", m.SKU).Str("item_id", m.ExternalItemID).Msg("Error al actualizar stock en MercadoLibre")
			failed++
			uc.emitInventoryEvent(ctx, businessID, integrationID, "meli.inventory.sync.item", map[string]interface{}{
				"correlation_id": correlationID,
				"sku":            m.SKU,
				"quantity":       qty,
				"action":         "failed",
				"error":          uerr.Error(),
			})
		} else {
			updated++
			if merr := uc.inventoryRepo.MarkPushedQty(ctx, integrationID, m.ProductID, m.ExternalVariantID, qty); merr != nil {
				uc.logger.Warn(ctx).Err(merr).Str("sku", m.SKU).Msg("No se pudo registrar la cantidad empujada")
			}
			uc.emitInventoryEvent(ctx, businessID, integrationID, "meli.inventory.sync.item", map[string]interface{}{
				"correlation_id": correlationID,
				"sku":            m.SKU,
				"quantity":       qty,
				"action":         "updated",
			})
		}
		uc.maybeInventoryProgress(ctx, businessID, integrationID, correlationID, i+1, total, updated, unchanged, skipped, failed)
	}

	uc.emitInventoryEvent(ctx, businessID, integrationID, "meli.inventory.sync.completed", map[string]interface{}{
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
