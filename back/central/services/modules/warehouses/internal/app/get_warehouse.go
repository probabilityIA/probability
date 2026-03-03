package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func (uc *UseCase) GetWarehouse(ctx context.Context, businessID, warehouseID uint) (*entities.Warehouse, error) {
	return uc.repo.GetByID(ctx, businessID, warehouseID)
}
