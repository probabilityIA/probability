//go:build legacy

package routes

import (
	"central/internal/infra/primary/http2/docs"
	"central/internal/pkg/env"
	"central/internal/pkg/log"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// BasicRoutes maneja las rutas básicas del servidor y configuración de Swagger
type BasicRoutes struct {
	router *gin.Engine
	logger log.ILogger
	env    env.IConfig
}

// NewBasicRoutes crea una nueva instancia de rutas básicas
func NewBasicRoutes(router *gin.Engine, logger log.ILogger, env env.IConfig) *BasicRoutes {
	return &BasicRoutes{
		router: router,
		logger: logger,
		env:    env,
	}
}

// Setup configura las rutas básicas del servidor y Swagger
func (br *BasicRoutes) Setup() {
	// Configurar Swagger
	br.configureSwagger()

	// Configurar rutas básicas
	br.setupBasicRoutes()
}

// configureSwagger configura la documentación de Swagger
func (br *BasicRoutes) configureSwagger() {
	// Configuración básica de Swagger
	docs.SwaggerInfo.Title = "Restaurant Reservation API"
	docs.SwaggerInfo.Description = "Servicio REST para la gestión de reservas multi-restaurante."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Configurar URL base de Swagger
	swaggerBaseURL := br.env.Get("URL_BASE_SWAGGER")
	if swaggerBaseURL == "" {
		swaggerBaseURL = fmt.Sprintf("localhost:%s", br.env.Get("HTTP_PORT"))
	}

	br.logger.Info().Str("swagger_base_url", swaggerBaseURL).Msg("Configurando Swagger con URL base")

	// Procesar URL para el host
	originalURL := swaggerBaseURL
	swaggerHost := br.processSwaggerURL(swaggerBaseURL)

	docs.SwaggerInfo.Host = swaggerHost
	br.logger.Info().Str("swagger_host", swaggerHost).Str("original_url", originalURL).Msg("Swagger configurado")

	// Configurar esquemas de URL
	br.configureSwaggerSchemes()
}

// processSwaggerURL procesa la URL de Swagger para obtener el host
func (br *BasicRoutes) processSwaggerURL(swaggerBaseURL string) string {
	if strings.HasPrefix(swaggerBaseURL, "http://") {
		return strings.TrimPrefix(swaggerBaseURL, "http://")
	} else if strings.HasPrefix(swaggerBaseURL, "https://") {
		return strings.TrimPrefix(swaggerBaseURL, "https://")
	}
	return swaggerBaseURL
}

// configureSwaggerSchemes configura los esquemas de URL para Swagger
func (br *BasicRoutes) configureSwaggerSchemes() {
	if strings.HasPrefix(br.env.Get("URL_BASE_SWAGGER"), "https://") {
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	}
}

// setupBasicRoutes configura las rutas básicas del servidor
func (br *BasicRoutes) setupBasicRoutes() {
	// Ruta de health check
	br.router.GET("/ping", br.pingHandler)

	// Rutas de documentación Swagger
	br.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	br.router.GET("/docs", br.docsRedirectHandler)
}

// pingHandler maneja la ruta de health check
func (br *BasicRoutes) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// docsRedirectHandler redirige a la documentación principal
func (br *BasicRoutes) docsRedirectHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
}
