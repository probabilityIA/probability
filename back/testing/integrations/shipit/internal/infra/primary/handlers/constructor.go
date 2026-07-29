package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/shipit/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type IHandler interface {
	RegisterRoutes(router *gin.Engine)
}

type Handler struct {
	apiSimulator *usecases.APISimulator
	log          log.ILogger
}

func New(apiSimulator *usecases.APISimulator, logger log.ILogger) IHandler {
	return &Handler{
		apiSimulator: apiSimulator,
		log:          logger,
	}
}
