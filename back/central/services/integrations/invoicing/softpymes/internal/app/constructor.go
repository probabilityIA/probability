package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// invoicingUseCase es el use case de facturación para Softpymes
type invoicingUseCase struct {
	client ports.ISoftpymesClient
	log    log.ILogger
}

// New crea el use case de facturación de Softpymes
func New(
	client ports.ISoftpymesClient,
	logger log.ILogger,
) ports.IInvoiceUseCase {
	return &invoicingUseCase{
		client: client,
		log:    logger.WithModule("softpymes.usecase"),
	}
}
