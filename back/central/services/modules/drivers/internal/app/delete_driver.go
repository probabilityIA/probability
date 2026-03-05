package app

import "context"

func (uc *UseCase) DeleteDriver(ctx context.Context, businessID, driverID uint) error {
	_, err := uc.repo.GetByID(ctx, businessID, driverID)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, businessID, driverID)
}
