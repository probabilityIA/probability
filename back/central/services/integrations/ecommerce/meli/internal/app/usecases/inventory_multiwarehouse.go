package usecases

import (
	"context"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

func (uc *meliUseCase) syncInventoryMultiWarehouse(ctx context.Context, cli domain.IMeliClient, businessID, integrationID uint, integrationIDStr, accessToken string, cfg domain.InventoryConfig, mapped []domain.MappedItem, rules []productmatch.Rule, correlationID string) error {
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
		skip, perr := uc.pushItemMultiWarehouse(ctx, cli, accessToken, m, cfg.WarehouseMappings, stockByWh, rules)
		if perr == domain.ErrTokenExpired {
			if newToken, rerr := uc.EnsureValidToken(ctx, integrationIDStr); rerr == nil {
				accessToken = newToken
				skip, perr = uc.pushItemMultiWarehouse(ctx, cli, accessToken, m, cfg.WarehouseMappings, stockByWh, rules)
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

func (uc *meliUseCase) pushItemMultiWarehouse(ctx context.Context, cli domain.IMeliClient, accessToken string, m domain.MappedItem, mappings []domain.WarehouseMapping, stockByWh map[uint]int, rules []productmatch.Rule) (bool, error) {
	item, err := cli.GetItem(ctx, accessToken, m.ExternalItemID)
	if err != nil {
		return false, err
	}

	userProductID := resolveUserProductID(item, m, rules)
	if userProductID == "" {
		total := 0
		for _, wm := range mappings {
			total += stockByWh[wm.InternalWarehouseID]
		}
		return false, cli.UpdateStock(ctx, accessToken, m.ExternalItemID, m.ExternalVariantID, total)
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

	current, err := cli.GetUserProductStock(ctx, accessToken, userProductID)
	if err != nil {
		return false, err
	}

	return false, cli.UpdateUserProductStock(ctx, accessToken, userProductID, current.Version, locations)
}

func resolveUserProductID(item *domain.MeliItemDetail, m domain.MappedItem, rules []productmatch.Rule) string {
	if len(item.Variations) > 0 {
		if id := matchVariationByExternalID(item, m.ExternalVariantID); id != "" {
			return id
		}
		if id := matchVariationByRules(item, m, rules); id != "" {
			return id
		}
		for _, v := range item.Variations {
			if v.UserProductID != "" {
				return v.UserProductID
			}
		}
	}
	return item.UserProductID
}

func matchVariationByExternalID(item *domain.MeliItemDetail, externalVariantID string) string {
	key := productmatch.Normalize(externalVariantID)
	if key == "" {
		return ""
	}
	for _, v := range item.Variations {
		if productmatch.Normalize(strconv.FormatInt(v.ID, 10)) == key && v.UserProductID != "" {
			return v.UserProductID
		}
	}
	return ""
}

func matchVariationByRules(item *domain.MeliItemDetail, m domain.MappedItem, rules []productmatch.Rule) string {
	probe := productmatch.Item{
		SKU:        firstNonEmpty(m.ExternalSKU, m.SKU),
		Barcode:    firstNonEmpty(m.ExternalBarcode, m.Barcode),
		ExternalID: m.ExternalItemID,
		VariantID:  m.ExternalVariantID,
	}
	values := make([]productmatch.Values, 0, len(item.Variations))
	for _, v := range item.Variations {
		values = append(values, productmatch.Item{
			SKU:        v.SellerSKU,
			ExternalID: item.ID,
			VariantID:  strconv.FormatInt(v.ID, 10),
		}.Values())
	}
	index := productmatch.IndexChannel(rules, values)
	pos, _, ok := index.Find(probe.Values())
	if !ok || item.Variations[pos].UserProductID == "" {
		return ""
	}
	return item.Variations[pos].UserProductID
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
