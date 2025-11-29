package routes

import (
	"central/shared/env"
	"central/shared/log"
	"net/http"
	"net/url"
	"strings"

	docs "central/shared/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupSwagger registra las rutas de Swagger UI
// Usa URL_BASE_SWAGGER para configurar Host y Schemes; JSON por defecto en /docs/doc.json
func SetupSwagger(r *gin.Engine, e env.IConfig, logger log.ILogger) {
	// Configurar Host y Schemes según URL_BASE_SWAGGER
	base := e.Get("URL_BASE_SWAGGER")
	if base == "" {
		// fallback razonable
		base = "http://localhost:" + e.Get("HTTP_PORT")
	}
	if u, err := url.Parse(base); err == nil && u.Host != "" {
		docs.SwaggerInfo.Host = u.Host
		if u.Scheme == "https" {
			docs.SwaggerInfo.Schemes = []string{"https"}
		} else if u.Scheme == "http" {
			docs.SwaggerInfo.Schemes = []string{"http"}
		}
	} else {
		// Si no es URL completa, asume host:port sin esquema
		docs.SwaggerInfo.Host = strings.TrimPrefix(strings.TrimPrefix(base, "http://"), "https://")
		// Mantener esquemas por defecto del docs
	}

	// BasePath por defecto si está vacío
	if docs.SwaggerInfo.BasePath == "" {
		docs.SwaggerInfo.BasePath = "/api/v1"
	}

	// Registrar UI apuntando al JSON
	jsonURL := e.Get("SWAGGER_JSON_URL")
	if jsonURL == "" {
		jsonURL = "/docs/doc.json"
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(jsonURL)))
	r.GET("/docs", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/docs/index.html") })
}
