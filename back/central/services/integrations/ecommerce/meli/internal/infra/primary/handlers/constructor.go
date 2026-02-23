package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define los endpoints HTTP de MercadoLibre.
type IHandler interface {
	// HandleNotification recibe notificaciones IPN de MercadoLibre.
	HandleNotification(c *gin.Context)
	// RegisterRoutes registra las rutas en el router.
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type meliHandler struct {
	useCase usecases.IMeliUseCase
	logger  log.ILogger
}

// New crea el handler HTTP de MercadoLibre.
func New(useCase usecases.IMeliUseCase, logger log.ILogger) IHandler {
	return &meliHandler{
		useCase: useCase,
		logger:  logger.WithModule("meli"),
	}
}

// RegisterRoutes registra las rutas de MercadoLibre en el router.
func (h *meliHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	meli := router.Group("/meli")
	{
		// IPN — notificaciones de MercadoLibre
		meli.POST("/notifications", h.HandleNotification)
	}
}
