package usecases

import (
	"github.com/secamc93/probability/back/integrationTest/integrations/softpymes/internal/domain"
	"github.com/secamc93/probability/back/integrationTest/shared/log"
)

// NewAPISimulator crea una nueva instancia del simulador de API
func NewAPISimulator(logger log.ILogger) domain.IAPIClient {
	return &APISimulator{
		logger:     logger,
		Repository: domain.NewInvoiceRepository(),
	}
}
