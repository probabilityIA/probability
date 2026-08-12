package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	SaveFn                         func(ctx context.Context, run *domain.SyncRun) error
	ListLastByBusinessFn           func(ctx context.Context, businessID uint) ([]domain.SyncRun, error)
	ListDetailFn                   func(ctx context.Context, query domain.DetailQuery) (*domain.DetailPage, error)
	IntegrationBelongsToBusinessFn func(ctx context.Context, integrationID, businessID uint) (bool, error)
	ChannelSummariesFn             func(ctx context.Context, businessID uint) ([]domain.ChannelSummary, error)
	CrossChannelFn                 func(ctx context.Context, businessID uint) (domain.CrossChannelCounts, error)

	Guardados       []domain.SyncRun
	ConsultasDueno  []ConsultaDueno
	ConsultasLast   []uint
	ConsultasDetail []domain.DetailQuery
}

type ConsultaDueno struct {
	IntegrationID uint
	BusinessID    uint
}

var _ domain.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) Save(ctx context.Context, run *domain.SyncRun) error {
	m.Guardados = append(m.Guardados, *run)
	if m.SaveFn != nil {
		return m.SaveFn(ctx, run)
	}
	return nil
}

func (m *RepositoryMock) ListLastByBusiness(ctx context.Context, businessID uint) ([]domain.SyncRun, error) {
	m.ConsultasLast = append(m.ConsultasLast, businessID)
	if m.ListLastByBusinessFn != nil {
		return m.ListLastByBusinessFn(ctx, businessID)
	}
	return []domain.SyncRun{}, nil
}

func (m *RepositoryMock) ListDetail(ctx context.Context, query domain.DetailQuery) (*domain.DetailPage, error) {
	m.ConsultasDetail = append(m.ConsultasDetail, query)
	if m.ListDetailFn != nil {
		return m.ListDetailFn(ctx, query)
	}
	return &domain.DetailPage{Page: query.Page, PageSize: query.PageSize}, nil
}

func (m *RepositoryMock) IntegrationBelongsToBusiness(ctx context.Context, integrationID, businessID uint) (bool, error) {
	m.ConsultasDueno = append(m.ConsultasDueno, ConsultaDueno{IntegrationID: integrationID, BusinessID: businessID})
	if m.IntegrationBelongsToBusinessFn != nil {
		return m.IntegrationBelongsToBusinessFn(ctx, integrationID, businessID)
	}
	return true, nil
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

func (m *RepositoryMock) ChannelSummaries(ctx context.Context, businessID uint) ([]domain.ChannelSummary, error) {
	if m.ChannelSummariesFn != nil {
		return m.ChannelSummariesFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) CrossChannel(ctx context.Context, businessID uint) (domain.CrossChannelCounts, error) {
	if m.CrossChannelFn != nil {
		return m.CrossChannelFn(ctx, businessID)
	}
	return domain.CrossChannelCounts{}, nil
}
