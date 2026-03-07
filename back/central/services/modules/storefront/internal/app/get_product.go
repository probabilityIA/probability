package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
)

func (uc *UseCase) GetProduct(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error) {
	return uc.repo.GetProductByID(ctx, businessID, productID)
}
