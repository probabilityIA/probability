package permissionhandler

import (
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// IPermissionHandler define la interfaz para el handler de Permission
type IPermissionHandler interface {
	GetPermissionsHandler(c *gin.Context)
	GetPermissionByIDHandler(c *gin.Context)
	GetPermissionsByScopeHandler(c *gin.Context)
	GetPermissionsByResourceHandler(c *gin.Context)
	CreatePermissionHandler(c *gin.Context)
	BulkCreatePermissionsHandler(c *gin.Context)
	UpdatePermissionHandler(c *gin.Context)
	DeletePermissionHandler(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, handler IPermissionHandler, logger log.ILogger)
}

type PermissionHandler struct {
	usecase app.Iapp
	logger  log.ILogger
}

// New crea una nueva instancia del handler de Permission
func New(usecase app.Iapp, logger log.ILogger) IPermissionHandler {
	return &PermissionHandler{
		usecase: usecase,
		logger:  logger,
	}
}
