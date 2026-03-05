package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/entities"
)

func (uc *UseCase) GetVehicle(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
	return uc.repo.GetByID(ctx, businessID, vehicleID)
}
