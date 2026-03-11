package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
)

func (uc *UseCase) GetMyOrder(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error) {
	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, domainerrors.ErrStorefrontNotActive
	}

	return uc.repo.GetOrderByIDAndUserID(ctx, orderID, businessID, userID)
}
