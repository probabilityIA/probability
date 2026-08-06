package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/ports"
)

type StopStatusCall struct {
	StopID        uint
	Status        string
	FailureReason *string
	SignatureURL  string
	PhotoURL      string
}

type OrderDriverCall struct {
	OrderID    string
	DriverID   *uint
	DriverName string
	IsLastMile bool
}

type DriverStatusCall struct {
	DriverID uint
	Status   string
}

type RepositoryMock struct {
	CreateRouteFn  func(ctx context.Context, route *entities.Route, stops []entities.RouteStop) (*entities.Route, error)
	GetRouteByIDFn func(ctx context.Context, businessID, routeID uint) (*entities.Route, error)
	ListRoutesFn   func(ctx context.Context, params dtos.ListRoutesParams) ([]entities.Route, int64, error)
	UpdateRouteFn  func(ctx context.Context, route *entities.Route) (*entities.Route, error)
	DeleteRouteFn  func(ctx context.Context, businessID, routeID uint) error

	UpdateRouteStatusFn   func(ctx context.Context, routeID uint, status string) error
	UpdateRouteCountersFn func(ctx context.Context, routeID uint) error
	SetRouteActualStartFn func(ctx context.Context, routeID uint) error
	SetRouteActualEndFn   func(ctx context.Context, routeID uint) error

	AddStopFn               func(ctx context.Context, stop *entities.RouteStop) (*entities.RouteStop, error)
	UpdateStopFn            func(ctx context.Context, stop *entities.RouteStop) (*entities.RouteStop, error)
	DeleteStopFn            func(ctx context.Context, routeID, stopID uint) error
	GetStopByIDFn           func(ctx context.Context, routeID, stopID uint) (*entities.RouteStop, error)
	UpdateStopStatusFn      func(ctx context.Context, stopID uint, status string, failureReason *string, signatureURL, photoURL string) error
	ReorderStopsFn          func(ctx context.Context, routeID uint, stopIDs []uint) error
	GetStopsByRouteIDFn     func(ctx context.Context, routeID uint) ([]entities.RouteStop, error)
	SetPendingStopsStatusFn func(ctx context.Context, routeID uint, status string) error

	ListAssignableOrdersFn    func(ctx context.Context, businessID uint) ([]dtos.AssignableOrder, error)
	ListDriversForBusinessFn  func(ctx context.Context, businessID uint) ([]dtos.DriverOption, error)
	ListVehiclesForBusinessFn func(ctx context.Context, businessID uint) ([]dtos.VehicleOption, error)

	GetDriverNameByIDFn     func(ctx context.Context, driverID uint) (string, error)
	UpdateDriverStatusFn    func(ctx context.Context, driverID uint, status string) error
	GetVehiclePlateByIDFn   func(ctx context.Context, vehicleID uint) (string, error)
	UpdateOrderDriverInfoFn func(ctx context.Context, orderID string, driverID *uint, driverName string, isLastMile bool) error
	ClearOrderDriverInfoFn  func(ctx context.Context, orderID string) error

	CreatedRoute      *entities.Route
	CreatedStops      []entities.RouteStop
	AddedStops        []*entities.RouteStop
	UpdatedStops      []*entities.RouteStop
	StopStatusCalls   []StopStatusCall
	OrderDriverCalls  []OrderDriverCall
	ClearedOrders     []string
	DriverStatusCalls []DriverStatusCall
	RouteStatusCalls  []string
	CountersRefreshed int
	ActualStartCalls  int
	ActualEndCalls    int
	PendingSetTo      []string
	ReorderedIDs      []uint
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) CreateRoute(ctx context.Context, route *entities.Route, stops []entities.RouteStop) (*entities.Route, error) {
	m.CreatedRoute = route
	m.CreatedStops = stops
	if m.CreateRouteFn != nil {
		return m.CreateRouteFn(ctx, route, stops)
	}
	route.ID = 1
	return route, nil
}

func (m *RepositoryMock) GetRouteByID(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
	if m.GetRouteByIDFn != nil {
		return m.GetRouteByIDFn(ctx, businessID, routeID)
	}
	return &entities.Route{ID: routeID, BusinessID: businessID, Status: "planned"}, nil
}

func (m *RepositoryMock) ListRoutes(ctx context.Context, params dtos.ListRoutesParams) ([]entities.Route, int64, error) {
	if m.ListRoutesFn != nil {
		return m.ListRoutesFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) UpdateRoute(ctx context.Context, route *entities.Route) (*entities.Route, error) {
	if m.UpdateRouteFn != nil {
		return m.UpdateRouteFn(ctx, route)
	}
	return route, nil
}

func (m *RepositoryMock) DeleteRoute(ctx context.Context, businessID, routeID uint) error {
	if m.DeleteRouteFn != nil {
		return m.DeleteRouteFn(ctx, businessID, routeID)
	}
	return nil
}

func (m *RepositoryMock) UpdateRouteStatus(ctx context.Context, routeID uint, status string) error {
	m.RouteStatusCalls = append(m.RouteStatusCalls, status)
	if m.UpdateRouteStatusFn != nil {
		return m.UpdateRouteStatusFn(ctx, routeID, status)
	}
	return nil
}

func (m *RepositoryMock) UpdateRouteCounters(ctx context.Context, routeID uint) error {
	m.CountersRefreshed++
	if m.UpdateRouteCountersFn != nil {
		return m.UpdateRouteCountersFn(ctx, routeID)
	}
	return nil
}

func (m *RepositoryMock) SetRouteActualStart(ctx context.Context, routeID uint) error {
	m.ActualStartCalls++
	if m.SetRouteActualStartFn != nil {
		return m.SetRouteActualStartFn(ctx, routeID)
	}
	return nil
}

func (m *RepositoryMock) SetRouteActualEnd(ctx context.Context, routeID uint) error {
	m.ActualEndCalls++
	if m.SetRouteActualEndFn != nil {
		return m.SetRouteActualEndFn(ctx, routeID)
	}
	return nil
}

func (m *RepositoryMock) AddStop(ctx context.Context, stop *entities.RouteStop) (*entities.RouteStop, error) {
	m.AddedStops = append(m.AddedStops, stop)
	if m.AddStopFn != nil {
		return m.AddStopFn(ctx, stop)
	}
	stop.ID = 1
	return stop, nil
}

func (m *RepositoryMock) UpdateStop(ctx context.Context, stop *entities.RouteStop) (*entities.RouteStop, error) {
	m.UpdatedStops = append(m.UpdatedStops, stop)
	if m.UpdateStopFn != nil {
		return m.UpdateStopFn(ctx, stop)
	}
	return stop, nil
}

func (m *RepositoryMock) DeleteStop(ctx context.Context, routeID, stopID uint) error {
	if m.DeleteStopFn != nil {
		return m.DeleteStopFn(ctx, routeID, stopID)
	}
	return nil
}

func (m *RepositoryMock) GetStopByID(ctx context.Context, routeID, stopID uint) (*entities.RouteStop, error) {
	if m.GetStopByIDFn != nil {
		return m.GetStopByIDFn(ctx, routeID, stopID)
	}
	return &entities.RouteStop{ID: stopID, RouteID: routeID}, nil
}

func (m *RepositoryMock) UpdateStopStatus(ctx context.Context, stopID uint, status string, failureReason *string, signatureURL, photoURL string) error {
	m.StopStatusCalls = append(m.StopStatusCalls, StopStatusCall{
		StopID: stopID, Status: status, FailureReason: failureReason,
		SignatureURL: signatureURL, PhotoURL: photoURL,
	})
	if m.UpdateStopStatusFn != nil {
		return m.UpdateStopStatusFn(ctx, stopID, status, failureReason, signatureURL, photoURL)
	}
	return nil
}

func (m *RepositoryMock) ReorderStops(ctx context.Context, routeID uint, stopIDs []uint) error {
	m.ReorderedIDs = stopIDs
	if m.ReorderStopsFn != nil {
		return m.ReorderStopsFn(ctx, routeID, stopIDs)
	}
	return nil
}

func (m *RepositoryMock) GetStopsByRouteID(ctx context.Context, routeID uint) ([]entities.RouteStop, error) {
	if m.GetStopsByRouteIDFn != nil {
		return m.GetStopsByRouteIDFn(ctx, routeID)
	}
	return nil, nil
}

func (m *RepositoryMock) SetPendingStopsStatus(ctx context.Context, routeID uint, status string) error {
	m.PendingSetTo = append(m.PendingSetTo, status)
	if m.SetPendingStopsStatusFn != nil {
		return m.SetPendingStopsStatusFn(ctx, routeID, status)
	}
	return nil
}

func (m *RepositoryMock) ListAssignableOrders(ctx context.Context, businessID uint) ([]dtos.AssignableOrder, error) {
	if m.ListAssignableOrdersFn != nil {
		return m.ListAssignableOrdersFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) ListDriversForBusiness(ctx context.Context, businessID uint) ([]dtos.DriverOption, error) {
	if m.ListDriversForBusinessFn != nil {
		return m.ListDriversForBusinessFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) ListVehiclesForBusiness(ctx context.Context, businessID uint) ([]dtos.VehicleOption, error) {
	if m.ListVehiclesForBusinessFn != nil {
		return m.ListVehiclesForBusinessFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) GetDriverNameByID(ctx context.Context, driverID uint) (string, error) {
	if m.GetDriverNameByIDFn != nil {
		return m.GetDriverNameByIDFn(ctx, driverID)
	}
	return "Conductor Demo", nil
}

func (m *RepositoryMock) UpdateDriverStatus(ctx context.Context, driverID uint, status string) error {
	m.DriverStatusCalls = append(m.DriverStatusCalls, DriverStatusCall{DriverID: driverID, Status: status})
	if m.UpdateDriverStatusFn != nil {
		return m.UpdateDriverStatusFn(ctx, driverID, status)
	}
	return nil
}

func (m *RepositoryMock) GetVehiclePlateByID(ctx context.Context, vehicleID uint) (string, error) {
	if m.GetVehiclePlateByIDFn != nil {
		return m.GetVehiclePlateByIDFn(ctx, vehicleID)
	}
	return "ABC123", nil
}

func (m *RepositoryMock) UpdateOrderDriverInfo(ctx context.Context, orderID string, driverID *uint, driverName string, isLastMile bool) error {
	m.OrderDriverCalls = append(m.OrderDriverCalls, OrderDriverCall{
		OrderID: orderID, DriverID: driverID, DriverName: driverName, IsLastMile: isLastMile,
	})
	if m.UpdateOrderDriverInfoFn != nil {
		return m.UpdateOrderDriverInfoFn(ctx, orderID, driverID, driverName, isLastMile)
	}
	return nil
}

func (m *RepositoryMock) ClearOrderDriverInfo(ctx context.Context, orderID string) error {
	m.ClearedOrders = append(m.ClearedOrders, orderID)
	if m.ClearOrderDriverInfoFn != nil {
		return m.ClearOrderDriverInfoFn(ctx, orderID)
	}
	return nil
}
