package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) ListOrderStatuses(ctx context.Context, isActive *bool) ([]entities.OrderStatusInfo, error) {
	return uc.repo.ListOrderStatuses(ctx, isActive)
}
