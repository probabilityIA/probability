package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) UpdateItemStock(ctx context.Context, integrationID, itemID string, quantity int) error {
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
