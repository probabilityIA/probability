package handlerintegrationtype

import (
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrationtype"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IntegrationTypeHandler maneja las peticiones HTTP relacionadas con tipos de integración
type IntegrationTypeHandler struct {
	usecase usecaseintegrationtype.IIntegrationTypeUseCase
	logger  log.ILogger
}

// NewIntegrationTypeHandler crea una nueva instancia del handler de tipos de integración
func New(usecase usecaseintegrationtype.IIntegrationTypeUseCase, logger log.ILogger) *IntegrationTypeHandler {
	contextualLogger := logger.WithModule("integration-types")
	return &IntegrationTypeHandler{
		usecase: usecase,
		logger:  contextualLogger,
	}
}
