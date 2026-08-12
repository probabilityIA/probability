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

	matches := make([]domain.MappedItem, 0)
	for _, m := range mapped {
		if m.ExternalItemID == itemID {
			matches = append(matches, m)
		}
	}
	if len(matches) == 0 {
		uc.logger.Info(ctx).Str("item_id", itemID).Msg("Item change ignored: no local product mapped")
		return nil
	}

	warehouseIDs := resolveWarehouseIDs(cfg)
	productIDs := make([]string, len(matches))
	for i, m := range matches {
		productIDs[i] = m.ProductID
	}
	stock, err := uc.inventoryRepo.GetStockForProducts(ctx, productIDs, warehouseIDs)
	if err != nil {
		return err
	}

	cli := uc.clientFor(ctx, integration)
	item, itemErr := cli.GetItem(ctx, accessToken, itemID)

	if itemErr == nil && item.LogisticType != "" {
		if uerr := uc.inventoryRepo.UpdateLogisticType(ctx, integration.ID, itemID, item.LogisticType); uerr != nil {
			uc.logger.Warn(ctx).Err(uerr).Str("item_id", itemID).Msg("No se pudo refrescar el tipo logistico de la publicacion")
		}
	}

	for _, m := range matches {
		qty := stock[m.ProductID]

		if itemErr == nil && len(item.Variations) == 0 && item.AvailableQuantity == qty {
			continue
		}
		if itemErr == nil && m.ExternalVariantID != "" {
			if v := findVariation(item, m.ExternalVariantID); v != nil && v.AvailableQuantity == qty {
				continue
			}
		}

		if uerr := cli.UpdateStock(ctx, accessToken, itemID, m.ExternalVariantID, qty); uerr != nil {
			if uerr == domain.ErrTokenExpired {
				newToken, rerr := uc.EnsureValidToken(ctx, strconv.FormatUint(uint64(integration.ID), 10))
				if rerr != nil {
					return rerr
				}
				accessToken = newToken
				if uerr = cli.UpdateStock(ctx, accessToken, itemID, m.ExternalVariantID, qty); uerr != nil {
					return uerr
				}
				continue
			}
			return uerr
		}
	}
	return nil
}

func findVariation(item *domain.MeliItemDetail, variantID string) *domain.MeliItemVariation {
	for i := range item.Variations {
		if strconv.FormatInt(item.Variations[i].ID, 10) == variantID {
			return &item.Variations[i]
		}
	}
	return nil
}
