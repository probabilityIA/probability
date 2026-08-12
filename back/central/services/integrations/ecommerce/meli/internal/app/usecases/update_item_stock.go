package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) UpdateItemStock(ctx context.Context, integrationID, itemID, variantID string, quantity int) error {
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

	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return err
	}
	cli := uc.clientFor(ctx, integration)
	if err := cli.UpdateStock(ctx, accessToken, itemID, variantID, quantity); err != nil {
		if err == domain.ErrTokenExpired {
			newToken, rerr := uc.EnsureValidToken(ctx, integrationID)
			if rerr != nil {
				return rerr
			}
			return cli.UpdateStock(ctx, newToken, itemID, variantID, quantity)
		}
		return err
	}
	return nil
}
