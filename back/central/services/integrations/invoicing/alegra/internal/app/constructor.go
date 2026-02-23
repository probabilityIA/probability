package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/alegra/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// invoicingUseCase es el use case de facturación para Alegra
type invoicingUseCase struct {
	alegraClient    ports.IAlegraClient
	integrationCore core.IIntegrationService
	log             log.ILogger
}

// New crea el use case de facturación de Alegra
func New(
	alegraClient ports.IAlegraClient,
	integrationCore core.IIntegrationService,
	logger log.ILogger,
) ports.IInvoiceUseCase {
	return &invoicingUseCase{
		alegraClient:    alegraClient,
		integrationCore: integrationCore,
		log:             logger.WithModule("alegra.usecase"),
	}
}
