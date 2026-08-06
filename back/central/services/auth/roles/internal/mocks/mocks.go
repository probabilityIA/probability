package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RoleRepositoryMock struct {
	GetRoleByIDFn              func(ctx context.Context, id uint) (*domain.Role, error)
	AssignPermissionsToRoleFn  func(ctx context.Context, roleID uint, permissionIDs []uint) error
	CreateRoleFn               func(ctx context.Context, role domain.CreateRoleDTO) (*domain.Role, error)
	UpdateRoleFn               func(ctx context.Context, id uint, role domain.UpdateRoleDTO) (*domain.Role, error)
	RoleExistsByNameFn         func(ctx context.Context, name string, excludeID *uint) (bool, error)
	GetRolePermissionsFn       func(ctx context.Context, roleID uint) ([]domain.Permission, error)
	GetRolesByLevelFn          func(ctx context.Context, level int) ([]domain.Role, error)
	GetRolesByScopeIDFn        func(ctx context.Context, scopeID uint) ([]domain.Role, error)
	GetRolesFn                 func(ctx context.Context) ([]domain.Role, error)
	GetSystemRolesFn           func(ctx context.Context) ([]domain.Role, error)
	RemovePermissionFromRoleFn func(ctx context.Context, roleID uint, permissionID uint) error
}

var _ domain.IRoleRepository = (*RoleRepositoryMock)(nil)

func (m *RoleRepositoryMock) GetRoleByID(ctx context.Context, id uint) (*domain.Role, error) {
	if m.GetRoleByIDFn != nil {
		return m.GetRoleByIDFn(ctx, id)
	}
	return &domain.Role{ID: id}, nil
}

func (m *RoleRepositoryMock) AssignPermissionsToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	if m.AssignPermissionsToRoleFn != nil {
		return m.AssignPermissionsToRoleFn(ctx, roleID, permissionIDs)
	}
	return nil
}

func (m *RoleRepositoryMock) CreateRole(ctx context.Context, role domain.CreateRoleDTO) (*domain.Role, error) {
	if m.CreateRoleFn != nil {
		return m.CreateRoleFn(ctx, role)
	}
	return &domain.Role{ID: 1, Name: role.Name}, nil
}

func (m *RoleRepositoryMock) UpdateRole(ctx context.Context, id uint, role domain.UpdateRoleDTO) (*domain.Role, error) {
	if m.UpdateRoleFn != nil {
		return m.UpdateRoleFn(ctx, id, role)
	}
	return &domain.Role{ID: id}, nil
}

func (m *RoleRepositoryMock) RoleExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error) {
	if m.RoleExistsByNameFn != nil {
		return m.RoleExistsByNameFn(ctx, name, excludeID)
	}
	return false, nil
}

func (m *RoleRepositoryMock) GetRolePermissions(ctx context.Context, roleID uint) ([]domain.Permission, error) {
	if m.GetRolePermissionsFn != nil {
		return m.GetRolePermissionsFn(ctx, roleID)
	}
	return nil, nil
}

func (m *RoleRepositoryMock) GetRolesByLevel(ctx context.Context, level int) ([]domain.Role, error) {
	if m.GetRolesByLevelFn != nil {
		return m.GetRolesByLevelFn(ctx, level)
	}
	return nil, nil
}

func (m *RoleRepositoryMock) GetRolesByScopeID(ctx context.Context, scopeID uint) ([]domain.Role, error) {
	if m.GetRolesByScopeIDFn != nil {
		return m.GetRolesByScopeIDFn(ctx, scopeID)
	}
	return nil, nil
}

func (m *RoleRepositoryMock) GetRoles(ctx context.Context) ([]domain.Role, error) {
	if m.GetRolesFn != nil {
		return m.GetRolesFn(ctx)
	}
	return nil, nil
}

func (m *RoleRepositoryMock) GetSystemRoles(ctx context.Context) ([]domain.Role, error) {
	if m.GetSystemRolesFn != nil {
		return m.GetSystemRolesFn(ctx)
	}
	return nil, nil
}

func (m *RoleRepositoryMock) RemovePermissionFromRole(ctx context.Context, roleID uint, permissionID uint) error {
	if m.RemovePermissionFromRoleFn != nil {
		return m.RemovePermissionFromRoleFn(ctx, roleID, permissionID)
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
