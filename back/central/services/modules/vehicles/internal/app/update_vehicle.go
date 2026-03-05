package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/errors"
)

func (uc *UseCase) UpdateVehicle(ctx context.Context, dto dtos.UpdateVehicleDTO) (*entities.Vehicle, error) {
	existing, err := uc.repo.GetByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.LicensePlate != "" && dto.LicensePlate != existing.LicensePlate {
		exists, err := uc.repo.ExistsByLicensePlate(ctx, dto.BusinessID, dto.LicensePlate, &dto.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domainerrors.ErrDuplicateLicensePlate
		}
	}

	existing.Type = dto.Type
	existing.LicensePlate = dto.LicensePlate
	existing.Brand = dto.Brand
	existing.VehicleModel = dto.VehicleModel
	existing.Year = dto.Year
	existing.Color = dto.Color
	existing.Status = dto.Status
	existing.WeightCapacityKg = dto.WeightCapacityKg
	existing.VolumeCapacityM3 = dto.VolumeCapacityM3
	existing.PhotoURL = dto.PhotoURL
	existing.InsuranceExpiry = dto.InsuranceExpiry
	existing.RegistrationExpiry = dto.RegistrationExpiry

	return uc.repo.Update(ctx, existing)
}
