package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
)

func (uc *UseCase) AddStop(ctx context.Context, dto dtos.AddStopDTO) (*entities.RouteStop, error) {
	route, err := uc.repo.GetRouteByID(ctx, dto.BusinessID, dto.RouteID)
	if err != nil {
		return nil, err
	}

	stop := &entities.RouteStop{
		RouteID:       dto.RouteID,
		OrderID:       dto.OrderID,
		Sequence:      route.TotalStops + 1,
		Status:        "pending",
		Address:       dto.Address,
		City:          dto.City,
		Lat:           dto.Lat,
		Lng:           dto.Lng,
		CustomerName:  dto.CustomerName,
		CustomerPhone: dto.CustomerPhone,
		DeliveryNotes: dto.DeliveryNotes,
	}

	created, err := uc.repo.AddStop(ctx, stop)
	if err != nil {
		return nil, err
	}

	_ = uc.repo.UpdateRouteCounters(ctx, dto.RouteID)

	// Set driver info on order if route has driver
	if dto.OrderID != nil && route.DriverID != nil {
		_ = uc.repo.UpdateOrderDriverInfo(ctx, *dto.OrderID, route.DriverID, route.DriverName, true)
	}

	return created, nil
}
