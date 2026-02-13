package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

// RepositoryMock es un mock del repositorio de órdenes usando testify/mock
type RepositoryMock struct {
	mock.Mock
}

// ============================================
// CRUD OPERATIONS
// ============================================

func (m *RepositoryMock) CreateOrder(ctx context.Context, order *entities.ProbabilityOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *RepositoryMock) GetOrderByID(ctx context.Context, id string) (*entities.ProbabilityOrder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProbabilityOrder), args.Error(1)
}

func (m *RepositoryMock) GetOrderByInternalNumber(ctx context.Context, internalNumber string) (*entities.ProbabilityOrder, error) {
	args := m.Called(ctx, internalNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProbabilityOrder), args.Error(1)
}

func (m *RepositoryMock) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*entities.ProbabilityOrder, error) {
	args := m.Called(ctx, orderNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProbabilityOrder), args.Error(1)
}

func (m *RepositoryMock) ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.ProbabilityOrder, int64, error) {
	args := m.Called(ctx, page, pageSize, filters)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]entities.ProbabilityOrder), args.Get(1).(int64), args.Error(2)
}

func (m *RepositoryMock) UpdateOrder(ctx context.Context, order *entities.ProbabilityOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *RepositoryMock) DeleteOrder(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RepositoryMock) GetOrderRaw(ctx context.Context, id string) (*entities.ProbabilityOrderChannelMetadata, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProbabilityOrderChannelMetadata), args.Error(1)
}

func (m *RepositoryMock) CountOrdersByClientID(ctx context.Context, clientID uint) (int64, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *RepositoryMock) GetLastManualOrderNumber(ctx context.Context, businessID uint) (int, error) {
	args := m.Called(ctx, businessID)
	return args.Int(0), args.Error(1)
}

func (m *RepositoryMock) GetFirstIntegrationIDByBusinessID(ctx context.Context, businessID uint) (uint, error) {
	args := m.Called(ctx, businessID)
	return args.Get(0).(uint), args.Error(1)
}

// ============================================
// VALIDATION
// ============================================

func (m *RepositoryMock) OrderExists(ctx context.Context, externalID string, integrationID uint) (bool, error) {
	args := m.Called(ctx, externalID, integrationID)
	return args.Bool(0), args.Error(1)
}

func (m *RepositoryMock) GetOrderByExternalID(ctx context.Context, externalID string, integrationID uint) (*entities.ProbabilityOrder, error) {
	args := m.Called(ctx, externalID, integrationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProbabilityOrder), args.Error(1)
}

// ============================================
// TABLAS RELACIONADAS
// ============================================

func (m *RepositoryMock) CreateOrderItems(ctx context.Context, items []*entities.ProbabilityOrderItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *RepositoryMock) CreateAddresses(ctx context.Context, addresses []*entities.ProbabilityAddress) error {
	args := m.Called(ctx, addresses)
	return args.Error(0)
}

func (m *RepositoryMock) CreatePayments(ctx context.Context, payments []*entities.ProbabilityPayment) error {
	args := m.Called(ctx, payments)
	return args.Error(0)
}

func (m *RepositoryMock) CreateShipments(ctx context.Context, shipments []*entities.ProbabilityShipment) error {
	args := m.Called(ctx, shipments)
	return args.Error(0)
}

func (m *RepositoryMock) CreateChannelMetadata(ctx context.Context, metadata *entities.ProbabilityOrderChannelMetadata) error {
	args := m.Called(ctx, metadata)
	return args.Error(0)
}

// ============================================
// CATÁLOGO (VALIDACIÓN)
// ============================================

func (m *RepositoryMock) GetProductBySKU(ctx context.Context, businessID uint, sku string) (*entities.Product, error) {
	args := m.Called(ctx, businessID, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *RepositoryMock) CreateProduct(ctx context.Context, product *entities.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *RepositoryMock) GetClientByEmail(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
	args := m.Called(ctx, businessID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *RepositoryMock) GetClientByDNI(ctx context.Context, businessID uint, dni string) (*entities.Client, error) {
	args := m.Called(ctx, businessID, dni)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *RepositoryMock) CreateClient(ctx context.Context, client *entities.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *RepositoryMock) CreateOrderError(ctx context.Context, orderError *entities.OrderError) error {
	args := m.Called(ctx, orderError)
	return args.Error(0)
}

// ============================================
// CONSULTAS A TABLAS DE ESTADOS
// ============================================

func (m *RepositoryMock) GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
	args := m.Called(ctx, integrationTypeID, originalStatus)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*uint), args.Error(1)
}

func (m *RepositoryMock) GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*uint), args.Error(1)
}

func (m *RepositoryMock) GetFulfillmentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*uint), args.Error(1)
}
