package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
)

func (uc *UseCase) ListMyOrders(ctx context.Context, businessID, userID uint, page, pageSize int) ([]entities.StorefrontOrder, int64, error) {
	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return nil, 0, err
	}
	if !active {
		return nil, 0, domainerrors.ErrStorefrontNotActive
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return uc.repo.ListOrdersByUserID(ctx, businessID, userID, page, pageSize)
}
