package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID         uint `json:"user_id"`
	BusinessID     uint `json:"business_id"`
	BusinessTypeID uint `json:"business_type_id"`
	RoleID         uint `json:"role_id"`
	jwt.RegisteredClaims
}

func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		parser := jwt.NewParser(jwt.WithLeeway(5 * time.Minute))
		token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("business_id", claims.BusinessID)
		c.Set("role_id", claims.RoleID)

		c.Next()
	}
}

func SuperAdminGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		businessID := c.GetUint("business_id")
		if businessID != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only super admins can access the testing platform"})
			c.Abort()
			return
		}
		c.Next()
	}
}
