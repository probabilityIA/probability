package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/stretchr/testify/mock"
)

// CacheMock es un mock de domain.IIntegrationCache usando testify/mock
type CacheMock struct {
	mock.Mock
}

func (m *CacheMock) SetIntegration(ctx context.Context, integration *domain.CachedIntegration) error {
	args := m.Called(ctx, integration)
	return args.Error(0)
}

func (m *CacheMock) GetIntegration(ctx context.Context, integrationID uint) (*domain.CachedIntegration, error) {
	args := m.Called(ctx, integrationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CachedIntegration), args.Error(1)
}

func (m *CacheMock) SetCredentials(ctx context.Context, creds *domain.CachedCredentials) error {
	args := m.Called(ctx, creds)
	return args.Error(0)
}

func (m *CacheMock) GetCredentials(ctx context.Context, integrationID uint) (*domain.CachedCredentials, error) {
	args := m.Called(ctx, integrationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CachedCredentials), args.Error(1)
}

func (m *CacheMock) GetCredentialField(ctx context.Context, integrationID uint, field string) (string, error) {
	args := m.Called(ctx, integrationID, field)
	return args.String(0), args.Error(1)
}

func (m *CacheMock) InvalidateIntegration(ctx context.Context, integrationID uint) error {
	args := m.Called(ctx, integrationID)
	return args.Error(0)
}

func (m *CacheMock) GetByCode(ctx context.Context, code string) (*domain.CachedIntegration, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CachedIntegration), args.Error(1)
}

func (m *CacheMock) GetByBusinessAndType(ctx context.Context, businessID, integrationTypeID uint) (*domain.CachedIntegration, error) {
	args := m.Called(ctx, businessID, integrationTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CachedIntegration), args.Error(1)
}

func (m *CacheMock) SetPlatformCredentials(ctx context.Context, integrationTypeID uint, creds map[string]interface{}) error {
	args := m.Called(ctx, integrationTypeID, creds)
	return args.Error(0)
}

func (m *CacheMock) GetPlatformCredentials(ctx context.Context, integrationTypeID uint) (map[string]interface{}, error) {
	args := m.Called(ctx, integrationTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
