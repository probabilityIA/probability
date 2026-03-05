package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
)

func (uc *UseCase) UpdateRoute(ctx context.Context, dto dtos.UpdateRouteDTO) (*entities.Route, error) {
	existing, err := uc.repo.GetRouteByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.DriverID != nil {
		name, err := uc.repo.GetDriverNameByID(ctx, *dto.DriverID)
		if err != nil {
			return nil, fmt.Errorf("driver not found: %w", err)
		}
		existing.DriverName = name
	} else {
		existing.DriverName = ""
	}

	if dto.VehicleID != nil {
		plate, err := uc.repo.GetVehiclePlateByID(ctx, *dto.VehicleID)
		if err != nil {
			return nil, fmt.Errorf("vehicle not found: %w", err)
		}
		existing.VehiclePlate = plate
	} else {
		existing.VehiclePlate = ""
	}

	existing.DriverID = dto.DriverID
	existing.VehicleID = dto.VehicleID
	existing.Date = dto.Date
	existing.StartTime = dto.StartTime
	existing.EndTime = dto.EndTime
	existing.OriginWarehouseID = dto.OriginWarehouseID
	existing.OriginAddress = dto.OriginAddress
	existing.OriginLat = dto.OriginLat
	existing.OriginLng = dto.OriginLng
	existing.Notes = dto.Notes

	return uc.repo.UpdateRoute(ctx, existing)
}
