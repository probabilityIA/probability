package handlerintegrations

import (
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrations"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IntegrationHandler struct {
	usecase usecaseintegrations.IIntegrationUseCase
	logger  log.ILogger
}

// New crea una nueva instancia del handler de integraciones
func New(usecase usecaseintegrations.IIntegrationUseCase, logger log.ILogger) *IntegrationHandler {
	contextualLogger := logger.WithModule("integrations")
	return &IntegrationHandler{
		usecase: usecase,
		logger:  contextualLogger,
	}
}
