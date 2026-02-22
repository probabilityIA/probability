package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// invoicingUseCase es el use case de facturación para Factus
// Por ahora es un stub — el procesamiento real se hace en el consumer directamente
type invoicingUseCase struct {
	factusClient    ports.IFactusClient
	integrationCore core.IIntegrationCore
	log             log.ILogger
}

// New crea el use case de facturación de Factus
func New(
	factusClient ports.IFactusClient,
	integrationCore core.IIntegrationCore,
	logger log.ILogger,
) ports.IInvoiceUseCase {
	return &invoicingUseCase{
		factusClient:    factusClient,
		integrationCore: integrationCore,
		log:             logger.WithModule("factus.usecase"),
	}
}
