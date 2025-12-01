package authhandler

import (
	"github.com/secamc93/probability/back/central/services/auth/login/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// IAuthHandler define la interfaz del handler de autenticación
type IAuthHandler interface {
	LoginHandler(c *gin.Context)
	VerifyHandler(c *gin.Context)
	GetUserRolesPermissionsHandler(c *gin.Context)
	ChangePasswordHandler(c *gin.Context)
	GeneratePasswordHandler(c *gin.Context)
	GenerateBusinessTokenHandler(c *gin.Context)
	RegisterRoutes(v1Group *gin.RouterGroup, handler IAuthHandler, logger log.ILogger)
	// GenerateAPIKeyHandler(c *gin.Context)
	// ValidateAPIKeyHandler(c *gin.Context)
}

type AuthHandler struct {
	usecase app.Iapp
	logger  log.ILogger
}

// New crea una nueva instancia del handler de autenticación
func New(usecase app.Iapp, logger log.ILogger) IAuthHandler {
	contextualLogger := logger.WithModule("autenticación")
	return &AuthHandler{
		usecase: usecase,
		logger:  contextualLogger,
	}
}
