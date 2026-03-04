package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// IHandler define la interfaz del handler HTTP del mock Shopify API.
type IHandler interface {
	RegisterRoutes(router *gin.Engine)
}

// Handler maneja las peticiones HTTP del mock de Shopify API.
type Handler struct {
	mockAPI *usecases.MockAPIServer
	logger  log.ILogger
}

// New crea una nueva instancia del handler.
func New(mockAPI *usecases.MockAPIServer, logger log.ILogger) IHandler {
	return &Handler{
		mockAPI: mockAPI,
		logger:  logger,
	}
}
