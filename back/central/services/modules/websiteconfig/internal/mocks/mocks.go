package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	GetConfigFn             func(ctx context.Context, businessID uint) (*entities.WebsiteConfig, error)
	UpsertConfigFn          func(ctx context.Context, businessID uint, dto *dtos.UpdateConfigDTO) (*entities.WebsiteConfig, error)
	ListProductCategoriesFn func(ctx context.Context, businessID uint) ([]entities.ProductCategory, error)

	ConsultasConfig     []uint
	ConsultasCategorias []uint
	Upserts             []Upsert
}

type Upsert struct {
	BusinessID uint
	DTO        *dtos.UpdateConfigDTO
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) GetConfig(ctx context.Context, businessID uint) (*entities.WebsiteConfig, error) {
	m.ConsultasConfig = append(m.ConsultasConfig, businessID)
	if m.GetConfigFn != nil {
		return m.GetConfigFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) UpsertConfig(ctx context.Context, businessID uint, dto *dtos.UpdateConfigDTO) (*entities.WebsiteConfig, error) {
	m.Upserts = append(m.Upserts, Upsert{BusinessID: businessID, DTO: dto})
	if m.UpsertConfigFn != nil {
		return m.UpsertConfigFn(ctx, businessID, dto)
	}
	return &entities.WebsiteConfig{BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListProductCategories(ctx context.Context, businessID uint) ([]entities.ProductCategory, error) {
	m.ConsultasCategorias = append(m.ConsultasCategorias, businessID)
	if m.ListProductCategoriesFn != nil {
		return m.ListProductCategoriesFn(ctx, businessID)
	}
	return []entities.ProductCategory{}, nil
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
