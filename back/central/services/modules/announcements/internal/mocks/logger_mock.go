package mocks

import (
	"context"
	"io"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/shared/log"
)

type LoggerMock struct{}

var _ log.ILogger = (*LoggerMock)(nil)

func nopEvent() *zerolog.Event {
	nop := zerolog.New(io.Discard)
	return nop.Log()
}

func (m *LoggerMock) Info(ctx ...context.Context) *zerolog.Event  { return nopEvent() }
func (m *LoggerMock) Error(ctx ...context.Context) *zerolog.Event { return nopEvent() }
func (m *LoggerMock) Warn(ctx ...context.Context) *zerolog.Event  { return nopEvent() }
func (m *LoggerMock) Debug(ctx ...context.Context) *zerolog.Event { return nopEvent() }
func (m *LoggerMock) Fatal(ctx ...context.Context) *zerolog.Event { return nopEvent() }
func (m *LoggerMock) Panic(ctx ...context.Context) *zerolog.Event { return nopEvent() }

func (m *LoggerMock) With() zerolog.Context {
	nop := zerolog.New(io.Discard)
	return nop.With()
}

func (m *LoggerMock) WithService(service string) log.ILogger  { return m }
func (m *LoggerMock) WithModule(module string) log.ILogger    { return m }
func (m *LoggerMock) WithBusinessID(id uint) log.ILogger      { return m }
