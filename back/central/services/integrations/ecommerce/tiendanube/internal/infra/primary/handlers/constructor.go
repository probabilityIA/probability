package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define los endpoints HTTP de Tiendanube.
type IHandler interface {
	// HandleWebhook recibe webhooks de Tiendanube.
	HandleWebhook(c *gin.Context)
	// RegisterRoutes registra las rutas en el router.
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type tiendanubeHandler struct {
	useCase usecases.ITiendanubeUseCase
	logger  log.ILogger
}

// New crea el handler HTTP de Tiendanube.
func New(useCase usecases.ITiendanubeUseCase, logger log.ILogger) IHandler {
	return &tiendanubeHandler{
		useCase: useCase,
		logger:  logger.WithModule("tiendanube"),
	}
}

// RegisterRoutes registra las rutas de Tiendanube en el router.
func (h *tiendanubeHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	tn := router.Group("/tiendanube")
	{
		tn.POST("/webhook", h.HandleWebhook)
	}
}
