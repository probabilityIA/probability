package mocks

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type AlertPublisherMock struct {
	PublishFn func(ctx context.Context, event entities.AlertEvent) error

	mu        sync.Mutex
	Published []entities.AlertEvent
}

var _ ports.IAlertPublisher = (*AlertPublisherMock)(nil)

func (m *AlertPublisherMock) Publish(ctx context.Context, event entities.AlertEvent) error {
	m.mu.Lock()
	m.Published = append(m.Published, event)
	m.mu.Unlock()
	if m.PublishFn != nil {
		return m.PublishFn(ctx, event)
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
