package rolehandler

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IRoleHandler define la interfaz del handler de roles
type IRoleHandler interface {
	GetRolesHandler(c *gin.Context)
	GetRoleByIDHandler(c *gin.Context)
	GetRolesByScopeHandler(c *gin.Context)
	GetRolesByLevelHandler(c *gin.Context)
	GetSystemRolesHandler(c *gin.Context)
	CreateRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	AssignPermissionsToRole(c *gin.Context)
	RemovePermissionFromRole(c *gin.Context)
	GetRolePermissions(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, handler IRoleHandler, logger log.ILogger)
}

// RoleHandler maneja las solicitudes HTTP para roles
type RoleHandler struct {
	usecase app.IUseCaseRole
	logger  log.ILogger
}

// NewRoleHandler crea una nueva instancia del handler de roles
func New(usecase app.IUseCaseRole, logger log.ILogger) IRoleHandler {
	return &RoleHandler{
		usecase: usecase,
		logger:  logger,
	}
}
