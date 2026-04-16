package app

import (
	"context"
)

func (uc *UseCase) DeleteRoute(ctx context.Context, businessID, routeID uint) error {
	route, err := uc.repo.GetRouteByID(ctx, businessID, routeID)
	if err != nil {
		return err
	}

	// Clear driver info from orders
	for _, stop := range route.Stops {
		if stop.OrderID != nil {
			_ = uc.repo.ClearOrderDriverInfo(ctx, *stop.OrderID)
		}
	}

	// Set driver back to active if assigned
	if route.DriverID != nil {
		_ = uc.repo.UpdateDriverStatus(ctx, *route.DriverID, "active")
	}

	return uc.repo.DeleteRoute(ctx, businessID, routeID)
}
