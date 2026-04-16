package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Ihandlers define la interfaz del handler de usuarios
type Ihandlers interface {
	GetUsersHandler(c *gin.Context)
	GetUserByIDHandler(c *gin.Context)
	Createhandlers(c *gin.Context)
	Updatehandlers(c *gin.Context)
	Deletehandlers(c *gin.Context)
	AssignRoleToUserBusinessHandler(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, handler Ihandlers, logger log.ILogger)
}

// handlers maneja las solicitudes HTTP para usuarios
type handlers struct {
	usecase app.Iapp
	logger  log.ILogger
}

// New crea una nueva instancia del handler de usuarios
func New(usecase app.Iapp, logger log.ILogger) Ihandlers {
	return &handlers{
		usecase: usecase,
		logger:  logger,
	}
}
