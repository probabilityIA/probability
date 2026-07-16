package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) UpdateItemStock(ctx context.Context, integrationID, itemID string, quantity int) error {
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
	if err := uc.client.UpdateStock(ctx, accessToken, itemID, quantity); err != nil {
		if err == domain.ErrTokenExpired {
			newToken, rerr := uc.EnsureValidToken(ctx, integrationID)
			if rerr != nil {
				return rerr
			}
			return uc.client.UpdateStock(ctx, newToken, itemID, quantity)
		}
		return err
	}
	return nil
}
