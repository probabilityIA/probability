package mocks

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/stretchr/testify/mock"
)

// RepositoryMock es un mock de domain.IRepository usando testify/mock
type RepositoryMock struct {
	mock.Mock
}

// ============================================
// INTEGRATIONS
// ============================================

func (m *RepositoryMock) CreateIntegration(ctx context.Context, integration *domain.Integration) error {
	args := m.Called(ctx, integration)
	return args.Error(0)
}

func (m *RepositoryMock) UpdateIntegration(ctx context.Context, id uint, integration *domain.Integration) error {
	args := m.Called(ctx, id, integration)
	return args.Error(0)
}

func (m *RepositoryMock) GetIntegrationByID(ctx context.Context, id uint) (*domain.Integration, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Integration), args.Error(1)
}

func (m *RepositoryMock) DeleteIntegration(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RepositoryMock) ListIntegrations(ctx context.Context, filters domain.IntegrationFilters) ([]*domain.Integration, int64, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Integration), args.Get(1).(int64), args.Error(2)
}

func (m *RepositoryMock) GetIntegrationByIntegrationTypeID(ctx context.Context, integrationTypeID uint, businessID *uint) (*domain.Integration, error) {
	args := m.Called(ctx, integrationTypeID, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Integration), args.Error(1)
}

func (m *RepositoryMock) GetActiveIntegrationByIntegrationTypeID(ctx context.Context, integrationTypeID uint, businessID *uint) (*domain.Integration, error) {
	args := m.Called(ctx, integrationTypeID, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Integration), args.Error(1)
}

func (m *RepositoryMock) ListIntegrationsByBusiness(ctx context.Context, businessID uint) ([]*domain.Integration, error) {
	args := m.Called(ctx, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Integration), args.Error(1)
}

func (m *RepositoryMock) ListIntegrationsByIntegrationTypeID(ctx context.Context, integrationTypeID uint) ([]*domain.Integration, error) {
	args := m.Called(ctx, integrationTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Integration), args.Error(1)
}

func (m *RepositoryMock) SetIntegrationAsDefault(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RepositoryMock) ExistsIntegrationByCode(ctx context.Context, code string, businessID *uint) (bool, error) {
	args := m.Called(ctx, code, businessID)
	return args.Bool(0), args.Error(1)
}

func (m *RepositoryMock) UpdateLastSync(ctx context.Context, id uint, lastSync time.Time) error {
	args := m.Called(ctx, id, lastSync)
	return args.Error(0)
}

// ============================================
// INTEGRATION TYPES
// ============================================

func (m *RepositoryMock) CreateIntegrationType(ctx context.Context, integrationType *domain.IntegrationType) error {
	args := m.Called(ctx, integrationType)
	return args.Error(0)
}

func (m *RepositoryMock) UpdateIntegrationType(ctx context.Context, id uint, integrationType *domain.IntegrationType) error {
	args := m.Called(ctx, id, integrationType)
	return args.Error(0)
}

func (m *RepositoryMock) GetIntegrationTypeByID(ctx context.Context, id uint) (*domain.IntegrationType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IntegrationType), args.Error(1)
}

func (m *RepositoryMock) GetIntegrationTypeByCode(ctx context.Context, code string) (*domain.IntegrationType, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IntegrationType), args.Error(1)
}

func (m *RepositoryMock) GetIntegrationTypeByName(ctx context.Context, name string) (*domain.IntegrationType, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IntegrationType), args.Error(1)
}

func (m *RepositoryMock) DeleteIntegrationType(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RepositoryMock) ListIntegrationTypes(ctx context.Context, categoryID *uint) ([]*domain.IntegrationType, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.IntegrationType), args.Error(1)
}

func (m *RepositoryMock) ListActiveIntegrationTypes(ctx context.Context) ([]*domain.IntegrationType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.IntegrationType), args.Error(1)
}

// ============================================
// INTEGRATION CATEGORIES
// ============================================

func (m *RepositoryMock) GetIntegrationCategoryByID(ctx context.Context, id uint) (*domain.IntegrationCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IntegrationCategory), args.Error(1)
}

func (m *RepositoryMock) ListIntegrationCategories(ctx context.Context) ([]*domain.IntegrationCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.IntegrationCategory), args.Error(1)
}
