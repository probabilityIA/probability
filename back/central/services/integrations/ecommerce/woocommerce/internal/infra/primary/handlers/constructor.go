package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define los endpoints HTTP de WooCommerce.
type IHandler interface {
	// HandleWebhook recibe webhooks de WooCommerce.
	HandleWebhook(c *gin.Context)
	// RegisterRoutes registra las rutas en el router.
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type wooCommerceHandler struct {
	useCase usecases.IWooCommerceUseCase
	logger  log.ILogger
}

// New crea el handler HTTP de WooCommerce.
func New(useCase usecases.IWooCommerceUseCase, logger log.ILogger) IHandler {
	return &wooCommerceHandler{
		useCase: useCase,
		logger:  logger.WithModule("woocommerce"),
	}
}

// RegisterRoutes registra las rutas de WooCommerce en el router.
func (h *wooCommerceHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	woo := router.Group("/woocommerce")
	{
		woo.POST("/webhook", h.HandleWebhook)
	}
}
