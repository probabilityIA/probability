package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
)

func (u *UseCase) GetOccupancy(ctx context.Context, businessID, warehouseID uint) ([]entities.OccupancyItem, error) {
	exists, err := u.repo.WarehouseExists(ctx, businessID, warehouseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainerrors.ErrWarehouseNotFound
	}
	return u.repo.GetWarehouseOccupancy(ctx, businessID, warehouseID)
}
