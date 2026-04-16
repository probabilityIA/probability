package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/amazon/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define los endpoints HTTP de Amazon.
type IHandler interface {
	// HandleNotification recibe notificaciones de Amazon SQS/SNS.
	HandleNotification(c *gin.Context)
	// RegisterRoutes registra las rutas en el router.
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type amazonHandler struct {
	useCase usecases.IAmazonUseCase
	logger  log.ILogger
}

// New crea el handler HTTP de Amazon.
func New(useCase usecases.IAmazonUseCase, logger log.ILogger) IHandler {
	return &amazonHandler{
		useCase: useCase,
		logger:  logger.WithModule("amazon"),
	}
}

// RegisterRoutes registra las rutas de Amazon en el router.
func (h *amazonHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	amazon := router.Group("/amazon")
	{
		amazon.POST("/notification", h.HandleNotification)
	}
}
