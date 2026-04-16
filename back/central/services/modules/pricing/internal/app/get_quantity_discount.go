package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

func (uc *UseCase) GetQuantityDiscount(ctx context.Context, businessID, discountID uint) (*entities.QuantityDiscount, error) {
	return uc.repo.GetQuantityDiscount(ctx, businessID, discountID)
}
