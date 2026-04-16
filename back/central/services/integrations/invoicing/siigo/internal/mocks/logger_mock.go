package mocks

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/shared/log"
)

// LoggerMock implementa log.ILogger descartando todos los logs durante los tests.
// Retorna eventos zerolog válidos que no producen salida, permitiendo que el
// código bajo test llame a métodos como .Str(), .Bool(), .Err(), .Msg() sin panic.
type LoggerMock struct{}

// newDiscardEvent construye un *zerolog.Event conectado a io.Discard.
// Esto es necesario porque zerolog.Event es una struct concreta sin interfaz,
// por lo que necesitamos un evento real (pero que no escribe nada).
func newDiscardEvent(level zerolog.Level) *zerolog.Event {
	logger := zerolog.New(io.Discard)
	switch level {
	case zerolog.InfoLevel:
		return logger.Info()
	case zerolog.WarnLevel:
		return logger.Warn()
	case zerolog.ErrorLevel:
		return logger.Error()
	case zerolog.DebugLevel:
		return logger.Debug()
	case zerolog.FatalLevel:
		// Usamos WithLevel para evitar que os.Exit sea llamado en tests
		return logger.WithLevel(zerolog.FatalLevel)
	case zerolog.PanicLevel:
		return logger.WithLevel(zerolog.PanicLevel)
	default:
		return logger.Info()
	}
}

// Info implementa log.ILogger.
func (l *LoggerMock) Info(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.InfoLevel)
}

// Error implementa log.ILogger.
func (l *LoggerMock) Error(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.ErrorLevel)
}

// Warn implementa log.ILogger.
func (l *LoggerMock) Warn(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.WarnLevel)
}

// Debug implementa log.ILogger.
func (l *LoggerMock) Debug(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.DebugLevel)
}

// Fatal implementa log.ILogger.
// Usa WithLevel en lugar de Fatal() para evitar llamar a os.Exit en tests.
func (l *LoggerMock) Fatal(ctx ...context.Context) *zerolog.Event {
	logger := zerolog.New(os.Stderr)
	return logger.WithLevel(zerolog.FatalLevel)
}

// Panic implementa log.ILogger.
func (l *LoggerMock) Panic(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.PanicLevel)
}

// With implementa log.ILogger.
func (l *LoggerMock) With() zerolog.Context {
	return zerolog.New(io.Discard).With()
}

// WithService implementa log.ILogger — retorna el mismo mock.
func (l *LoggerMock) WithService(service string) log.ILogger {
	return l
}

// WithModule implementa log.ILogger — retorna el mismo mock.
func (l *LoggerMock) WithModule(module string) log.ILogger {
	return l
}

// WithBusinessID implementa log.ILogger — retorna el mismo mock.
func (l *LoggerMock) WithBusinessID(businessID uint) log.ILogger {
	return l
}
