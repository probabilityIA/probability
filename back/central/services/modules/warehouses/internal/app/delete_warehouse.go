package app

import "context"

func (uc *UseCase) DeleteWarehouse(ctx context.Context, businessID, warehouseID uint) error {
	return uc.repo.Delete(ctx, businessID, warehouseID)
}
