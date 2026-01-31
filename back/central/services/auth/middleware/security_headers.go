package middleware

import "github.com/gin-gonic/gin"

// SecurityHeadersMiddleware agrega headers de seguridad para proteger la aplicación
// cuando se ejecuta en iframe de Shopify
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Permitir iframe solo desde Shopify
		c.Header("Content-Security-Policy",
			"frame-ancestors https://admin.shopify.com https://*.myshopify.com")

		// HTTPS enforcement - forzar HTTPS en producción
		c.Header("Strict-Transport-Security",
			"max-age=63072000; includeSubDomains; preload")

		// Prevenir MIME sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		c.Next()
	}
}
