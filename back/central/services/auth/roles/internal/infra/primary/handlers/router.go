package rolehandler

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

func (h *RoleHandler) RegisterRoutes(router *gin.RouterGroup, handler IRoleHandler, logger log.ILogger) {
	rolesGroup := router.Group("/roles")
	{
		rolesGroup.GET("", middleware.JWT(), handler.GetRolesHandler)
		rolesGroup.POST("", middleware.JWT(), middleware.RequireSuperAdmin(), handler.CreateRole)
		rolesGroup.GET("/:id", middleware.JWT(), handler.GetRoleByIDHandler)
		rolesGroup.PUT("/:id", middleware.JWT(), middleware.RequireSuperAdmin(), handler.UpdateRole)
		rolesGroup.GET("/scope/:scope_id", middleware.JWT(), handler.GetRolesByScopeHandler)
		rolesGroup.GET("/level/:level", middleware.JWT(), handler.GetRolesByLevelHandler)
		rolesGroup.GET("/system", middleware.JWT(), handler.GetSystemRolesHandler)

		rolesGroup.POST("/:id/permissions", middleware.JWT(), middleware.RequireSuperAdmin(), handler.AssignPermissionsToRole)
		rolesGroup.GET("/:id/permissions", middleware.JWT(), handler.GetRolePermissions)
		rolesGroup.DELETE("/:id/permissions/:permission_id", middleware.JWT(), middleware.RequireSuperAdmin(), handler.RemovePermissionFromRole)
	}
}
