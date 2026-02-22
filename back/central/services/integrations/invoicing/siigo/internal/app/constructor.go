package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// invoicingUseCase es el use case de facturación para Siigo
// Por ahora es un stub — el procesamiento real se hace en el consumer directamente
type invoicingUseCase struct {
	siigoClient     ports.ISiigoClient
	integrationCore core.IIntegrationCore
	log             log.ILogger
}

// New crea el use case de facturación de Siigo
func New(
	siigoClient ports.ISiigoClient,
	integrationCore core.IIntegrationCore,
	logger log.ILogger,
) ports.IInvoiceUseCase {
	return &invoicingUseCase{
		siigoClient:     siigoClient,
		integrationCore: integrationCore,
		log:             logger.WithModule("siigo.usecase"),
	}
}
