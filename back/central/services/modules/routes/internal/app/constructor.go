package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/ports"
)

type IUseCase interface {
	CreateRoute(ctx context.Context, dto dtos.CreateRouteDTO) (*entities.Route, error)
	GetRoute(ctx context.Context, businessID, routeID uint) (*entities.Route, error)
	ListRoutes(ctx context.Context, params dtos.ListRoutesParams) ([]entities.Route, int64, error)
	UpdateRoute(ctx context.Context, dto dtos.UpdateRouteDTO) (*entities.Route, error)
	DeleteRoute(ctx context.Context, businessID, routeID uint) error
	StartRoute(ctx context.Context, businessID, routeID uint) error
	CompleteRoute(ctx context.Context, businessID, routeID uint) error
	AddStop(ctx context.Context, dto dtos.AddStopDTO) (*entities.RouteStop, error)
	UpdateStop(ctx context.Context, dto dtos.UpdateStopDTO) (*entities.RouteStop, error)
	DeleteStop(ctx context.Context, businessID, routeID, stopID uint) error
	UpdateStopStatus(ctx context.Context, dto dtos.UpdateStopStatusDTO) error
	ReorderStops(ctx context.Context, dto dtos.ReorderStopsDTO) error
	ListDriversForBusiness(ctx context.Context, businessID uint) ([]dtos.DriverOption, error)
	ListVehiclesForBusiness(ctx context.Context, businessID uint) ([]dtos.VehicleOption, error)
	ListAssignableOrders(ctx context.Context, businessID uint) ([]dtos.AssignableOrder, error)
}

type UseCase struct {
	repo ports.IRepository
}

func New(repo ports.IRepository) IUseCase {
	return &UseCase{repo: repo}
}
