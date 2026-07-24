package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	ListConcepts(ctx context.Context) ([]entities.Concept, error)
	CreateConcept(ctx context.Context, dto dtos.CreateConceptDTO) (*entities.Concept, error)
	UpdateConcept(ctx context.Context, dto dtos.UpdateConceptDTO) (*entities.Concept, error)

	ListTaxes(ctx context.Context) ([]entities.Tax, error)
	CreateTax(ctx context.Context, dto dtos.CreateTaxDTO) (*entities.Tax, error)
	UpdateTax(ctx context.Context, dto dtos.UpdateTaxDTO) (*entities.Tax, error)

	SetConceptTax(ctx context.Context, dto dtos.SetConceptTaxDTO) error

	ListEntries(ctx context.Context, params dtos.ListEntriesParams) ([]entities.Entry, int64, error)
	CreateManualEntry(ctx context.Context, dto dtos.CreateEntryDTO) (*entities.Entry, error)
	DeleteManualEntry(ctx context.Context, id uint) error

	Report(ctx context.Context, params dtos.ReportParams) (*dtos.ReportResponse, error)
	Sync(ctx context.Context) (*dtos.SyncResult, error)

	CreateInvoice(ctx context.Context, dto dtos.CreateInvoiceDTO) (*entities.Invoice, error)
	UpdateInvoice(ctx context.Context, dto dtos.UpdateInvoiceDTO) (*entities.Invoice, error)
	GetInvoice(ctx context.Context, id uint) (*entities.Invoice, error)
	ListInvoices(ctx context.Context, params dtos.ListInvoicesParams) ([]entities.Invoice, int64, error)
	DeleteInvoice(ctx context.Context, id uint) error
	SendInvoice(ctx context.Context, dto dtos.SendInvoiceDTO) (*entities.Invoice, error)
	MarkInvoicePaid(ctx context.Context, id uint) (*entities.Invoice, error)
	CancelInvoice(ctx context.Context, id uint) (*entities.Invoice, error)
}

type UseCase struct {
	repo  ports.IRepository
	email ports.IEmailSender
	log   log.ILogger
}

func New(repo ports.IRepository, emailSender ports.IEmailSender, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, email: emailSender, log: logger.WithModule("accounting")}
}
