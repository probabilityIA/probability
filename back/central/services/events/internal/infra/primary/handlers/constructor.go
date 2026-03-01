package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// ISSEHandler define la interfaz del handler SSE
type ISSEHandler interface {
	HandleSSE(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
	GetManager() ports.ISSEPublisher
}

// SSEHandler maneja las conexiones Server-Sent Events
type SSEHandler struct {
	eventManager ports.ISSEPublisher
	logger       log.ILogger
}

// New crea un nuevo handler SSE
func New(eventManager ports.ISSEPublisher, logger log.ILogger) ISSEHandler {
	return &SSEHandler{
		eventManager: eventManager,
		logger:       logger,
	}
}

// GetManager retorna el manager interno
func (h *SSEHandler) GetManager() ports.ISSEPublisher {
	return h.eventManager
}
