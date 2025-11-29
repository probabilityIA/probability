package routes

import (
	"central/internal/pkg/log"
	"strings"
	"time"

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

// SetupGinLogging configura el logging de Gin y un middleware HTTP mínimo
func SetupGinLogging(r *gin.Engine, logger log.ILogger) {
	// Redirigir el writer de Gin a nuestro logger personalizado
	gin.DefaultWriter = NewGinLogger(logger)

	// Middleware de log de Gin
	r.Use(gin.Logger())

	// Middleware HTTP propio con latencia y status
	r.Use(func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		c.Next()
		status := c.Writer.Status()
		lat := time.Since(start)
		logger.Info(c.Request.Context()).
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("latency", lat).
			Msg("HTTP")
	})
}
