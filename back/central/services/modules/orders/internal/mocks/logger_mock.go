package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/stretchr/testify/mock"
)

// LoggerMock es un mock del logger
type LoggerMock struct {
	mock.Mock
}

func (m *LoggerMock) Info(ctx ...context.Context) *zerolog.Event {
	m.Called(ctx)
	// Retornar un evento disabled que no hace nada
	nop := zerolog.Nop()
	return nop.Info()
}

func (m *LoggerMock) Error(ctx ...context.Context) *zerolog.Event {
	m.Called(ctx)
	nop := zerolog.Nop()
	return nop.Error()
}

func (m *LoggerMock) Debug(ctx ...context.Context) *zerolog.Event {
	m.Called(ctx)
	nop := zerolog.Nop()
	return nop.Debug()
}

func (m *LoggerMock) Warn(ctx ...context.Context) *zerolog.Event {
	m.Called(ctx)
	nop := zerolog.Nop()
	return nop.Warn()
}

func (m *LoggerMock) Fatal(ctx ...context.Context) *zerolog.Event {
	m.Called(ctx)
	nop := zerolog.Nop()
	return nop.Fatal()
}

func (m *LoggerMock) Panic(ctx ...context.Context) *zerolog.Event {
	m.Called(ctx)
	nop := zerolog.Nop()
	return nop.Panic()
}

func (m *LoggerMock) With() zerolog.Context {
	m.Called()
	nop := zerolog.Nop()
	return nop.With()
}

func (m *LoggerMock) WithService(service string) log.ILogger {
	m.Called(service)
	return m
}

func (m *LoggerMock) WithModule(module string) log.ILogger {
	m.Called(module)
	return m
}

func (m *LoggerMock) WithBusinessID(businessID uint) log.ILogger {
	m.Called(businessID)
	return m
}
