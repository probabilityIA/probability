package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/falabella/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define los endpoints HTTP de Falabella.
type IHandler interface {
	// HandleWebhook recibe webhooks de Falabella Seller Center.
	HandleWebhook(c *gin.Context)
	// RegisterRoutes registra las rutas en el router.
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type falabellaHandler struct {
	useCase usecases.IFalabellaUseCase
	logger  log.ILogger
}

// New crea el handler HTTP de Falabella.
func New(useCase usecases.IFalabellaUseCase, logger log.ILogger) IHandler {
	return &falabellaHandler{
		useCase: useCase,
		logger:  logger.WithModule("falabella"),
	}
}

// RegisterRoutes registra las rutas de Falabella en el router.
func (h *falabellaHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	falabella := router.Group("/falabella")
	{
		falabella.POST("/webhook", h.HandleWebhook)
	}
}
