package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
	"github.com/stretchr/testify/mock"
)

// RepositoryMock es un mock del repositorio de permisos usando testify/mock
type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) GetPermissions(ctx context.Context, businessTypeID *uint, name *string, scopeID *uint, resource *string) ([]domain.Permission, error) {
	args := m.Called(ctx, businessTypeID, name, scopeID, resource)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Permission), args.Error(1)
}

func (m *RepositoryMock) GetPermissionByID(ctx context.Context, id uint) (*domain.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Permission), args.Error(1)
}

func (m *RepositoryMock) GetPermissionsByScopeID(ctx context.Context, scopeID uint) ([]domain.Permission, error) {
	args := m.Called(ctx, scopeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Permission), args.Error(1)
}

func (m *RepositoryMock) GetPermissionsByResource(ctx context.Context, resource string) ([]domain.Permission, error) {
	args := m.Called(ctx, resource)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Permission), args.Error(1)
}

func (m *RepositoryMock) PermissionExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *RepositoryMock) CreatePermission(ctx context.Context, permission domain.Permission) (string, error) {
	args := m.Called(ctx, permission)
	return args.String(0), args.Error(1)
}

func (m *RepositoryMock) UpdatePermission(ctx context.Context, id uint, permission domain.Permission) (string, error) {
	args := m.Called(ctx, id, permission)
	return args.String(0), args.Error(1)
}

func (m *RepositoryMock) DeletePermission(ctx context.Context, id uint) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}
