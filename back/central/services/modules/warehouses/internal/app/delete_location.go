package app

import "context"

func (uc *UseCase) DeleteLocation(ctx context.Context, warehouseID, locationID uint, businessID uint) error {
	// Verificar que la bodega pertenece al negocio
	_, err := uc.repo.GetByID(ctx, businessID, warehouseID)
	if err != nil {
		return err
	}

	return uc.repo.DeleteLocation(ctx, warehouseID, locationID)
}
