package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
	"github.com/stretchr/testify/mock"
)

// UseCaseMock es un mock de app.Iapp usando testify/mock
type UseCaseMock struct {
	mock.Mock
}

func (m *UseCaseMock) GetPermissions(ctx context.Context, businessTypeID *uint, name *string, scopeID *uint, resource *string) ([]domain.PermissionDTO, error) {
	args := m.Called(ctx, businessTypeID, name, scopeID, resource)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PermissionDTO), args.Error(1)
}

func (m *UseCaseMock) GetPermissionByID(ctx context.Context, id uint) (*domain.PermissionDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PermissionDTO), args.Error(1)
}

func (m *UseCaseMock) CreatePermission(ctx context.Context, permission domain.CreatePermissionDTO) (string, error) {
	args := m.Called(ctx, permission)
	return args.String(0), args.Error(1)
}

func (m *UseCaseMock) BulkCreatePermissions(ctx context.Context, permissions []domain.CreatePermissionDTO) ([]domain.BulkCreateResult, error) {
	args := m.Called(ctx, permissions)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.BulkCreateResult), args.Error(1)
}

func (m *UseCaseMock) UpdatePermission(ctx context.Context, id uint, permission domain.UpdatePermissionDTO) (string, error) {
	args := m.Called(ctx, id, permission)
	return args.String(0), args.Error(1)
}

func (m *UseCaseMock) DeletePermission(ctx context.Context, id uint) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *UseCaseMock) GetPermissionsByResource(ctx context.Context, resource string) ([]domain.PermissionDTO, error) {
	args := m.Called(ctx, resource)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PermissionDTO), args.Error(1)
}

func (m *UseCaseMock) GetPermissionsByScopeID(ctx context.Context, scopeID uint) ([]domain.PermissionDTO, error) {
	args := m.Called(ctx, scopeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PermissionDTO), args.Error(1)
}
