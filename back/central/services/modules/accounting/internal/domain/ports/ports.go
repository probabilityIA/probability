package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
)

type IEmailSender interface {
	SendHTML(ctx context.Context, to, subject, html string) error
}

type IDianEmitter interface {
	Emit(ctx context.Context, integrationID uint, inv *entities.Invoice, config map[string]interface{}) (*entities.DianEmitResult, error)
}

type IDianIntegration interface {
	GetGlobalIntegrationID(ctx context.Context) (uint, error)
	GetConfig(ctx context.Context, integrationID uint) (map[string]interface{}, error)
	Ensure(ctx context.Context, name string, config, credentials map[string]interface{}) (uint, error)
	TestConnection(ctx context.Context, config, credentials map[string]interface{}) error
}

type IRepository interface {
	ListConcepts(ctx context.Context) ([]entities.Concept, error)
	GetConceptByID(ctx context.Context, id uint) (*entities.Concept, error)
	ConceptCodeExists(ctx context.Context, code string, excludeID *uint) (bool, error)
	CreateConcept(ctx context.Context, dto dtos.CreateConceptDTO) (*entities.Concept, error)
	UpdateConcept(ctx context.Context, dto dtos.UpdateConceptDTO) (*entities.Concept, error)

	ListTaxes(ctx context.Context) ([]entities.Tax, error)
	GetTaxByID(ctx context.Context, id uint) (*entities.Tax, error)
	TaxCodeExists(ctx context.Context, code string, excludeID *uint) (bool, error)
	CreateTax(ctx context.Context, dto dtos.CreateTaxDTO) (*entities.Tax, error)
	UpdateTax(ctx context.Context, dto dtos.UpdateTaxDTO) (*entities.Tax, error)

	SetConceptTax(ctx context.Context, dto dtos.SetConceptTaxDTO) error

	ListEntries(ctx context.Context, params dtos.ListEntriesParams) ([]entities.Entry, int64, error)
	GetEntryByID(ctx context.Context, id uint) (*entities.Entry, error)
	CreateEntry(ctx context.Context, e *entities.Entry) (bool, error)
	DeleteEntry(ctx context.Context, id uint) error

	Report(ctx context.Context, params dtos.ReportParams) ([]dtos.ReportConceptRow, []dtos.ReportTaxRow, error)

	FindSyncCandidates(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error)

	CreateInvoice(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error)
	UpdateInvoice(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error)
	GetInvoiceByID(ctx context.Context, id uint) (*entities.Invoice, error)
	ListInvoices(ctx context.Context, params dtos.ListInvoicesParams) ([]entities.Invoice, int64, error)
	DeleteInvoice(ctx context.Context, id uint) error
	SetInvoiceStatus(ctx context.Context, id uint, status string, sentAt, paidAt *time.Time, emailTo string) error
	SetInvoiceCustomer(ctx context.Context, id uint, document, phone, address string) error
	SetInvoiceDian(ctx context.Context, id uint, result *entities.DianEmitResult) error
	DeleteEntryBySource(ctx context.Context, sourceType, sourceID string) error
	GetTaxesByIDs(ctx context.Context, ids []uint) ([]entities.Tax, error)
	GetBusinessName(ctx context.Context, id uint) (string, error)
	GetBusinessDefaultEmail(ctx context.Context, id uint) (string, error)
}
