package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
)

type IRepository interface {
	// Route CRUD
	CreateRoute(ctx context.Context, route *entities.Route, stops []entities.RouteStop) (*entities.Route, error)
	GetRouteByID(ctx context.Context, businessID, routeID uint) (*entities.Route, error)
	ListRoutes(ctx context.Context, params dtos.ListRoutesParams) ([]entities.Route, int64, error)
	UpdateRoute(ctx context.Context, route *entities.Route) (*entities.Route, error)
	DeleteRoute(ctx context.Context, businessID, routeID uint) error

	// Route lifecycle
	UpdateRouteStatus(ctx context.Context, routeID uint, status string) error
	UpdateRouteCounters(ctx context.Context, routeID uint) error
	SetRouteActualStart(ctx context.Context, routeID uint) error
	SetRouteActualEnd(ctx context.Context, routeID uint) error

	// Stop CRUD
	AddStop(ctx context.Context, stop *entities.RouteStop) (*entities.RouteStop, error)
	UpdateStop(ctx context.Context, stop *entities.RouteStop) (*entities.RouteStop, error)
	DeleteStop(ctx context.Context, routeID, stopID uint) error
	GetStopByID(ctx context.Context, routeID, stopID uint) (*entities.RouteStop, error)
	UpdateStopStatus(ctx context.Context, stopID uint, status string, failureReason *string, signatureURL, photoURL string) error
	ReorderStops(ctx context.Context, routeID uint, stopIDs []uint) error
	GetStopsByRouteID(ctx context.Context, routeID uint) ([]entities.RouteStop, error)
	SetPendingStopsStatus(ctx context.Context, routeID uint, status string) error

	// Assignable orders (processing status, no driver assigned)
	ListAssignableOrders(ctx context.Context, businessID uint) ([]dtos.AssignableOrder, error)
	ListDriversForBusiness(ctx context.Context, businessID uint) ([]dtos.DriverOption, error)
	ListVehiclesForBusiness(ctx context.Context, businessID uint) ([]dtos.VehicleOption, error)

	// Cross-module queries (replicated locally per isolation rule)
	GetDriverNameByID(ctx context.Context, driverID uint) (string, error)
	UpdateDriverStatus(ctx context.Context, driverID uint, status string) error
	GetVehiclePlateByID(ctx context.Context, vehicleID uint) (string, error)
	UpdateOrderDriverInfo(ctx context.Context, orderID string, driverID *uint, driverName string, isLastMile bool) error
	ClearOrderDriverInfo(ctx context.Context, orderID string) error
}
