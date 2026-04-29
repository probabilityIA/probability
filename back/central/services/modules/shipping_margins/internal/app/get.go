package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
)

func (uc *UseCase) Get(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
	return uc.repo.GetByID(ctx, businessID, id)
}
