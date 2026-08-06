package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/auth/actions/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	GetActionsFn      func(ctx context.Context, page, pageSize int, name string) ([]domain.Action, int64, error)
	GetActionByIDFn   func(ctx context.Context, id uint) (*domain.Action, error)
	GetActionByNameFn func(ctx context.Context, name string) (*domain.Action, error)
	CreateActionFn    func(ctx context.Context, action domain.Action) (uint, error)
	UpdateActionFn    func(ctx context.Context, id uint, action domain.Action) (string, error)
	DeleteActionFn    func(ctx context.Context, id uint) (string, error)

	CreatedActions []domain.Action
	UpdatedActions []domain.Action
}

var _ domain.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) GetActions(ctx context.Context, page, pageSize int, name string) ([]domain.Action, int64, error) {
	if m.GetActionsFn != nil {
		return m.GetActionsFn(ctx, page, pageSize, name)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetActionByID(ctx context.Context, id uint) (*domain.Action, error) {
	if m.GetActionByIDFn != nil {
		return m.GetActionByIDFn(ctx, id)
	}
	return &domain.Action{ID: id}, nil
}

func (m *RepositoryMock) GetActionByName(ctx context.Context, name string) (*domain.Action, error) {
	if m.GetActionByNameFn != nil {
		return m.GetActionByNameFn(ctx, name)
	}
	return nil, nil
}

func (m *RepositoryMock) CreateAction(ctx context.Context, action domain.Action) (uint, error) {
	m.CreatedActions = append(m.CreatedActions, action)
	if m.CreateActionFn != nil {
		return m.CreateActionFn(ctx, action)
	}
	return 1, nil
}

func (m *RepositoryMock) UpdateAction(ctx context.Context, id uint, action domain.Action) (string, error) {
	m.UpdatedActions = append(m.UpdatedActions, action)
	if m.UpdateActionFn != nil {
		return m.UpdateActionFn(ctx, id, action)
	}
	return "ok", nil
}

func (m *RepositoryMock) DeleteAction(ctx context.Context, id uint) (string, error) {
	if m.DeleteActionFn != nil {
		return m.DeleteActionFn(ctx, id)
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
