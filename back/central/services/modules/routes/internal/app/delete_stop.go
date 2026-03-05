package app

import (
	"context"
)

func (uc *UseCase) DeleteStop(ctx context.Context, businessID, routeID, stopID uint) error {
	_, err := uc.repo.GetRouteByID(ctx, businessID, routeID)
	if err != nil {
		return err
	}

	stop, err := uc.repo.GetStopByID(ctx, routeID, stopID)
	if err != nil {
		return err
	}

	// Clear driver info from order
	if stop.OrderID != nil {
		_ = uc.repo.ClearOrderDriverInfo(ctx, *stop.OrderID)
	}

	if err := uc.repo.DeleteStop(ctx, routeID, stopID); err != nil {
		return err
	}

	_ = uc.repo.UpdateRouteCounters(ctx, routeID)
	return nil
}
