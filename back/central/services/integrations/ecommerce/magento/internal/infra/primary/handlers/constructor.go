package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/magento/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define los endpoints HTTP de Magento.
type IHandler interface {
	// HandleWebhook recibe webhooks de Magento.
	HandleWebhook(c *gin.Context)
	// RegisterRoutes registra las rutas en el router.
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type magentoHandler struct {
	useCase usecases.IMagentoUseCase
	logger  log.ILogger
}

// New crea el handler HTTP de Magento.
func New(useCase usecases.IMagentoUseCase, logger log.ILogger) IHandler {
	return &magentoHandler{
		useCase: useCase,
		logger:  logger.WithModule("magento"),
	}
}

// RegisterRoutes registra las rutas de Magento en el router.
func (h *magentoHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	magento := router.Group("/magento")
	{
		magento.POST("/webhook", h.HandleWebhook)
	}
}
