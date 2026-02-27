package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) GetOrderStatus(ctx context.Context, id uint) (*entities.OrderStatusInfo, error) {
	return uc.repo.GetOrderStatusByID(ctx, id)
}
