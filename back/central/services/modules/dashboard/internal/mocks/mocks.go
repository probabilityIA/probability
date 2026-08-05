package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	GetTotalOrdersFn               func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) (int64, error)
	GetOrdersTodayFn               func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) (int64, error)
	GetOrdersByIntegrationTypeFn   func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrderCountByIntegrationType, error)
	GetTopCustomersFn              func(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.TopCustomer, error)
	GetOrdersByLocationFn          func(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.OrderCountByLocation, error)
	GetTopDriversFn                func(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.TopDriver, error)
	GetDriversByLocationFn         func(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.DriverByLocation, error)
	GetTopProductsFn               func(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.TopProduct, error)
	GetProductsByCategoryFn        func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ProductByCategory, error)
	GetProductsByBrandFn           func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ProductByBrand, error)
	GetShipmentsByStatusFn         func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByStatus, error)
	GetShipmentsByStatusFilteredFn func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByStatus, error)
	GetShipmentsByCarrierFn        func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByCarrier, error)
	GetShipmentsByCarrierTodayFn   func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByCarrier, error)
	GetShipmentsByWarehouseFn      func(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByWarehouse, error)
	GetShipmentsByDayOfWeekFn      func(ctx context.Context, businessID *uint, integrationID *uint, weekStartDate *time.Time) ([]domain.ShipmentsByDayOfWeek, error)
	GetOrdersByDepartmentFn        func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByDepartment, error)
	GetOrdersByBusinessFn          func(ctx context.Context, limit int, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByBusiness, error)
	GetOrdersByMonthFn             func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByMonth, error)
	GetOrdersByWeekFn              func(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByWeek, error)
	GetTopSellingDaysFn            func(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]domain.TopSellingDay, error)

	mu     sync.Mutex
	called []string
}

var _ domain.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) record(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.called = append(m.called, name)
}

func (m *RepositoryMock) Called() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]string, len(m.called))
	copy(out, m.called)
	return out
}

func (m *RepositoryMock) WasCalled(name string) bool {
	for _, c := range m.Called() {
		if c == name {
			return true
		}
	}
	return false
}

func (m *RepositoryMock) GetTotalOrders(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) (int64, error) {
	m.record("GetTotalOrders")
	if m.GetTotalOrdersFn != nil {
		return m.GetTotalOrdersFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return 0, nil
}

func (m *RepositoryMock) GetOrdersToday(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) (int64, error) {
	m.record("GetOrdersToday")
	if m.GetOrdersTodayFn != nil {
		return m.GetOrdersTodayFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return 0, nil
}

func (m *RepositoryMock) GetOrdersByIntegrationType(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrderCountByIntegrationType, error) {
	m.record("GetOrdersByIntegrationType")
	if m.GetOrdersByIntegrationTypeFn != nil {
		return m.GetOrdersByIntegrationTypeFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetTopCustomers(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.TopCustomer, error) {
	m.record("GetTopCustomers")
	if m.GetTopCustomersFn != nil {
		return m.GetTopCustomersFn(ctx, businessID, integrationID, limit, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetOrdersByLocation(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.OrderCountByLocation, error) {
	m.record("GetOrdersByLocation")
	if m.GetOrdersByLocationFn != nil {
		return m.GetOrdersByLocationFn(ctx, businessID, integrationID, limit, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetTopDrivers(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.TopDriver, error) {
	m.record("GetTopDrivers")
	if m.GetTopDriversFn != nil {
		return m.GetTopDriversFn(ctx, businessID, integrationID, limit, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetDriversByLocation(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.DriverByLocation, error) {
	m.record("GetDriversByLocation")
	if m.GetDriversByLocationFn != nil {
		return m.GetDriversByLocationFn(ctx, businessID, integrationID, limit, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetTopProducts(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.TopProduct, error) {
	m.record("GetTopProducts")
	if m.GetTopProductsFn != nil {
		return m.GetTopProductsFn(ctx, businessID, integrationID, limit, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetProductsByCategory(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ProductByCategory, error) {
	m.record("GetProductsByCategory")
	if m.GetProductsByCategoryFn != nil {
		return m.GetProductsByCategoryFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetProductsByBrand(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ProductByBrand, error) {
	m.record("GetProductsByBrand")
	if m.GetProductsByBrandFn != nil {
		return m.GetProductsByBrandFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetShipmentsByStatus(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByStatus, error) {
	m.record("GetShipmentsByStatus")
	if m.GetShipmentsByStatusFn != nil {
		return m.GetShipmentsByStatusFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetShipmentsByStatusFiltered(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByStatus, error) {
	m.record("GetShipmentsByStatusFiltered")
	if m.GetShipmentsByStatusFilteredFn != nil {
		return m.GetShipmentsByStatusFilteredFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetShipmentsByCarrier(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByCarrier, error) {
	m.record("GetShipmentsByCarrier")
	if m.GetShipmentsByCarrierFn != nil {
		return m.GetShipmentsByCarrierFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetShipmentsByCarrierToday(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByCarrier, error) {
	m.record("GetShipmentsByCarrierToday")
	if m.GetShipmentsByCarrierTodayFn != nil {
		return m.GetShipmentsByCarrierTodayFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetShipmentsByWarehouse(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByWarehouse, error) {
	m.record("GetShipmentsByWarehouse")
	if m.GetShipmentsByWarehouseFn != nil {
		return m.GetShipmentsByWarehouseFn(ctx, businessID, integrationID, limit, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetShipmentsByDayOfWeek(ctx context.Context, businessID *uint, integrationID *uint, weekStartDate *time.Time) ([]domain.ShipmentsByDayOfWeek, error) {
	m.record("GetShipmentsByDayOfWeek")
	if m.GetShipmentsByDayOfWeekFn != nil {
		return m.GetShipmentsByDayOfWeekFn(ctx, businessID, integrationID, weekStartDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetOrdersByDepartment(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByDepartment, error) {
	m.record("GetOrdersByDepartment")
	if m.GetOrdersByDepartmentFn != nil {
		return m.GetOrdersByDepartmentFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetOrdersByBusiness(ctx context.Context, limit int, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByBusiness, error) {
	m.record("GetOrdersByBusiness")
	if m.GetOrdersByBusinessFn != nil {
		return m.GetOrdersByBusinessFn(ctx, limit, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetOrdersByMonth(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByMonth, error) {
	m.record("GetOrdersByMonth")
	if m.GetOrdersByMonthFn != nil {
		return m.GetOrdersByMonthFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetOrdersByWeek(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByWeek, error) {
	m.record("GetOrdersByWeek")
	if m.GetOrdersByWeekFn != nil {
		return m.GetOrdersByWeekFn(ctx, businessID, integrationID, startDate, endDate)
	}
	return nil, nil
}

func (m *RepositoryMock) GetTopSellingDays(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]domain.TopSellingDay, error) {
	m.record("GetTopSellingDays")
	if m.GetTopSellingDaysFn != nil {
		return m.GetTopSellingDaysFn(ctx, businessID, integrationID, limit)
	}
	return nil, nil
}

type SilentLogger struct{}

func NewSilentLogger() log.ILogger { return &SilentLogger{} }

func (l *SilentLogger) nop() zerolog.Logger { return zerolog.Nop() }

func (l *SilentLogger) Info(ctx ...context.Context) *zerolog.Event  { n := l.nop(); return n.Info() }
func (l *SilentLogger) Error(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Error() }
func (l *SilentLogger) Warn(ctx ...context.Context) *zerolog.Event  { n := l.nop(); return n.Warn() }
func (l *SilentLogger) Debug(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Debug() }
func (l *SilentLogger) Fatal(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Fatal() }
func (l *SilentLogger) Panic(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Panic() }
func (l *SilentLogger) With() zerolog.Context                       { n := l.nop(); return n.With() }
func (l *SilentLogger) WithService(service string) log.ILogger      { return l }
func (l *SilentLogger) WithModule(module string) log.ILogger        { return l }
func (l *SilentLogger) WithBusinessID(businessID uint) log.ILogger  { return l }
