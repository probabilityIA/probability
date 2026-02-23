package mocks

import (
	"context"
	"io"

	"github.com/rs/zerolog"
	sharedlog "github.com/secamc93/probability/back/central/shared/log"
)

// LoggerMock implementa shared/log.ILogger descartando toda salida.
// Permite que los tests no produzcan logs en consola y no fallen
// por dependencias del sistema de logging real.
type LoggerMock struct {
	zl zerolog.Logger
}

// NewLoggerMock crea un LoggerMock que descarta todos los logs.
func NewLoggerMock() sharedlog.ILogger {
	zl := zerolog.New(io.Discard)
	return &LoggerMock{zl: zl}
}

func (l *LoggerMock) Info(ctx ...context.Context) *zerolog.Event {
	return l.zl.Info()
}

func (l *LoggerMock) Error(ctx ...context.Context) *zerolog.Event {
	return l.zl.Error()
}

func (l *LoggerMock) Warn(ctx ...context.Context) *zerolog.Event {
	return l.zl.Warn()
}

func (l *LoggerMock) Debug(ctx ...context.Context) *zerolog.Event {
	return l.zl.Debug()
}

func (l *LoggerMock) Fatal(ctx ...context.Context) *zerolog.Event {
	return l.zl.WithLevel(zerolog.FatalLevel)
}

func (l *LoggerMock) Panic(ctx ...context.Context) *zerolog.Event {
	return l.zl.WithLevel(zerolog.PanicLevel)
}

func (l *LoggerMock) With() zerolog.Context {
	return l.zl.With()
}

func (l *LoggerMock) WithService(_ string) sharedlog.ILogger {
	return l
}

func (l *LoggerMock) WithModule(_ string) sharedlog.ILogger {
	return l
}

func (l *LoggerMock) WithBusinessID(_ uint) sharedlog.ILogger {
	return l
}
