package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) ListOrderStatusMappings(ctx context.Context, filters map[string]interface{}) ([]entities.OrderStatusMapping, int64, error) {
	return uc.repo.List(ctx, filters)
}
