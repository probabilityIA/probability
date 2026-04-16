package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) CreateOrderStatus(ctx context.Context, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error) {
	return uc.repo.CreateOrderStatus(ctx, status)
}
