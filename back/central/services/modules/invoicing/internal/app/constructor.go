package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// useCase implementa todos los casos de uso del módulo de facturación
type useCase struct {
<<<<<<< HEAD
	// Repositorios
	invoiceRepo     ports.IInvoiceRepository
	invoiceItemRepo ports.IInvoiceItemRepository
	configRepo      ports.IInvoicingConfigRepository
	syncLogRepo     ports.IInvoiceSyncLogRepository
	creditNoteRepo  ports.ICreditNoteRepository
	orderRepo       ports.IOrderRepository
=======
	// Repositorio único (implementa TODAS las operaciones de persistencia)
	repo ports.IRepository
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

	// Servicios externos
	integrationCore core.IIntegrationCore // Reemplaza providerRepo, providerTypeRepo y providerClient
	encryption      ports.IEncryptionService
	eventPublisher  ports.IEventPublisher
<<<<<<< HEAD
=======
	ssePublisher    ports.IInvoiceSSEPublisher
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

	// Logger
	log log.ILogger
}

// New crea una nueva instancia del use case de facturación
func New(
<<<<<<< HEAD
	invoiceRepo ports.IInvoiceRepository,
	invoiceItemRepo ports.IInvoiceItemRepository,
	configRepo ports.IInvoicingConfigRepository,
	syncLogRepo ports.IInvoiceSyncLogRepository,
	creditNoteRepo ports.ICreditNoteRepository,
	orderRepo ports.IOrderRepository,
	integrationCore core.IIntegrationCore,
	encryption ports.IEncryptionService,
	eventPublisher ports.IEventPublisher,
	logger log.ILogger,
) ports.IUseCase {
	return &useCase{
		invoiceRepo:     invoiceRepo,
		invoiceItemRepo: invoiceItemRepo,
		configRepo:      configRepo,
		syncLogRepo:     syncLogRepo,
		creditNoteRepo:  creditNoteRepo,
		orderRepo:       orderRepo,
		integrationCore: integrationCore,
		encryption:      encryption,
		eventPublisher:  eventPublisher,
=======
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
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
		log:             logger.WithModule("invoicing.usecase"),
	}
}
