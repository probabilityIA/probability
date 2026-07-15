package usecases

import (
	"context"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) processItemNotification(ctx context.Context, notification *domain.MeliNotification) error {
	itemID := extractResourceStringID(notification.Resource)
	if itemID == "" {
		return nil
	}
	integration, accessToken, err := uc.resolveIntegrationAndToken(ctx, notification.UserID)
	if err != nil {
		return err
	}
	return uc.reconcileItemStock(ctx, integration, accessToken, itemID)
}

func (uc *meliUseCase) reconcileItemStock(ctx context.Context, integration *domain.Integration, accessToken, itemID string) error {
	cfg := parseInventoryConfig(integration.Config)
	if !cfg.Enabled {
		return nil
	}

	mapped, err := uc.inventoryRepo.ListMappedItems(ctx, integration.ID)
	if err != nil {
		return err
	}

	var productID string
	for _, m := range mapped {
		if m.ExternalItemID == itemID {
			productID = m.ProductID
			break
		}
	}
	if productID == "" {
		uc.logger.Info(ctx).Str("item_id", itemID).Msg("Item change ignored: no local product mapped")
		return nil
	}

	warehouseIDs := resolveWarehouseIDs(cfg)
	stock, err := uc.inventoryRepo.GetStockForProducts(ctx, []string{productID}, warehouseIDs)
	if err != nil {
		return err
	}
	qty := stock[productID]

	item, err := uc.client.GetItem(ctx, accessToken, itemID)
	if err == nil && item.AvailableQuantity == qty && len(item.Variations) == 0 {
		return nil
	}

	if uerr := uc.client.UpdateStock(ctx, accessToken, itemID, qty); uerr != nil {
		if uerr == domain.ErrTokenExpired {
			newToken, rerr := uc.EnsureValidToken(ctx, strconv.FormatUint(uint64(integration.ID), 10))
			if rerr == nil {
				return uc.client.UpdateStock(ctx, newToken, itemID, qty)
			}
		}
		return uerr
	}
	return nil
}
