package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/errors"
)

func (uc *UseCase) CreateVehicle(ctx context.Context, dto dtos.CreateVehicleDTO) (*entities.Vehicle, error) {
	if dto.LicensePlate != "" {
		exists, err := uc.repo.ExistsByLicensePlate(ctx, dto.BusinessID, dto.LicensePlate, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domainerrors.ErrDuplicateLicensePlate
		}
	}

	status := dto.Status
	if status == "" {
		status = "active"
	}

	vehicle := &entities.Vehicle{
		BusinessID:         dto.BusinessID,
		Type:               dto.Type,
		LicensePlate:       dto.LicensePlate,
		Brand:              dto.Brand,
		VehicleModel:       dto.VehicleModel,
		Year:               dto.Year,
		Color:              dto.Color,
		Status:             status,
		WeightCapacityKg:   dto.WeightCapacityKg,
		VolumeCapacityM3:   dto.VolumeCapacityM3,
		PhotoURL:           dto.PhotoURL,
		InsuranceExpiry:    dto.InsuranceExpiry,
		RegistrationExpiry: dto.RegistrationExpiry,
	}

	return uc.repo.Create(ctx, vehicle)
}
