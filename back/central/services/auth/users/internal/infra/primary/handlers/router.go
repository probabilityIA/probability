package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

// RegisterRoutes registra las rutas del handler de usuarios
func (h *handlers) RegisterRoutes(router *gin.RouterGroup, handler Ihandlers, logger log.ILogger) {
	usersGroup := router.Group("/users")

	{
		usersGroup.GET("", middleware.JWT(), handler.GetUsersHandler)
		usersGroup.GET("/:id", middleware.JWT(), handler.GetUserByIDHandler)
		usersGroup.POST("", middleware.JWT(), handler.Createhandlers)
		usersGroup.PUT("/:id", middleware.JWT(), handler.Updatehandlers)
		usersGroup.DELETE("/:id", middleware.JWT(), handler.Deletehandlers)
		usersGroup.POST("/:id/assign-role", middleware.JWT(), handler.AssignRoleToUserBusinessHandler)
	}
}
