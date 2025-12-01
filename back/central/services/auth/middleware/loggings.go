package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

func LoggingMiddleware(logger log.ILogger, env env.IConfig) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		isDebug := env.Get("LOG_LEVEL") == "debug"

		// En modo debug, loguear todas las solicitudes. De lo contrario, solo los errores.
		if isDebug || param.StatusCode >= 400 {
			var event *zerolog.Event
			if param.StatusCode >= 500 {
				event = logger.Error()
			} else if param.StatusCode >= 400 {
				event = logger.Warn()
			} else {
				event = logger.Info()
			}

			// Agregar campos estructurados al log
			event.
				Str("method", param.Method).
				Str("path", param.Path).
				Int("status_code", param.StatusCode).
				Dur("latency", param.Latency).
				Str("client_ip", param.ClientIP).
				Str("error", param.ErrorMessage).
				Msg("Solicitud HTTP procesada")
		}

		// Retornar string vacío para que Gin no loguee por su cuenta
		return ""
	})
}
