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

	if isLocal {
		domainName = ""                 // Dejar que el navegador maneje el dominio en localhost
		sameSite = http.SameSiteLaxMode // Lax es suficiente para localhost y no requiere Secure
	} else {
		domainName = ".probabilityia.com.co"
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
