package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// useCase implementa todos los casos de uso del módulo de facturación
type useCase struct {
	// Repositorios
	invoiceRepo     ports.IInvoiceRepository
	invoiceItemRepo ports.IInvoiceItemRepository
	configRepo      ports.IInvoicingConfigRepository
	syncLogRepo     ports.IInvoiceSyncLogRepository
	creditNoteRepo  ports.ICreditNoteRepository
	orderRepo       ports.IOrderRepository

	// Servicios externos
	integrationCore core.IIntegrationCore // Reemplaza providerRepo, providerTypeRepo y providerClient
	encryption      ports.IEncryptionService
	eventPublisher  ports.IEventPublisher

	// Logger
	log log.ILogger
}

// New crea una nueva instancia del use case de facturación
func New(
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
		log:             logger.WithModule("invoicing.usecase"),
	}
}
