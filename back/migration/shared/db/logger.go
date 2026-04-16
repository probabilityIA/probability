package db

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBLogger interface {
	LogMode(level logger.LogLevel) logger.Interface
	Info(ctx context.Context, msg string, data ...any)
	Warn(ctx context.Context, msg string, data ...any)
	Error(ctx context.Context, msg string, data ...any)
	Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error)
}

type dbLogger struct {
	level         logger.LogLevel
	slowThreshold time.Duration
	logger        zerolog.Logger
}

func NewDBLogger(log zerolog.Logger) DBLogger {
	return &dbLogger{
		level:         logger.Error,
		slowThreshold: 200 * time.Millisecond,
		logger:        log,
	}
}

func (l *dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.level = level
	return &newlogger
}

func (l *dbLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.level >= logger.Info {
		l.logger.Info().Msgf(msg, data...)
	}
}

func (l *dbLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.level >= logger.Warn {
		l.logger.Warn().Msgf(msg, data...)
	}
}

func (l *dbLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.level >= logger.Error {
		l.logger.Error().Msgf(msg, data...)
	}
}

// getFunctionName obtiene el nombre de la función que está ejecutando la consulta SQL
func getFunctionName() string {
	// Obtener el stack trace
	pc := make([]uintptr, 10)
	n := runtime.Callers(0, pc)
	if n == 0 {
		return "unknown"
	}

	// Buscar la función que no sea del paquete gorm o del logger
	for i := 0; i < n; i++ {
		fn := runtime.FuncForPC(pc[i])
		if fn == nil {
			continue
		}

		funcName := fn.Name()
		// Filtrar funciones de gorm y del logger
		if !strings.Contains(funcName, "gorm.io/gorm") &&
			!strings.Contains(funcName, "db.(*dbLogger)") &&
			!strings.Contains(funcName, "runtime.") {
			// Extraer solo el nombre de la función sin el paquete completo
			parts := strings.Split(funcName, ".")
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
			return funcName
		}
	}
	return "unknown"
}

func (l *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	functionName := getFunctionName()

	switch {
	case err != nil && l.level >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		sql, rows := fc()
		l.logger.Error().
			Err(err).
			Dur("duration", elapsed).
			Str("sql", sql).
			Int64("rows", rows).
			Str("function", functionName).
			Msg("query failed")
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.level >= logger.Warn:
		sql, rows := fc()
		l.logger.Warn().
			Dur("duration", elapsed).
			Str("sql", sql).
			Int64("rows", rows).
			Str("function", functionName).
			Msg("slow query")
	case l.level >= logger.Info:
		sql, rows := fc()
		l.logger.Info().
			Dur("duration", elapsed).
			Str("sql", sql).
			Int64("rows", rows).
			Str("function", functionName).
			Msg("SQL Query")
	}
}
