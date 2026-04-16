package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define los métodos HTTP del handler de monitoreo
type IHandler interface {
	WebhookGrafana(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type handler struct {
	useCase ports.IUseCase
	log     log.ILogger
	env     env.IConfig
}

// New crea una nueva instancia del handler de monitoreo
func New(useCase ports.IUseCase, logger log.ILogger, environment env.IConfig) IHandler {
	return &handler{
		useCase: useCase,
		log:     logger,
		env:     environment,
	}
}
