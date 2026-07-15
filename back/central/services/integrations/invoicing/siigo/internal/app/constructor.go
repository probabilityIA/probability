package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type invoicingUseCase struct {
	siigoClient     ports.ISiigoClient
	integrationCore core.IIntegrationService
	rabbit          rabbitmq.IQueue
	productRepo     ports.IProductReadRepository
	log             log.ILogger
}

func New(
	siigoClient ports.ISiigoClient,
	integrationCore core.IIntegrationService,
	rabbit rabbitmq.IQueue,
	productRepo ports.IProductReadRepository,
	logger log.ILogger,
) ports.IInvoiceUseCase {
	return &invoicingUseCase{
		siigoClient:     siigoClient,
		integrationCore: integrationCore,
		rabbit:          rabbit,
		productRepo:     productRepo,
		log:             logger.WithModule("siigo.usecase"),
	}
}
