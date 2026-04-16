package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// invoicingUseCase es el use case de facturación para World Office
type invoicingUseCase struct {
	worldOfficeClient ports.IWorldOfficeClient
	integrationCore   core.IIntegrationService
	log               log.ILogger
}

// New crea el use case de facturación de World Office
func New(
	worldOfficeClient ports.IWorldOfficeClient,
	integrationCore core.IIntegrationService,
	logger log.ILogger,
) ports.IInvoiceUseCase {
	return &invoicingUseCase{
		worldOfficeClient: worldOfficeClient,
		integrationCore:   integrationCore,
		log:               logger.WithModule("world_office.usecase"),
	}
}
