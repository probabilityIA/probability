package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

func (m *RepositoryMock) GetIntegrationStats(ctx context.Context, businessID uint) ([]domain.IntegrationStats, error) {
	args := m.Called(ctx, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.IntegrationStats), args.Error(1)
}

func (m *CacheMock) SetIntegrationStats(ctx context.Context, businessID uint, stats []domain.IntegrationStats) error {
	args := m.Called(ctx, businessID, stats)
	return args.Error(0)
}

func (m *CacheMock) GetIntegrationStats(ctx context.Context, businessID uint) ([]domain.IntegrationStats, error) {
	args := m.Called(ctx, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.IntegrationStats), args.Error(1)
}
