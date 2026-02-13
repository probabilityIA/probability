package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// useCase implementa todos los casos de uso del módulo de facturación
type useCase struct {
	// Repositorio único (implementa TODAS las operaciones de persistencia)
	repo ports.IRepository

	// Servicios externos
	integrationCore core.IIntegrationCore // Reemplaza providerRepo, providerTypeRepo y providerClient
	encryption      ports.IEncryptionService
	eventPublisher  ports.IEventPublisher
	ssePublisher    ports.IInvoiceSSEPublisher

	// Logger
	log log.ILogger
}

// New crea una nueva instancia del use case de facturación
func New(
	repo ports.IRepository,
	integrationCore core.IIntegrationCore,
	encryption ports.IEncryptionService,
	eventPublisher ports.IEventPublisher,
	ssePublisher ports.IInvoiceSSEPublisher,
	logger log.ILogger,
) ports.IUseCase {
	return &useCase{
		repo:            repo,
		integrationCore: integrationCore,
		encryption:      encryption,
		eventPublisher:  eventPublisher,
		ssePublisher:    ssePublisher,
		log:             logger.WithModule("invoicing.usecase"),
	}
}
