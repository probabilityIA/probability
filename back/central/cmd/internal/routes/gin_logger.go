package routes

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// GinLogger implementa el logger personalizado para Gin
type GinLogger struct {
	logger log.ILogger
}

// NewGinLogger crea una nueva instancia del logger de Gin
func NewGinLogger(logger log.ILogger) *GinLogger {
	return &GinLogger{
		logger: logger,
	}
}

// Write implementa la interfaz io.Writer para el logger de Gin
func (gl *GinLogger) Write(p []byte) (n int, err error) {
	message := strings.TrimSpace(string(p))

	// Filtrar mensajes innecesarios
	if message != "" &&
		!strings.Contains(message, "-->") &&
		!strings.Contains(message, "Running in \"debug\" mode") {
		gl.logger.Info().Msg(message)
	}

	return len(p), nil
}

// printSimpleHTTPLog imprime un log simple en consola similar a Next.js
func printSimpleHTTPLog(method, path string, status int, latency time.Duration) {
	// Color codes ANSI para status
	var statusColor string
	switch {
	case status >= 500:
		statusColor = "\033[31m" // Red
	case status >= 400:
		statusColor = "\033[33m" // Yellow
	case status >= 300:
		statusColor = "\033[36m" // Cyan
	case status >= 200:
		statusColor = "\033[32m" // Green
	default:
		statusColor = "\033[37m" // White
	}
	resetColor := "\033[0m"

	// Formatear latencia en ms
	latencyMs := float64(latency.Microseconds()) / 1000.0

	// Imprimir en formato: GET /path 200 in 123ms
	fmt.Fprintf(os.Stdout, " %s %s %s%d%s in %.0fms\n",
		method,
		path,
		statusColor,
		status,
		resetColor,
		latencyMs,
	)
	// Flush para asegurar que se imprima inmediatamente
	os.Stdout.Sync()
}

// SetupGinLogging configura el logging de Gin y un middleware HTTP mínimo
func SetupGinLogging(r *gin.Engine, logger log.ILogger) {
	// Redirigir el writer de Gin a nuestro logger personalizado
	gin.DefaultWriter = NewGinLogger(logger)

	// Middleware HTTP propio con latencia y status
	r.Use(func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		c.Next()
		status := c.Writer.Status()
		lat := time.Since(start)

		// Simple console log similar a Next.js
		printSimpleHTTPLog(method, path, status, lat)

		// Log estructurado para errores
		if status >= 400 {
			logger.Warn(c.Request.Context()).
				Str("method", method).
				Str("path", path).
				Int("status", status).
				Dur("latency", lat).
				Msg("HTTP Error")
		}
	})
}
