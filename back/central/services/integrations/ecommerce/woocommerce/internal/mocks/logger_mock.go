package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/shared/log"
)

// LoggerMock mock del logger para tests unitarios del módulo WooCommerce.
// Descarta todos los eventos por defecto (null logger).
type LoggerMock struct {
	InfoFn  func(ctx ...context.Context) *zerolog.Event
	ErrorFn func(ctx ...context.Context) *zerolog.Event
	WarnFn  func(ctx ...context.Context) *zerolog.Event
	DebugFn func(ctx ...context.Context) *zerolog.Event
	FatalFn func(ctx ...context.Context) *zerolog.Event
	PanicFn func(ctx ...context.Context) *zerolog.Event
}

// NewLoggerMock crea un LoggerMock que descarta todos los eventos (null logger).
func NewLoggerMock() log.ILogger {
	noop := zerolog.Nop()
	return &LoggerMock{
		InfoFn: func(ctx ...context.Context) *zerolog.Event {
			return noop.Info()
		},
		ErrorFn: func(ctx ...context.Context) *zerolog.Event {
			return noop.Error()
		},
		WarnFn: func(ctx ...context.Context) *zerolog.Event {
			return noop.Warn()
		},
		DebugFn: func(ctx ...context.Context) *zerolog.Event {
			return noop.Debug()
		},
		FatalFn: func(ctx ...context.Context) *zerolog.Event {
			return noop.Fatal()
		},
		PanicFn: func(ctx ...context.Context) *zerolog.Event {
			return noop.Panic()
		},
	}
}

func (m *LoggerMock) Info(ctx ...context.Context) *zerolog.Event {
	if m.InfoFn != nil {
		return m.InfoFn(ctx...)
	}
	noop := zerolog.Nop()
	return noop.Info()
}

func (m *LoggerMock) Error(ctx ...context.Context) *zerolog.Event {
	if m.ErrorFn != nil {
		return m.ErrorFn(ctx...)
	}
	noop := zerolog.Nop()
	return noop.Error()
}

func (m *LoggerMock) Warn(ctx ...context.Context) *zerolog.Event {
	if m.WarnFn != nil {
		return m.WarnFn(ctx...)
	}
	noop := zerolog.Nop()
	return noop.Warn()
}

func (m *LoggerMock) Debug(ctx ...context.Context) *zerolog.Event {
	if m.DebugFn != nil {
		return m.DebugFn(ctx...)
	}
	noop := zerolog.Nop()
	return noop.Debug()
}

func (m *LoggerMock) Fatal(ctx ...context.Context) *zerolog.Event {
	if m.FatalFn != nil {
		return m.FatalFn(ctx...)
	}
	noop := zerolog.Nop()
	return noop.Fatal()
}

func (m *LoggerMock) Panic(ctx ...context.Context) *zerolog.Event {
	if m.PanicFn != nil {
		return m.PanicFn(ctx...)
	}
	noop := zerolog.Nop()
	return noop.Panic()
}

func (m *LoggerMock) With() zerolog.Context {
	noop := zerolog.Nop()
	return noop.With()
}

func (m *LoggerMock) WithService(_ string) log.ILogger  { return m }
func (m *LoggerMock) WithModule(_ string) log.ILogger   { return m }
func (m *LoggerMock) WithBusinessID(_ uint) log.ILogger { return m }
