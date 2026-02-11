package app

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// useCase implementa todos los casos de uso del módulo de facturación
type useCase struct {
	// Repositorio único (implementa TODAS las operaciones de persistencia)
	repo ports.IRepository

	// Servicios externos
	encryption            ports.IEncryptionService
	eventPublisher        ports.IEventPublisher
	ssePublisher          ports.IInvoiceSSEPublisher
	invoiceRequestPub     ports.IInvoiceRequestPublisher // Publisher para requests a proveedores

	// Logger
	log log.ILogger
}

// New crea una nueva instancia del use case de facturación
func New(
	repo ports.IRepository,
	encryption ports.IEncryptionService,
	eventPublisher ports.IEventPublisher,
	ssePublisher ports.IInvoiceSSEPublisher,
	invoiceRequestPub ports.IInvoiceRequestPublisher,
	logger log.ILogger,
) ports.IUseCase {
	return &useCase{
		repo:              repo,
		encryption:        encryption,
		eventPublisher:    eventPublisher,
		ssePublisher:      ssePublisher,
		invoiceRequestPub: invoiceRequestPub,
		log:               logger.WithModule("invoicing.usecase"),
	}
}
