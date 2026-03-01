package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (uc *useCase) ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}
	return uc.repo.ListMovementTypes(ctx, params)
}
