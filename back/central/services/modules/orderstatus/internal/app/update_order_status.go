package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) UpdateOrderStatus(ctx context.Context, id uint, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error) {
	return uc.repo.UpdateOrderStatus(ctx, id, status)
}
