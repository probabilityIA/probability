package app

import "context"

func (uc *UseCase) DeleteVehicle(ctx context.Context, businessID, vehicleID uint) error {
	_, err := uc.repo.GetByID(ctx, businessID, vehicleID)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, businessID, vehicleID)
}
