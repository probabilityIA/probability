package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/shared/log"
)

type LoggerMock struct{}

func NewLoggerMock() log.ILogger {
	return &LoggerMock{}
}

func (m *LoggerMock) Info(ctx ...context.Context) *zerolog.Event {
	noop := zerolog.Nop()
	return noop.Info()
}

func (m *LoggerMock) Error(ctx ...context.Context) *zerolog.Event {
	noop := zerolog.Nop()
	return noop.Error()
}

func (m *LoggerMock) Warn(ctx ...context.Context) *zerolog.Event {
	noop := zerolog.Nop()
	return noop.Warn()
}

func (m *LoggerMock) Debug(ctx ...context.Context) *zerolog.Event {
	noop := zerolog.Nop()
	return noop.Debug()
}

func (m *LoggerMock) Fatal(ctx ...context.Context) *zerolog.Event {
	noop := zerolog.Nop()
	return noop.Fatal()
}

func (m *LoggerMock) Panic(ctx ...context.Context) *zerolog.Event {
	noop := zerolog.Nop()
	return noop.Panic()
}

func (m *LoggerMock) With() zerolog.Context {
	noop := zerolog.Nop()
	return noop.With()
}

func (m *LoggerMock) WithService(service string) log.ILogger {
	return m
}

func (m *LoggerMock) WithModule(module string) log.ILogger {
	return m
}

func (m *LoggerMock) WithBusinessID(businessID uint) log.ILogger {
	return m
}
