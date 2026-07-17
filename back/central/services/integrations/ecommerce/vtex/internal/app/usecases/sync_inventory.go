package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

func (uc *vtexUseCase) GetWarehouses(ctx context.Context, integrationID string, businessID uint) (*domain.WarehousesInfo, error) {
	integration, err := uc.integrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	cred, err := uc.resolveCredential(ctx, integration, integrationID)
	if err != nil {
		return nil, err
	}

	warehouses, err := uc.client.GetWarehouses(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("listing vtex warehouses: %w", err)
	}

	return &domain.WarehousesInfo{Warehouses: warehouses}, nil
}

func (uc *vtexUseCase) SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integration, err := uc.integrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	cred, err := uc.resolveCredential(ctx, integration, integrationID)
	if err != nil {
		return err
	}

	invCfg := uc.inventoryConfigFrom(integration.Config)
	if !invCfg.Enabled {
		return domain.ErrInventorySyncDisabled
	}

	groups := invCfg.WarehouseGroups()
	if len(groups) == 0 {
		return domain.ErrNoWarehousesMapped
	}

	mapped, err := uc.productRepo.ListMappedItems(ctx, integration.ID)
	if err != nil {
		return fmt.Errorf("listing mapped items: %w", err)
	}
	if len(mapped) == 0 {
		uc.emitSyncEvent(ctx, businessID, integration.ID, "vtex.inventory.sync.completed", map[string]interface{}{
			"correlation_id": correlationID,
			"total":          0,
			"synced":         0,
			"failed":         0,
		})
		return nil
	}

	productIDs := make([]string, 0, len(mapped))
	for _, m := range mapped {
		productIDs = append(productIDs, m.ProductID)
	}

	total := len(mapped) * len(groups)

	uc.emitSyncEvent(ctx, businessID, integration.ID, "vtex.inventory.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
	})

	synced := 0
	processed := 0
	fails := &failedSKUs{}

	for vtexWarehouseID, internalWarehouseIDs := range groups {
		stock, err := uc.productRepo.GetStockForProducts(ctx, productIDs, internalWarehouseIDs)
		if err != nil {
			uc.logger.Error(ctx).Err(err).
				Str("vtex_warehouse", vtexWarehouseID).
				Msg("Error obteniendo stock para bodega VTEX")
			for range mapped {
				processed++
				fails.add("")
			}
			continue
		}

		for _, item := range mapped {
			processed++
			quantity := stock[item.ProductID]

			if err := uc.client.UpdateSKUInventory(ctx, cred, item.ExternalItemID, vtexWarehouseID, quantity); err != nil {
				uc.logger.Error(ctx).Err(err).
					Str("sku", item.SKU).
					Str("vtex_sku_id", item.ExternalItemID).
					Str("vtex_warehouse", vtexWarehouseID).
					Msg("Error actualizando inventario en VTEX")
				fails.add(item.SKU)
				continue
			}

			synced++
			uc.maybeProductProgress(ctx, businessID, integration.ID, correlationID, domain.DirectionToVTEX, processed, total, synced, 0, fails.count())
		}
	}

	uc.emitSyncEvent(ctx, businessID, integration.ID, "vtex.inventory.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
		"synced":         synced,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	uc.logger.Info(ctx).
		Int("total", total).
		Int("synced", synced).
		Int("failed", fails.count()).
		Uint("integration_id", integration.ID).
		Msg("VTEX inventory sync completed")

	if synced == 0 && fails.count() > 0 {
		return fmt.Errorf("vtex: no se pudo sincronizar inventario de ningun producto")
	}

	return nil
}
