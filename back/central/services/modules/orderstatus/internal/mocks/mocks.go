package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	CreateFn       func(ctx context.Context, mapping *entities.OrderStatusMapping) error
	GetByIDFn      func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error)
	ListFn         func(ctx context.Context, filters map[string]interface{}) ([]entities.OrderStatusMapping, int64, error)
	UpdateFn       func(ctx context.Context, mapping *entities.OrderStatusMapping) error
	DeleteFn       func(ctx context.Context, id uint) error
	ToggleActiveFn func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error)
	ExistsFn       func(ctx context.Context, integrationTypeID uint, originalStatus string) (bool, error)

	GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error)
	ListOrderStatusesFn                                  func(ctx context.Context, isActive *bool) ([]entities.OrderStatusInfo, error)
	ListFulfillmentStatusesFn                            func(ctx context.Context, isActive *bool) ([]entities.FulfillmentStatusInfo, error)

	CreateOrderStatusFn  func(ctx context.Context, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error)
	GetOrderStatusByIDFn func(ctx context.Context, id uint) (*entities.OrderStatusInfo, error)
	UpdateOrderStatusFn  func(ctx context.Context, id uint, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error)
	DeleteOrderStatusFn  func(ctx context.Context, id uint) error

	ListEcommerceIntegrationTypesFn func(ctx context.Context, businessID uint) ([]entities.IntegrationTypeInfo, error)
	ListChannelStatusesFn           func(ctx context.Context, integrationTypeID uint, isActive *bool) ([]entities.ChannelStatusInfo, error)
	CreateChannelStatusFn           func(ctx context.Context, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error)
	GetChannelStatusByIDFn          func(ctx context.Context, id uint) (*entities.ChannelStatusInfo, error)
	UpdateChannelStatusFn           func(ctx context.Context, id uint, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error)
	DeleteChannelStatusFn           func(ctx context.Context, id uint) error
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) Create(ctx context.Context, mapping *entities.OrderStatusMapping) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, mapping)
	}
	return nil
}

func (m *RepositoryMock) GetByID(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return &entities.OrderStatusMapping{ID: id}, nil
}

func (m *RepositoryMock) List(ctx context.Context, filters map[string]interface{}) ([]entities.OrderStatusMapping, int64, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filters)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) Update(ctx context.Context, mapping *entities.OrderStatusMapping) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, mapping)
	}
	return nil
}

func (m *RepositoryMock) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) ToggleActive(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
	if m.ToggleActiveFn != nil {
		return m.ToggleActiveFn(ctx, id)
	}
	return &entities.OrderStatusMapping{ID: id}, nil
}

func (m *RepositoryMock) Exists(ctx context.Context, integrationTypeID uint, originalStatus string) (bool, error) {
	if m.ExistsFn != nil {
		return m.ExistsFn(ctx, integrationTypeID, originalStatus)
	}
	return false, nil
}

func (m *RepositoryMock) GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
	if m.GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn != nil {
		return m.GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn(ctx, integrationTypeID, originalStatus)
	}
	return nil, nil
}

func (m *RepositoryMock) ListOrderStatuses(ctx context.Context, isActive *bool) ([]entities.OrderStatusInfo, error) {
	if m.ListOrderStatusesFn != nil {
		return m.ListOrderStatusesFn(ctx, isActive)
	}
	return nil, nil
}

func (m *RepositoryMock) ListFulfillmentStatuses(ctx context.Context, isActive *bool) ([]entities.FulfillmentStatusInfo, error) {
	if m.ListFulfillmentStatusesFn != nil {
		return m.ListFulfillmentStatusesFn(ctx, isActive)
	}
	return nil, nil
}

func (m *RepositoryMock) CreateOrderStatus(ctx context.Context, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error) {
	if m.CreateOrderStatusFn != nil {
		return m.CreateOrderStatusFn(ctx, status)
	}
	return status, nil
}

func (m *RepositoryMock) GetOrderStatusByID(ctx context.Context, id uint) (*entities.OrderStatusInfo, error) {
	if m.GetOrderStatusByIDFn != nil {
		return m.GetOrderStatusByIDFn(ctx, id)
	}
	return &entities.OrderStatusInfo{ID: id}, nil
}

func (m *RepositoryMock) UpdateOrderStatus(ctx context.Context, id uint, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error) {
	if m.UpdateOrderStatusFn != nil {
		return m.UpdateOrderStatusFn(ctx, id, status)
	}
	return status, nil
}

func (m *RepositoryMock) DeleteOrderStatus(ctx context.Context, id uint) error {
	if m.DeleteOrderStatusFn != nil {
		return m.DeleteOrderStatusFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) ListEcommerceIntegrationTypes(ctx context.Context, businessID uint) ([]entities.IntegrationTypeInfo, error) {
	if m.ListEcommerceIntegrationTypesFn != nil {
		return m.ListEcommerceIntegrationTypesFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) ListChannelStatuses(ctx context.Context, integrationTypeID uint, isActive *bool) ([]entities.ChannelStatusInfo, error) {
	if m.ListChannelStatusesFn != nil {
		return m.ListChannelStatusesFn(ctx, integrationTypeID, isActive)
	}
	return nil, nil
}

func (m *RepositoryMock) CreateChannelStatus(ctx context.Context, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error) {
	if m.CreateChannelStatusFn != nil {
		return m.CreateChannelStatusFn(ctx, status)
	}
	return status, nil
}

func (m *RepositoryMock) GetChannelStatusByID(ctx context.Context, id uint) (*entities.ChannelStatusInfo, error) {
	if m.GetChannelStatusByIDFn != nil {
		return m.GetChannelStatusByIDFn(ctx, id)
	}
	return &entities.ChannelStatusInfo{ID: id}, nil
}

func (m *RepositoryMock) UpdateChannelStatus(ctx context.Context, id uint, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error) {
	if m.UpdateChannelStatusFn != nil {
		return m.UpdateChannelStatusFn(ctx, id, status)
	}
	return status, nil
}

func (m *RepositoryMock) DeleteChannelStatus(ctx context.Context, id uint) error {
	if m.DeleteChannelStatusFn != nil {
		return m.DeleteChannelStatusFn(ctx, id)
	}
	return nil
}

type SilentLogger struct{}

func NewSilentLogger() log.ILogger {
	return &SilentLogger{}
}

func (l *SilentLogger) nop() zerolog.Logger {
	return zerolog.Nop()
}

func (l *SilentLogger) Info(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Info()
}

func (l *SilentLogger) Error(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Error()
}

func (l *SilentLogger) Warn(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Warn()
}

func (l *SilentLogger) Debug(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Debug()
}

func (l *SilentLogger) Fatal(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Fatal()
}

func (l *SilentLogger) Panic(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Panic()
}

func (l *SilentLogger) With() zerolog.Context {
	n := l.nop()
	return n.With()
}

func (l *SilentLogger) WithService(service string) log.ILogger {
	return l
}

func (l *SilentLogger) WithModule(module string) log.ILogger {
	return l
}

func (l *SilentLogger) WithBusinessID(businessID uint) log.ILogger {
	return l
}
