package app

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// useCase implementa todos los casos de uso del módulo de facturación
type useCase struct {
	// Repositorios
	invoiceRepo         ports.IInvoiceRepository
	invoiceItemRepo     ports.IInvoiceItemRepository
	providerRepo        ports.IInvoicingProviderRepository
	providerTypeRepo    ports.IInvoicingProviderTypeRepository
	configRepo          ports.IInvoicingConfigRepository
	syncLogRepo         ports.IInvoiceSyncLogRepository
	creditNoteRepo      ports.ICreditNoteRepository
	orderRepo           ports.IOrderRepository

	// Servicios externos
	providerClient  ports.IInvoicingProviderClient
	encryption      ports.IEncryptionService
	eventPublisher  ports.IEventPublisher

	// Logger
	log log.ILogger
}

// New crea una nueva instancia del use case de facturación
func New(
	invoiceRepo ports.IInvoiceRepository,
	invoiceItemRepo ports.IInvoiceItemRepository,
	providerRepo ports.IInvoicingProviderRepository,
	providerTypeRepo ports.IInvoicingProviderTypeRepository,
	configRepo ports.IInvoicingConfigRepository,
	syncLogRepo ports.IInvoiceSyncLogRepository,
	creditNoteRepo ports.ICreditNoteRepository,
	orderRepo ports.IOrderRepository,
	providerClient ports.IInvoicingProviderClient,
	encryption ports.IEncryptionService,
	eventPublisher ports.IEventPublisher,
	logger log.ILogger,
) ports.IUseCase {
	return &useCase{
		invoiceRepo:      invoiceRepo,
		invoiceItemRepo:  invoiceItemRepo,
		providerRepo:     providerRepo,
		providerTypeRepo: providerTypeRepo,
		configRepo:       configRepo,
		syncLogRepo:      syncLogRepo,
		creditNoteRepo:   creditNoteRepo,
		orderRepo:        orderRepo,
		providerClient:   providerClient,
		encryption:       encryption,
		eventPublisher:   eventPublisher,
		log:              logger.WithModule("invoicing.usecase"),
	}
}
