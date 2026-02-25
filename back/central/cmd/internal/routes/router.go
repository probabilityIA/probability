package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// BuildRouter construye y configura el *gin.Engine del monolito en un solo lugar
func BuildRouter(ctx context.Context, logger log.ILogger, environment env.IConfig) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// CORS - DEBE IR PRIMERO
	r.Use(middleware.CorsMiddleware())

	// Security headers para iframe de Shopify
	r.Use(middleware.SecurityHeadersMiddleware())

	// Prometheus metrics middleware
	r.Use(metrics.PrometheusMiddleware())

	// Logging centralizado
	SetupGinLogging(r, logger)

	// Recovery
	r.Use(gin.Recovery())

	// Metrics endpoint (sin autenticación - solo métricas técnicas)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "central-reserve",
			"version":   "1.0.0",
		})
	})

	// Test endpoint
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// 404 JSON explícito + log WARN
	r.NoRoute(func(c *gin.Context) {
		logger.Warn(ctx).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", 404).
			Msg("Ruta no encontrada")
		c.JSON(404, gin.H{"error": "not_found"})
	})

	return r
}
