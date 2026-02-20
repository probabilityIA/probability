package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// fmt.Printf("CORS Debug: Origin=%s Method=%s\n", origin, c.Request.Method)

		allowedOrigins := []string{
			"https://www.probabilityia.com.co",
			"https://admin.shopify.com",
			"http://localhost:3000",
			"http://localhost:3001",
		}

		isShopifyOrigin := strings.HasSuffix(origin, ".myshopify.com")
		isAllowed := false

		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}

		// Permissive localhost/dev check
		if strings.HasPrefix(origin, "http://localhost") ||
			strings.HasPrefix(origin, "http://127.0.0.1") ||
			strings.Contains(origin, "ngrok") ||
			origin == "null" {
			isAllowed = true
		}

		if isAllowed || isShopifyOrigin {
			finalOrigin := origin
			// Fix for null origin with credentials
			if finalOrigin == "null" {
				finalOrigin = "http://localhost:3000"
			}
			c.Header("Access-Control-Allow-Origin", finalOrigin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With, X-API-Key, sh-token")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
