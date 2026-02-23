package mocks

import (
	"context"
	"io"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Compilación estática: LoggerMock debe satisfacer log.ILogger.
var _ log.ILogger = (*LoggerMock)(nil)

// LoggerMock es el mock de log.ILogger.
// Descarta toda la salida escribiendo a io.Discard para evitar ruido en los tests.
type LoggerMock struct {
	zl zerolog.Logger
}

// NewLoggerMock crea un LoggerMock que descarta todos los mensajes.
func NewLoggerMock() *LoggerMock {
	return &LoggerMock{
		zl: zerolog.New(io.Discard),
	}
}

func (m *LoggerMock) Info(ctx ...context.Context) *zerolog.Event {
	return m.zl.Info()
}

func (m *LoggerMock) Error(ctx ...context.Context) *zerolog.Event {
	return m.zl.Error()
}

func (m *LoggerMock) Warn(ctx ...context.Context) *zerolog.Event {
	return m.zl.Warn()
}

func (m *LoggerMock) Debug(ctx ...context.Context) *zerolog.Event {
	return m.zl.Debug()
}

func (m *LoggerMock) Fatal(ctx ...context.Context) *zerolog.Event {
	return m.zl.WithLevel(zerolog.FatalLevel)
}

func (m *LoggerMock) Panic(ctx ...context.Context) *zerolog.Event {
	return m.zl.WithLevel(zerolog.PanicLevel)
}

func (m *LoggerMock) With() zerolog.Context {
	return m.zl.With()
}

func (m *LoggerMock) WithService(_ string) log.ILogger {
	return m
}

func (m *LoggerMock) WithModule(_ string) log.ILogger {
	return m
}

func (m *LoggerMock) WithBusinessID(_ uint) log.ILogger {
	return m
}
