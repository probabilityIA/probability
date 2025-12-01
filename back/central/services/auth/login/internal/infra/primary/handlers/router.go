package authhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

func (h *AuthHandler) RegisterRoutes(v1Group *gin.RouterGroup, handler IAuthHandler, logger log.ILogger) {
	authGroup := v1Group.Group("/auth")
	{
		authGroup.POST("/login", handler.LoginHandler)
		authGroup.GET("/verify", middleware.JWT(), handler.VerifyHandler)
		authGroup.GET("/roles-permissions", middleware.JWT(), handler.GetUserRolesPermissionsHandler)
		authGroup.POST("/change-password", middleware.JWT(), handler.ChangePasswordHandler)
		authGroup.POST("/generate-password", middleware.JWT(), handler.GeneratePasswordHandler)
		authGroup.POST("/business-token", middleware.BusinessTokenAuth(), handler.GenerateBusinessTokenHandler)
	}
}
