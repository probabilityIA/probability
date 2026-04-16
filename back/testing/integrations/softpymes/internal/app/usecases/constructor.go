package usecases

import (
	"github.com/secamc93/probability/back/testing/integrations/softpymes/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// NewAPISimulator crea una nueva instancia del simulador de API
func NewAPISimulator(logger log.ILogger) *APISimulator {
	return &APISimulator{
		logger:     logger,
		Repository: domain.NewInvoiceRepository(),
	}
}
