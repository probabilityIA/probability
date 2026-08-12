package usecases

import (
	"context"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) UpdateItemStock(ctx context.Context, integrationID, productID, itemID, variantID string, quantity int) error {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return err
	}
	if integration == nil {
		return nil
	}
	if enabled, _ := integration.Config["inventory_sync_enabled"].(bool); !enabled {
		uc.logger.Info(ctx).Str("integration_id", integrationID).Msg("Sync de inventario desactivado para MercadoLibre, push omitido")
		return nil
	}

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	if productID != "" {
		state, serr := uc.inventoryRepo.GetPushState(ctx, uint(integIDUint), productID, variantID)
		if serr != nil {
			uc.logger.Warn(ctx).Err(serr).Str("product_id", productID).Msg("No se pudo leer el estado del push, se empuja igual")
		} else if state != nil {
			if state.LogisticType == domain.LogisticFulfillment {
				uc.logger.Info(ctx).Str("item_id", itemID).Msg("Publicacion de fulfillment, el stock lo administra MercadoLibre, push omitido")
				return nil
			}
			if state.LastPushedQty != nil && *state.LastPushedQty == quantity {
				return nil
			}
		}
	}

	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return err
	}
	cli := uc.clientFor(ctx, integration)

	if err := cli.UpdateStock(ctx, accessToken, itemID, variantID, quantity); err != nil {
		if err != domain.ErrTokenExpired {
			return err
		}
		newToken, rerr := uc.EnsureValidToken(ctx, integrationID)
		if rerr != nil {
			return rerr
		}
		if err := cli.UpdateStock(ctx, newToken, itemID, variantID, quantity); err != nil {
			return err
		}
	}

	if productID != "" {
		if merr := uc.inventoryRepo.MarkPushedQty(ctx, uint(integIDUint), productID, variantID, quantity); merr != nil {
			uc.logger.Warn(ctx).Err(merr).Str("product_id", productID).Msg("No se pudo registrar la cantidad empujada")
		}
	}
	return nil
}
