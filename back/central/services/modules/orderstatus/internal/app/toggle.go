package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) ToggleOrderStatusMappingActive(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
	return uc.repo.ToggleActive(ctx, id)
}
