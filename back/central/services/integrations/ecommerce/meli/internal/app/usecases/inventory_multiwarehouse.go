package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) syncInventoryMultiWarehouse(ctx context.Context, businessID, integrationID uint, integrationIDStr, accessToken string, cfg domain.InventoryConfig, mapped []domain.MappedItem, correlationID string) error {
	warehouseIDs := make([]uint, 0, len(cfg.WarehouseMappings))
	for _, wm := range cfg.WarehouseMappings {
		warehouseIDs = append(warehouseIDs, wm.InternalWarehouseID)
	}

	productIDs := make([]string, 0, len(mapped))
	for _, m := range mapped {
		productIDs = append(productIDs, m.ProductID)
	}

	byWarehouse, err := uc.inventoryRepo.GetInventoryByWarehouses(ctx, productIDs, warehouseIDs)
	if err != nil {
		return err
	}

	total := len(mapped)
	uc.emitInventoryEvent(ctx, businessID, integrationID, "meli.inventory.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
		"mode":           "multi_warehouse",
	})

	updated, unchanged, skipped, failed := 0, 0, 0, 0
	for i, m := range mapped {
		stockByWh := byWarehouse[m.ProductID]
		skip, perr := uc.pushItemMultiWarehouse(ctx, accessToken, m, cfg.WarehouseMappings, stockByWh)
		if perr == domain.ErrTokenExpired {
			if newToken, rerr := uc.EnsureValidToken(ctx, integrationIDStr); rerr == nil {
				accessToken = newToken
				skip, perr = uc.pushItemMultiWarehouse(ctx, accessToken, m, cfg.WarehouseMappings, stockByWh)
			}
		}
		switch {
		case perr != nil:
			uc.logger.Error(ctx).Err(perr).Str("sku", m.SKU).Str("item_id", m.ExternalItemID).Msg("Error al actualizar stock multi-bodega en MercadoLibre")
			failed++
		case skip:
			skipped++
		default:
			updated++
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

func (uc *meliUseCase) pushItemMultiWarehouse(ctx context.Context, accessToken string, m domain.MappedItem, mappings []domain.WarehouseMapping, stockByWh map[uint]int) (bool, error) {
	item, err := uc.client.GetItem(ctx, accessToken, m.ExternalItemID)
	if err != nil {
		return false, err
	}

	userProductID := resolveUserProductID(item, m.SKU)
	if userProductID == "" {
		total := 0
		for _, wm := range mappings {
			total += stockByWh[wm.InternalWarehouseID]
		}
		return false, uc.client.UpdateStock(ctx, accessToken, m.ExternalItemID, total)
	}

	locations := make([]domain.StockLocation, 0, len(mappings))
	for _, wm := range mappings {
		locations = append(locations, domain.StockLocation{
			Type:          "seller_warehouse",
			StoreID:       wm.MLStoreID,
			NetworkNodeID: wm.MLNetworkNodeID,
			Quantity:      stockByWh[wm.InternalWarehouseID],
		})
	}

	current, err := uc.client.GetUserProductStock(ctx, accessToken, userProductID)
	if err != nil {
		return false, err
	}

	return false, uc.client.UpdateUserProductStock(ctx, accessToken, userProductID, current.Version, locations)
}

func resolveUserProductID(item *domain.MeliItemDetail, sku string) string {
	if len(item.Variations) > 0 {
		for _, v := range item.Variations {
			if v.SellerSKU == sku && v.UserProductID != "" {
				return v.UserProductID
			}
		}
		for _, v := range item.Variations {
			if v.UserProductID != "" {
				return v.UserProductID
			}
		}
	}
	return item.UserProductID
}
