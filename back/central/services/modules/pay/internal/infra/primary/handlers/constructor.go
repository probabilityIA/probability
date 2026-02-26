package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define la interfaz del handler de pagos
type IHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	CreatePayment(c *gin.Context)
	GetPayment(c *gin.Context)
	ListPayments(c *gin.Context)
}

// handler implementa IHandler
type handler struct {
	useCase ports.IUseCase
	log     log.ILogger
}

// New crea una nueva instancia del handler de pagos
func New(useCase ports.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		useCase: useCase,
		log:     logger.WithModule("pay.handler"),
	}
}
