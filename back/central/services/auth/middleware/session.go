package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetSecureCookie establece una cookie segura compatible con iframe de Shopify
// Usa SameSite=None con Secure=true para permitir cookies de terceros en iframe
func SetSecureCookie(c *gin.Context, name, value string, maxAge int) {
	// Detectar si la conexión es segura (HTTPS)
	secure := c.Request.URL.Scheme == "https" ||
		c.GetHeader("X-Forwarded-Proto") == "https"

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   ".probabilityia.com.co", // Con punto inicial para subdominios y iframes
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteNoneMode, // CRÍTICO para iframe - permite cookies de terceros
	})
}
