package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/auth/resources/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	GetResourcesFn      func(ctx context.Context, filters domain.ResourceFilters) ([]domain.Resource, int64, error)
	GetResourceByIDFn   func(ctx context.Context, id uint) (*domain.Resource, error)
	GetResourceByNameFn func(ctx context.Context, name string) (*domain.Resource, error)
	CreateResourceFn    func(ctx context.Context, resource domain.Resource) (uint, error)
	UpdateResourceFn    func(ctx context.Context, id uint, resource domain.Resource) (string, error)
	DeleteResourceFn    func(ctx context.Context, id uint) (string, error)

	CreatedResources []domain.Resource
	UpdatedResources []domain.Resource
}

var _ domain.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) GetResources(ctx context.Context, filters domain.ResourceFilters) ([]domain.Resource, int64, error) {
	if m.GetResourcesFn != nil {
		return m.GetResourcesFn(ctx, filters)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetResourceByID(ctx context.Context, id uint) (*domain.Resource, error) {
	if m.GetResourceByIDFn != nil {
		return m.GetResourceByIDFn(ctx, id)
	}
	return &domain.Resource{ID: id}, nil
}

func (m *RepositoryMock) GetResourceByName(ctx context.Context, name string) (*domain.Resource, error) {
	if m.GetResourceByNameFn != nil {
		return m.GetResourceByNameFn(ctx, name)
	}
	return nil, nil
}

func (m *RepositoryMock) CreateResource(ctx context.Context, resource domain.Resource) (uint, error) {
	m.CreatedResources = append(m.CreatedResources, resource)
	if m.CreateResourceFn != nil {
		return m.CreateResourceFn(ctx, resource)
	}
	return 1, nil
}

func (m *RepositoryMock) UpdateResource(ctx context.Context, id uint, resource domain.Resource) (string, error) {
	m.UpdatedResources = append(m.UpdatedResources, resource)
	if m.UpdateResourceFn != nil {
		return m.UpdateResourceFn(ctx, id, resource)
	}
	return "ok", nil
}

func (m *RepositoryMock) DeleteResource(ctx context.Context, id uint) (string, error) {
	if m.DeleteResourceFn != nil {
		return m.DeleteResourceFn(ctx, id)
	}
	return "eliminado", nil
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
