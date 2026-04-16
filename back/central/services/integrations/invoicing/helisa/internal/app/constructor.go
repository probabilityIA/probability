package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// invoicingUseCase es el use case de facturación para Helisa
type invoicingUseCase struct {
	helisaClient    ports.IHelisaClient
	integrationCore core.IIntegrationService
	log             log.ILogger
}

// New crea el use case de facturación de Helisa
func New(
	helisaClient ports.IHelisaClient,
	integrationCore core.IIntegrationService,
	logger log.ILogger,
) ports.IInvoiceUseCase {
	return &invoicingUseCase{
		helisaClient:    helisaClient,
		integrationCore: integrationCore,
		log:             logger.WithModule("helisa.usecase"),
	}
}
