package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
)

func (uc *UseCase) GetMyOrder(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error) {
	return uc.repo.GetOrderByIDAndUserID(ctx, orderID, businessID, userID)
}
