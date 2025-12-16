package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware/internal/app"
	"github.com/secamc93/probability/back/central/services/auth/middleware/internal/domain"
	httpinfra "github.com/secamc93/probability/back/central/services/auth/middleware/internal/infra/primary/http"
	jwtinfra "github.com/secamc93/probability/back/central/services/auth/middleware/internal/infra/secondary/jwt"
	"github.com/secamc93/probability/back/central/shared/env"
	sharedjwt "github.com/secamc93/probability/back/central/shared/jwt"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Re-export types for backward compatibility
type AuthType = domain.AuthType
type AuthInfo = domain.AuthInfo
type AuthError = domain.AuthError

const (
	AuthTypeUnknown = domain.AuthTypeUnknown
	AuthTypeJWT     = domain.AuthTypeJWT
	AuthTypeAPIKey  = domain.AuthTypeAPIKey
)

var (
	defaultMiddleware  *httpinfra.Middleware
	defaultJWTService  domain.IJWTService
	defaultAuthUseCase domain.IAuthUseCase
	defaultLogger      log.ILogger
	initialized        bool
)

func InitFromEnv(cfg env.IConfig, logger log.ILogger) {
	if logger == nil {
		logger = log.New()
	}
	if cfg == nil {
		cfg = env.New(logger)
	}
	secret := cfg.Get("JWT_SECRET")
	shared := sharedjwt.New(secret)

	// Create adapters
	jwtAdapter := jwtinfra.NewAdapter(shared)

	// Configure global
	Configure(jwtAdapter, nil, logger)
}

func Configure(jwtService domain.IJWTService, authUseCase domain.IAuthUseCase, logger log.ILogger) {
	defaultJWTService = jwtService
	defaultAuthUseCase = authUseCase
	defaultLogger = logger

	authService := app.NewAuthService(jwtService)
	defaultMiddleware = httpinfra.NewMiddleware(authService, authUseCase, logger)

	initialized = true
}

func ensureInitialized() {
	if !initialized {
		panic("auth middleware not configured: call middleware.Configure(...) during service bootstrap")
	}
}

func GetJWTService() domain.IJWTService {
	ensureInitialized()
	return defaultJWTService
}

// Middlewares
func JWT() gin.HandlerFunc {
	ensureInitialized()
	return defaultMiddleware.AuthMiddleware()
}

func APIKey() gin.HandlerFunc {
	ensureInitialized()
	return defaultMiddleware.APIKeyMiddleware()
}

func Auto() gin.HandlerFunc {
	ensureInitialized()
	return defaultMiddleware.AutoAuthMiddleware()
}

func BusinessTokenAuth() gin.HandlerFunc {
	ensureInitialized()
	return defaultMiddleware.BusinessTokenAuthMiddleware()
}

// Deprecated: Use Configure and JWT() instead, or use this for custom instances
func AuthMiddleware(jwtService domain.IJWTService, logger log.ILogger) gin.HandlerFunc {
	svc := app.NewAuthService(jwtService)
	mw := httpinfra.NewMiddleware(svc, nil, logger)
	return mw.AuthMiddleware()
}

// Deprecated: Use Configure and APIKey() instead
func APIKeyMiddleware(authUseCase domain.IAuthUseCase, logger log.ILogger) gin.HandlerFunc {
	mw := httpinfra.NewMiddleware(nil, authUseCase, logger)
	return mw.APIKeyMiddleware()
}

// Deprecated: Use Configure and Auto() instead
func AutoAuthMiddleware(jwtService domain.IJWTService, authUseCase domain.IAuthUseCase, logger log.ILogger) gin.HandlerFunc {
	svc := app.NewAuthService(jwtService)
	mw := httpinfra.NewMiddleware(svc, authUseCase, logger)
	return mw.AutoAuthMiddleware()
}

// Helpers
func GetAuthInfo(c *gin.Context) (*domain.AuthInfo, bool) {
	authInfo, exists := c.Get("auth_info")
	if !exists {
		return nil, false
	}
	if a, ok := authInfo.(*domain.AuthInfo); ok {
		return a, true
	}
	return nil, false
}

func GetAuthType(c *gin.Context) (domain.AuthType, bool) {
	authType, exists := c.Get("auth_type")
	if !exists {
		return domain.AuthTypeUnknown, false
	}
	if a, ok := authType.(domain.AuthType); ok {
		return a, true
	}
	return domain.AuthTypeUnknown, false
}

func GetAPIKey(c *gin.Context) (string, bool) {
	authInfo, exists := GetAuthInfo(c)
	if !exists || authInfo.Type != domain.AuthTypeAPIKey {
		return "", false
	}
	return authInfo.APIKey, true
}

func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	if id, ok := userID.(uint); ok {
		return id, true
	}
	return 0, false
}

func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}
	if e, ok := email.(string); ok {
		return e, true
	}
	return "", false
}

func GetUserRoles(c *gin.Context) ([]string, bool) {
	roles, exists := c.Get("user_roles")
	if !exists {
		return nil, false
	}
	if r, ok := roles.([]string); ok {
		return r, true
	}
	return nil, false
}

func GetJWTClaims(c *gin.Context) (*domain.JWTClaims, bool) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		return nil, false
	}
	if c, ok := claims.(*domain.JWTClaims); ok {
		return c, true
	}
	return nil, false
}

func GetBusinessID(c *gin.Context) (uint, bool) {
	authInfo, exists := GetAuthInfo(c)
	if !exists {
		return 0, false
	}
	return authInfo.BusinessID, true
}

func GetBusinessIDFromContext(c *gin.Context) (uint, bool) {
	businessID, exists := c.Get("business_id")
	if !exists {
		return 0, false
	}
	if id, ok := businessID.(uint); ok {
		return id, true
	}
	return 0, false
}

func GetBusinessTokenClaims(c *gin.Context) (*domain.BusinessTokenClaims, bool) {
	claims, exists := c.Get("business_token_claims")
	if !exists {
		return nil, false
	}
	if c, ok := claims.(*domain.BusinessTokenClaims); ok {
		return c, true
	}
	return nil, false
}

func GetBusinessTypeID(c *gin.Context) (uint, bool) {
	typeID, exists := c.Get("business_type_id")
	if !exists {
		return 0, false
	}
	if id, ok := typeID.(uint); ok {
		return id, true
	}
	return 0, false
}

func GetRoleID(c *gin.Context) (uint, bool) {
	roleID, exists := c.Get("role_id")
	if !exists {
		return 0, false
	}
	if id, ok := roleID.(uint); ok {
		return id, true
	}
	return 0, false
}

func IsSuperAdmin(c *gin.Context) bool {
	businessID, exists := c.Get("business_id")
	if !exists {
		return false
	}
	if id, ok := businessID.(uint); ok {
		return id == 0
	}
	return false
}

// GetScope retorna el scope del usuario autenticado ("platform" o "business")
func GetScope(c *gin.Context) string {
	// Primero intentar obtener del auth_info
	authInfo, exists := GetAuthInfo(c)
	if exists && authInfo.Scope != "" {
		return authInfo.Scope
	}

	// Si no hay scope en auth_info, determinar por business_id
	// Si business_id == 0, es scope platform (super admin)
	// Si business_id > 0, es scope business
	businessID, exists := GetBusinessID(c)
	if exists {
		if businessID == 0 {
			return "platform"
		}
		return "business"
	}

	// Por defecto, retornar business
	return "business"
}

// RequireRole and RequireAnyRole logic
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := GetUserRoles(c)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: roles no disponibles",
			})
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: rol requerido",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := GetUserRoles(c)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: roles no disponibles",
			})
			c.Abort()
			return
		}

		hasRole := false
		for _, userRole := range roles {
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: rol requerido",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireAuthType(authType domain.AuthType) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentAuthType, exists := GetAuthType(c)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: tipo de autenticación no disponible",
			})
			c.Abort()
			return
		}

		if currentAuthType != authType {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: tipo de autenticación requerido",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireJWT() gin.HandlerFunc {
	return RequireAuthType(domain.AuthTypeJWT)
}

func RequireAPIKey() gin.HandlerFunc {
	return RequireAuthType(domain.AuthTypeAPIKey)
}

// AuthBuilder
type AuthBuilder struct {
	jwtService  domain.IJWTService
	authUseCase domain.IAuthUseCase
	logger      log.ILogger
}

func NewAuthBuilder(jwtService domain.IJWTService, authUseCase domain.IAuthUseCase, logger log.ILogger) *AuthBuilder {
	return &AuthBuilder{
		jwtService:  jwtService,
		authUseCase: authUseCase,
		logger:      logger,
	}
}

func (ab *AuthBuilder) JWT() gin.HandlerFunc {
	return AuthMiddleware(ab.jwtService, ab.logger)
}

func (ab *AuthBuilder) APIKey() gin.HandlerFunc {
	return APIKeyMiddleware(ab.authUseCase, ab.logger)
}

func (ab *AuthBuilder) Auto() gin.HandlerFunc {
	return AutoAuthMiddleware(ab.jwtService, ab.authUseCase, ab.logger)
}
