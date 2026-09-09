package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type handler struct {
	uc  app.IUseCase
	log log.ILogger
}

func New(useCase app.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		uc:  useCase,
		log: logger.WithModule("push-handler"),
	}
}
