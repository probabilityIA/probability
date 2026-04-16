package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware agrega headers de seguridad para proteger la aplicación
// cuando se ejecuta en iframe de Shopify
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el dominio de la tienda de Shopify
		// Shopify envía esto en query params o headers
		shop := c.Query("shop")

		if shop == "" {
			shop = c.GetHeader("X-Shopify-Shop-Domain")
		}

		// Construir CSP dinámico basado en la tienda específica
		var csp string
		if shop != "" {
			// CRÍTICO: Validar que el shop termine en .myshopify.com para prevenir ataques
			if strings.HasSuffix(shop, ".myshopify.com") {
				// Frame-ancestors DEBE ser específico por tienda (NO wildcards)
				// Ref: https://shopify.dev/docs/apps/build/security/set-up-iframe-protection
				csp = fmt.Sprintf("frame-ancestors https://%s https://admin.shopify.com", shop)
			} else {
				// Si el dominio no es válido, denegar todo framing
				csp = "frame-ancestors 'none'"
			}
		} else {
			// Si no hay shop, permitir solo admin.shopify.com (para OAuth)
			// O denegar completamente si no es una ruta de Shopify
			csp = "frame-ancestors https://admin.shopify.com"
		}

		c.Header("Content-Security-Policy", csp)

		// HTTPS enforcement - forzar HTTPS en producción
		c.Header("Strict-Transport-Security",
			"max-age=63072000; includeSubDomains; preload")

		// Prevenir MIME sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevenir clickjacking (legacy fallback para navegadores viejos)
		// X-Frame-Options es obsoleto pero algunos navegadores viejos lo necesitan
		if shop != "" && strings.HasSuffix(shop, ".myshopify.com") {
			c.Header("X-Frame-Options", fmt.Sprintf("ALLOW-FROM https://%s", shop))
		}

		c.Next()
	}
}
