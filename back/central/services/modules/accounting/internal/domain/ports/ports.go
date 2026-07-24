package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
)

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
}
