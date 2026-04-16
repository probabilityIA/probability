package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
)

func (uc *UseCase) ListDriversForBusiness(ctx context.Context, businessID uint) ([]dtos.DriverOption, error) {
	return uc.repo.ListDriversForBusiness(ctx, businessID)
}

func (uc *UseCase) ListVehiclesForBusiness(ctx context.Context, businessID uint) ([]dtos.VehicleOption, error) {
	return uc.repo.ListVehiclesForBusiness(ctx, businessID)
}

func (uc *UseCase) ListAssignableOrders(ctx context.Context, businessID uint) ([]dtos.AssignableOrder, error) {
	return uc.repo.ListAssignableOrders(ctx, businessID)
}
