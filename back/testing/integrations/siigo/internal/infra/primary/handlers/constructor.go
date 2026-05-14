package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/siigo/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type IHandler interface {
	RegisterRoutes(router *gin.Engine)
}

type Handler struct {
	apiSimulator *usecases.APISimulator
	logger       log.ILogger
}

func New(apiSimulator *usecases.APISimulator, logger log.ILogger) IHandler {
	return &Handler{
		apiSimulator: apiSimulator,
		logger:       logger,
	}
}
