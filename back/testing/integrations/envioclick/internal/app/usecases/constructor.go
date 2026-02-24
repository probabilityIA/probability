package usecases

import (
	"github.com/secamc93/probability/back/testing/integrations/envioclick/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// APISimulator simula el API de EnvioClick
type APISimulator struct {
	logger     log.ILogger
	Repository *domain.ShipmentRepository // Exportado para acceso desde bundle
}

// NewAPISimulator crea una nueva instancia del simulador de API
func NewAPISimulator(logger log.ILogger) *APISimulator {
	return &APISimulator{
		logger:     logger,
		Repository: domain.NewShipmentRepository(),
	}
}
