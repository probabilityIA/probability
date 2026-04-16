package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetSecureCookie establece una cookie segura compatible con iframe de Shopify
// Usa SameSite=None con Secure=true para permitir cookies de terceros en iframe
func SetSecureCookie(c *gin.Context, name, value string, maxAge int) {
	// Detectar si la conexión es segura (HTTPS)
	secure := c.Request.URL.Scheme == "https" ||
		c.GetHeader("X-Forwarded-Proto") == "https"

	// Detectar si es localhost
	host := c.Request.Host
	isLocal := strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1")

	var domainName string
	var sameSite http.SameSite

	// Si es localhost o dominios de desarrollo (ngrok), no forzar el dominio de producción
	if isLocal || strings.Contains(host, "ngrok") || strings.Contains(host, ".dev") {
		domainName = "" // Dejar que el navegador maneje el dominio (host-only)

		if secure {
			sameSite = http.SameSiteNoneMode // Necesario para iframes en https (ngrok)
		} else {
			sameSite = http.SameSiteLaxMode // Localhost http
		}
	} else {
		// En producción, usar el dominio raíz para compartir cookies entre subdominios
		// Aseguramos que solo se use si realmente estamos en ese dominio
		if strings.Contains(host, "probabilityia.com.co") {
			domainName = ".probabilityia.com.co"
		} else {
			domainName = "" // Fallback para otros dominios de producción o staging
		}
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   domainName,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure, // En producción (HTTPS) será true, en local (HTTP) false
		SameSite: sameSite,
	})
}
