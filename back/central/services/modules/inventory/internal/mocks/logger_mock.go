package mocks

import (
	"context"
	"io"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/shared/log"
)

// LoggerMock implementa log.ILogger descartando todos los mensajes
type LoggerMock struct{}

// newDiscardEvent crea un zerolog.Event que descarta la salida
func newDiscardEvent(level zerolog.Level) *zerolog.Event {
	logger := zerolog.New(io.Discard).Level(level)
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
		return logger.WithLevel(zerolog.FatalLevel)
	case zerolog.PanicLevel:
		return logger.WithLevel(zerolog.PanicLevel)
	default:
		return logger.Info()
	}
}

func (l *LoggerMock) Info(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.InfoLevel)
}

func (l *LoggerMock) Error(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.ErrorLevel)
}

func (l *LoggerMock) Warn(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.WarnLevel)
}

func (l *LoggerMock) Debug(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.DebugLevel)
}

func (l *LoggerMock) Fatal(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.FatalLevel)
}

func (l *LoggerMock) Panic(ctx ...context.Context) *zerolog.Event {
	return newDiscardEvent(zerolog.PanicLevel)
}

func (l *LoggerMock) With() zerolog.Context {
	return zerolog.New(io.Discard).With()
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
