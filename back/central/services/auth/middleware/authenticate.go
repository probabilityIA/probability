package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware/internal/app"
)

func Authenticate(c *gin.Context) bool {
	ensureInitialized()

	token, err := c.Cookie("session_token")
	if err != nil || token == "" {
		token = c.GetHeader("Authorization")
	}
	if token == "" {
		return false
	}

	authInfo, err := app.NewAuthService(defaultJWTService).ValidateBusinessToken(token)
	if err != nil || authInfo == nil {
		return false
	}

	c.Set("auth_info", authInfo)
	c.Set("auth_type", authInfo.Type)
	c.Set("user_id", authInfo.UserID)
	c.Set("user_email", authInfo.Email)
	c.Set("business_id", authInfo.BusinessID)
	c.Set("business_type_id", authInfo.BusinessTypeID)
	c.Set("role_id", authInfo.RoleID)
	c.Set("business_token_claims", authInfo.BusinessTokenClaims)
	c.Set("jwt_claims", authInfo.JWTClaims)
	c.Set("is_super_admin", authInfo.BusinessID == 0)
	return true
}
