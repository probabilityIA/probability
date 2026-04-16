package permissionhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

// RegisterRoutes registra las rutas para el handler de permisos
func (h *PermissionHandler) RegisterRoutes(router *gin.RouterGroup, handler IPermissionHandler, logger log.ILogger) {
	permissionsGroup := router.Group("/permissions")
	{
		permissionsGroup.GET("", middleware.JWT(), handler.GetPermissionsHandler)
		permissionsGroup.GET("/:id", middleware.JWT(), handler.GetPermissionByIDHandler)
		permissionsGroup.GET("/scope/:scope_id", middleware.JWT(), handler.GetPermissionsByScopeHandler)
		permissionsGroup.GET("/resource/:resource", middleware.JWT(), handler.GetPermissionsByResourceHandler)
		permissionsGroup.POST("", middleware.JWT(), handler.CreatePermissionHandler)
		permissionsGroup.POST("/bulk", middleware.JWT(), handler.BulkCreatePermissionsHandler)
		permissionsGroup.PUT("/:id", middleware.JWT(), handler.UpdatePermissionHandler)
		permissionsGroup.DELETE("/:id", middleware.JWT(), handler.DeletePermissionHandler)
	}
}
