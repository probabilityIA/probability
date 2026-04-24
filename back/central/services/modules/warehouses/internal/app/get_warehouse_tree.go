package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
)

func (u *UseCase) GetWarehouseTree(ctx context.Context, businessID, warehouseID uint) (*dtos.WarehouseTreeDTO, error) {
	exists, err := u.repo.WarehouseExists(ctx, businessID, warehouseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainerrors.ErrWarehouseNotFound
	}
	return u.repo.GetWarehouseTree(ctx, businessID, warehouseID)
}
