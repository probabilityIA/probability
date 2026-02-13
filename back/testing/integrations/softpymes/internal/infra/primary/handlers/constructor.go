package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/softpymes/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// IHandler define la interfaz del handler HTTP
type IHandler interface {
	RegisterRoutes(router *gin.Engine)
}

// Handler maneja las peticiones HTTP del simulador de Softpymes
type Handler struct {
	apiSimulator *usecases.APISimulator
	logger       log.ILogger
}

// New crea una nueva instancia del handler
func New(apiSimulator *usecases.APISimulator, logger log.ILogger) IHandler {
	return &Handler{
		apiSimulator: apiSimulator,
		logger:       logger,
	}
}
