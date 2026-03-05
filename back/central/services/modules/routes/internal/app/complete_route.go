package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
)

func (uc *UseCase) CompleteRoute(ctx context.Context, businessID, routeID uint) error {
	route, err := uc.repo.GetRouteByID(ctx, businessID, routeID)
	if err != nil {
		return err
	}

	if route.Status != "in_progress" {
		return domainerrors.ErrRouteNotInProgress
	}

	// Mark remaining pending stops as skipped
	_ = uc.repo.SetPendingStopsStatus(ctx, routeID, "skipped")

	if err := uc.repo.UpdateRouteStatus(ctx, routeID, "completed"); err != nil {
		return err
	}

	if err := uc.repo.SetRouteActualEnd(ctx, routeID); err != nil {
		return err
	}

	if err := uc.repo.UpdateRouteCounters(ctx, routeID); err != nil {
		return err
	}

	// Set driver back to active
	if route.DriverID != nil {
		_ = uc.repo.UpdateDriverStatus(ctx, *route.DriverID, "active")
	}

	return nil
}
