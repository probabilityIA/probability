package handlers

import (
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IntegrationHandler struct {
	usecase domain.IIntegrationUseCase
	logger  log.ILogger
}

// New crea una nueva instancia del handler de integraciones
func New(usecase domain.IIntegrationUseCase, logger log.ILogger) IIntegrationHandler {
	contextualLogger := logger.WithModule("integrations")
	return &IntegrationHandler{
		usecase: usecase,
		logger:  contextualLogger,
	}
}
