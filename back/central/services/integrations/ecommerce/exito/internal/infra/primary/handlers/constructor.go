package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/exito/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define los endpoints HTTP de Exito.
type IHandler interface {
	// HandleWebhook recibe webhooks de Exito.
	HandleWebhook(c *gin.Context)
	// RegisterRoutes registra las rutas en el router.
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type exitoHandler struct {
	useCase usecases.IExitoUseCase
	logger  log.ILogger
}

// New crea el handler HTTP de Exito.
func New(useCase usecases.IExitoUseCase, logger log.ILogger) IHandler {
	return &exitoHandler{
		useCase: useCase,
		logger:  logger.WithModule("exito"),
	}
}

// RegisterRoutes registra las rutas de Exito en el router.
func (h *exitoHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	exito := router.Group("/exito")
	{
		exito.POST("/webhook", h.HandleWebhook)
	}
}
