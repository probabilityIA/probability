package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware/internal/app"
	"github.com/secamc93/probability/back/central/services/auth/middleware/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Middleware struct {
	authService *app.AuthService
	authUseCase domain.IAuthUseCase
	logger      log.ILogger
}

func NewMiddleware(authService *app.AuthService, authUseCase domain.IAuthUseCase, logger log.ILogger) *Middleware {
	return &Middleware{
		authService: authService,
		authUseCase: authUseCase,
		logger:      logger,
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		authInfo, err := m.authService.ValidateBusinessToken(token)
		if err != nil {
			m.logger.Error().Err(err).Msg("Token inválido")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("auth_info", authInfo)
		c.Set("auth_type", authInfo.Type)
		c.Set("user_id", authInfo.UserID)

		if authInfo.BusinessTokenClaims != nil {
			c.Set("business_id", authInfo.BusinessTokenClaims.BusinessID)
			c.Set("business_type_id", authInfo.BusinessTokenClaims.BusinessTypeID)
			c.Set("role_id", authInfo.BusinessTokenClaims.RoleID)
			c.Set("business_token_claims", authInfo.BusinessTokenClaims)
			c.Set("is_super_admin", authInfo.BusinessTokenClaims.BusinessID == 0)

			if authInfo.BusinessTokenClaims.BusinessID == 0 {
				m.logger.Debug().
					Uint("user_id", authInfo.UserID).
					Msg("Business token de SUPER ADMIN validado exitosamente")
			} else {
				m.logger.Debug().
					Uint("user_id", authInfo.UserID).
					Uint("business_id", authInfo.BusinessTokenClaims.BusinessID).
					Msg("Business token validado exitosamente")
			}
		}

		c.Next()
	}
}

func (m *Middleware) BusinessTokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		authInfo, err := m.authService.ValidateMainToken(token)
		if err != nil {
			m.logger.Error().Err(err).Msg("Token inválido")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		m.logger.Debug().
			Uint("user_id", authInfo.UserID).
			Msg("Token principal validado exitosamente")

		c.Set("user_id", authInfo.UserID)
		c.Set("jwt_claims", authInfo.JWTClaims)
		c.Next()
	}
}

func (m *Middleware) APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := extractAPIKey(c)
		if apiKey == "" {
			m.logger.Error().Msg("API Key requerida")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API Key requerida",
			})
			c.Abort()
			return
		}

		request := domain.ValidateAPIKeyRequest{
			APIKey: apiKey,
		}

		response, err := m.authUseCase.ValidateAPIKey(c.Request.Context(), request)
		if err != nil {
			m.logger.Error().Err(err).Msg("Error al validar API Key")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		if !response.Success {
			m.logger.Error().Str("message", response.Message).Msg("API Key inválida")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": response.Message,
			})
			c.Abort()
			return
		}

		authInfo := &domain.AuthInfo{
			Type:       domain.AuthTypeAPIKey,
			UserID:     response.UserID,
			Email:      response.Email,
			Roles:      response.Roles,
			BusinessID: response.BusinessID,
			APIKey:     apiKey,
		}

		c.Set("auth_info", authInfo)
		c.Set("auth_type", authInfo.Type)
		c.Set("user_id", authInfo.UserID)
		c.Set("user_email", authInfo.Email)
		c.Set("user_roles", authInfo.Roles)
		c.Set("business_id", authInfo.BusinessID)
		c.Set("jwt_claims", nil)

		m.logger.Debug().
			Str("auth_type", string(authInfo.Type)).
			Uint("user_id", authInfo.UserID).
			Msg("Usuario autenticado con API Key")

		c.Next()
	}
}

func (m *Middleware) AutoAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		apiKey := extractAPIKey(c)

		if authHeader != "" {
			m.AuthMiddleware()(c)
			if c.IsAborted() {
				return
			}
			c.Next()
			return
		}
		if apiKey != "" {
			m.APIKeyMiddleware()(c)
			if c.IsAborted() {
				return
			}
			c.Next()
			return
		}

		m.logger.Error().Msg("No se encontró método de autenticación")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Se requiere autenticación (JWT o API Key)",
		})
		c.Abort()
	}
}

func extractAPIKey(c *gin.Context) string {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		apiKey = c.Query("api_key")
	}
	return apiKey
}
