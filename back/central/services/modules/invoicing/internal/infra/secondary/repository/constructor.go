package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Repository agrupa todos los repositorios del módulo de facturación
type Repository struct {
	db  db.IDatabase
	log log.ILogger
}

// Repositories contiene todos los repositorios del módulo
type Repositories struct {
	Invoice         ports.IInvoiceRepository
	InvoiceItem     ports.IInvoiceItemRepository
	Provider        ports.IInvoicingProviderRepository
	ProviderType    ports.IInvoicingProviderTypeRepository
	Config          ports.IInvoicingConfigRepository
	SyncLog         ports.IInvoiceSyncLogRepository
	CreditNote      ports.ICreditNoteRepository
}

// New crea una nueva instancia de todos los repositorios
func New(database db.IDatabase, logger log.ILogger) *Repositories {
	baseRepo := &Repository{
		db:  database,
		log: logger.WithModule("invoicing.repository"),
	}

	return &Repositories{
		Invoice:      NewInvoiceRepository(baseRepo),
		InvoiceItem:  NewInvoiceItemRepository(baseRepo),
		Provider:     NewInvoicingProviderRepository(baseRepo),
		ProviderType: NewInvoicingProviderTypeRepository(baseRepo),
		Config:       NewInvoicingConfigRepository(baseRepo),
		SyncLog:      NewInvoiceSyncLogRepository(baseRepo),
		CreditNote:   NewCreditNoteRepository(baseRepo),
	}
}
