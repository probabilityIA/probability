package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
)

func (uc *UseCase) CreateRoute(ctx context.Context, dto dtos.CreateRouteDTO) (*entities.Route, error) {
	var driverName string
	if dto.DriverID != nil {
		name, err := uc.repo.GetDriverNameByID(ctx, *dto.DriverID)
		if err != nil {
			return nil, fmt.Errorf("driver not found: %w", err)
		}
		driverName = name
	}

	var vehiclePlate string
	if dto.VehicleID != nil {
		plate, err := uc.repo.GetVehiclePlateByID(ctx, *dto.VehicleID)
		if err != nil {
			return nil, fmt.Errorf("vehicle not found: %w", err)
		}
		vehiclePlate = plate
	}

	route := &entities.Route{
		BusinessID:        dto.BusinessID,
		DriverID:          dto.DriverID,
		VehicleID:         dto.VehicleID,
		Status:            "planned",
		Date:              dto.Date,
		StartTime:         dto.StartTime,
		EndTime:           dto.EndTime,
		OriginWarehouseID: dto.OriginWarehouseID,
		OriginAddress:     dto.OriginAddress,
		OriginLat:         dto.OriginLat,
		OriginLng:         dto.OriginLng,
		TotalStops:        len(dto.Stops),
		Notes:             dto.Notes,
		DriverName:        driverName,
		VehiclePlate:      vehiclePlate,
	}

	stops := make([]entities.RouteStop, len(dto.Stops))
	for i, s := range dto.Stops {
		stops[i] = entities.RouteStop{
			Sequence:      i + 1,
			OrderID:       s.OrderID,
			Status:        "pending",
			Address:       s.Address,
			City:          s.City,
			Lat:           s.Lat,
			Lng:           s.Lng,
			CustomerName:  s.CustomerName,
			CustomerPhone: s.CustomerPhone,
			DeliveryNotes: s.DeliveryNotes,
		}
	}

	created, err := uc.repo.CreateRoute(ctx, route, stops)
	if err != nil {
		return nil, err
	}

	// Set driver info on orders
	if dto.DriverID != nil {
		for _, s := range dto.Stops {
			if s.OrderID != nil {
				_ = uc.repo.UpdateOrderDriverInfo(ctx, *s.OrderID, dto.DriverID, driverName, true)
			}
		}
	}

	return created, nil
}
