package routes

import (
	"net/http"
	"net/url"
	"strings"

	authdocs "github.com/secamc93/probability/back/central/shared/docs/auth"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupSwagger registra las rutas de Swagger UI
// Usa URL_BASE_SWAGGER para configurar Host y Schemes
func SetupSwagger(r *gin.Engine, e env.IConfig, logger log.ILogger) {
	// Configurar Host y Schemes según URL_BASE_SWAGGER
	base := e.Get("URL_BASE_SWAGGER")
	if base == "" {
		// fallback razonable
		base = "http://localhost:" + e.Get("HTTP_PORT")
	}

	// Configuración común para Auth Docs
	if u, err := url.Parse(base); err == nil && u.Host != "" {
		authdocs.SwaggerInfo.Host = u.Host
		if u.Scheme == "https" {
			authdocs.SwaggerInfo.Schemes = []string{"https"}
		} else if u.Scheme == "http" {
			authdocs.SwaggerInfo.Schemes = []string{"http"}
		}
	} else {
		authdocs.SwaggerInfo.Host = strings.TrimPrefix(strings.TrimPrefix(base, "http://"), "https://")
	}

	// BasePath por defecto
	if authdocs.SwaggerInfo.BasePath == "" {
		authdocs.SwaggerInfo.BasePath = "/api/v1"
	}

	// Registrar UI para Auth
	// La URL del JSON debe ser relativa o absoluta accesible desde el navegador
	r.GET("/docs/auth/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/auth/doc.json")))

	// Redirección de /docs a /docs/auth/index.html por defecto
	r.GET("/docs", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/docs/auth/index.html") })
}
