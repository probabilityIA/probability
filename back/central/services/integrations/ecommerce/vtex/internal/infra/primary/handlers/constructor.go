package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define los endpoints HTTP de VTEX.
type IHandler interface {
	// HandleWebhook recibe webhooks de VTEX.
	HandleWebhook(c *gin.Context)
	// RegisterRoutes registra las rutas en el router.
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type vtexHandler struct {
	useCase usecases.IVTEXUseCase
	logger  log.ILogger
}

// New crea el handler HTTP de VTEX.
func New(useCase usecases.IVTEXUseCase, logger log.ILogger) IHandler {
	return &vtexHandler{
		useCase: useCase,
		logger:  logger.WithModule("vtex"),
	}
}

// RegisterRoutes registra las rutas de VTEX en el router.
func (h *vtexHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	vtex := router.Group("/vtex")
	{
		vtex.POST("/webhook", h.HandleWebhook)
	}
}
