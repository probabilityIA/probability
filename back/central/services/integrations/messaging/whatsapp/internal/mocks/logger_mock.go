package mocks

import (
	"context"
	"io"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/shared/log"
)

// discardLog es un logger de zerolog que descarta todos los mensajes.
// Se usa como singleton para construir eventos sin salida.
var discardLog = zerolog.New(io.Discard)

// LoggerMock implementa log.ILogger para tests unitarios (descarta todos los logs)
type LoggerMock struct{}

// Verificar en tiempo de compilación que implementa la interfaz
var _ log.ILogger = (*LoggerMock)(nil)

func (l *LoggerMock) Info(ctx ...context.Context) *zerolog.Event {
	return discardLog.Info()
}

func (l *LoggerMock) Error(ctx ...context.Context) *zerolog.Event {
	return discardLog.Error()
}

func (l *LoggerMock) Warn(ctx ...context.Context) *zerolog.Event {
	return discardLog.Warn()
}

func (l *LoggerMock) Debug(ctx ...context.Context) *zerolog.Event {
	return discardLog.Debug()
}

func (l *LoggerMock) Fatal(ctx ...context.Context) *zerolog.Event {
	// En tests, Fatal no debe terminar el proceso
	return discardLog.WithLevel(zerolog.NoLevel)
}

func (l *LoggerMock) Panic(ctx ...context.Context) *zerolog.Event {
	return discardLog.WithLevel(zerolog.NoLevel)
}

func (l *LoggerMock) With() zerolog.Context {
	return discardLog.With()
}

func (l *LoggerMock) WithService(service string) log.ILogger {
	return l
}

func (l *LoggerMock) WithModule(module string) log.ILogger {
	return l
}

func (l *LoggerMock) WithBusinessID(businessID uint) log.ILogger {
	return l
}
