package repository

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type entryRow struct {
	models.AccountingEntry
	ConceptCode  string
	ConceptName  string
	BusinessName string
}

func (r *Repository) entriesQuery(ctx context.Context, params dtos.ListEntriesParams) *gorm.DB {
	q := r.db.Conn(ctx).Model(&models.AccountingEntry{}).
		Where("accounting_entries.deleted_at IS NULL")
	if params.From != nil {
		q = q.Where("accounting_entries.entry_date >= ?", params.From.Format("2006-01-02"))
	}
	if params.To != nil {
		q = q.Where("accounting_entries.entry_date <= ?", params.To.Format("2006-01-02"))
	}
	if params.ConceptID != nil {
		q = q.Where("accounting_entries.concept_id = ?", *params.ConceptID)
	}
	if params.Kind != "" {
		q = q.Where("accounting_entries.kind = ?", params.Kind)
	}
	return q
}

func (r *Repository) ListEntries(ctx context.Context, params dtos.ListEntriesParams) ([]entities.Entry, int64, error) {
	var total int64
	if err := r.entriesQuery(ctx, params).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []entryRow
	err := r.entriesQuery(ctx, params).
		Select("accounting_entries.*, c.code AS concept_code, c.name AS concept_name, COALESCE(b.name, '') AS business_name").
		Joins("JOIN accounting_concepts c ON c.id = accounting_entries.concept_id").
		Joins("LEFT JOIN business b ON b.id = accounting_entries.business_id").
		Order("accounting_entries.entry_date DESC, accounting_entries.id DESC").
		Offset(params.Offset()).Limit(params.PageSize).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]entities.Entry, len(rows))
	for i := range rows {
		out[i] = *mappers.EntryToEntity(&rows[i].AccountingEntry, rows[i].ConceptCode, rows[i].ConceptName, rows[i].BusinessName)
	}
	return out, total, nil
}

func (r *Repository) GetEntryByID(ctx context.Context, id uint) (*entities.Entry, error) {
	var row entryRow
	err := r.db.Conn(ctx).Model(&models.AccountingEntry{}).
		Select("accounting_entries.*, c.code AS concept_code, c.name AS concept_name, COALESCE(b.name, '') AS business_name").
		Joins("JOIN accounting_concepts c ON c.id = accounting_entries.concept_id").
		Joins("LEFT JOIN business b ON b.id = accounting_entries.business_id").
		Where("accounting_entries.id = ? AND accounting_entries.deleted_at IS NULL", id).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, domainerrors.ErrEntryNotFound
	}
	return mappers.EntryToEntity(&row.AccountingEntry, row.ConceptCode, row.ConceptName, row.BusinessName), nil
}

func (r *Repository) CreateEntry(ctx context.Context, e *entities.Entry) (bool, error) {
	model := &models.AccountingEntry{
		ConceptID:   e.ConceptID,
		BusinessID:  e.BusinessID,
		EntryDate:   e.EntryDate,
		Amount:      e.Amount,
		Kind:        e.Kind,
		SourceType:  e.SourceType,
		SourceID:    e.SourceID,
		Description: e.Description,
		IsAutomatic: e.IsAutomatic,
		TaxTotal:    e.TaxTotal,
		TaxDetail:   datatypes.JSON(mappers.TaxLinesToJSON(e.TaxDetail)),
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		if isUniqueViolation(err) {
			return false, nil
		}
		return false, err
	}
	e.ID = model.ID
	return true, nil
}

func (r *Repository) DeleteEntry(ctx context.Context, id uint) error {
	res := r.db.Conn(ctx).Delete(&models.AccountingEntry{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrEntryNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	if stderrors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	return strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "SQLSTATE 23505")
}
