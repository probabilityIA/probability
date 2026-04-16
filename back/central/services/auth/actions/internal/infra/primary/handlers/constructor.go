package handlers

import (
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/app"

	"github.com/gin-gonic/gin"
)

// IActionHandler define la interfaz para el handler de Action
type IActionHandler interface {
	GetActionsHandler(c *gin.Context)
	GetActionByIDHandler(c *gin.Context)
	CreateActionHandler(c *gin.Context)
	UpdateActionHandler(c *gin.Context)
	DeleteActionHandler(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, handler IActionHandler, logger log.ILogger)
}

type ActionHandler struct {
	usecase app.IUseCaseAction
	logger  log.ILogger
}

// New crea una nueva instancia del handler de Action
func New(usecase app.IUseCaseAction, logger log.ILogger) IActionHandler {
	return &ActionHandler{
		usecase: usecase,
		logger:  logger,
	}
}
