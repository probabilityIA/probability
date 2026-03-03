package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	// Verificar que la bodega existe
	exists, err := uc.repo.WarehouseExists(ctx, params.WarehouseID, params.BusinessID)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, domainerrors.ErrWarehouseNotFound
	}

	return uc.repo.ListWarehouseInventory(ctx, params)
}
