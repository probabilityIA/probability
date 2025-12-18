package log

import (
	"os"

	"github.com/rs/zerolog"
)

type ILogger interface {
	Info() *zerolog.Event
	Error() *zerolog.Event
	Debug() *zerolog.Event
	Warn() *zerolog.Event
	Fatal() *zerolog.Event
}

type logger struct {
	log zerolog.Logger
}

func New() ILogger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	log := zerolog.New(output).With().Timestamp().Logger()

	return &logger{log: log}
}

func (l *logger) Info() *zerolog.Event {
	return l.log.Info()
}

func (l *logger) Error() *zerolog.Event {
	return l.log.Error()
}

func (l *logger) Debug() *zerolog.Event {
	return l.log.Debug()
}

func (l *logger) Warn() *zerolog.Event {
	return l.log.Warn()
}

func (l *logger) Fatal() *zerolog.Event {
	return l.log.Fatal()
}


