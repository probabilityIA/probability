package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)



func (r *Repository) CreateCreditNote(ctx context.Context, note *entities.CreditNote) error {
	model := mappers.CreditNoteToModel(note)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create credit note: %w", err)
	}

	note.ID = model.ID
	note.InternalNumber = model.InternalNumber // Generado por BeforeCreate
	return nil
}

func (r *Repository) GetCreditNoteByID(ctx context.Context, id uint) (*entities.CreditNote, error) {
	var model models.CreditNote

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		return nil, fmt.Errorf("credit note not found: %w", err)
	}

	return mappers.CreditNoteToDomain(&model), nil
}

func (r *Repository) GetCreditNotesByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.CreditNote, error) {
	var models []*models.CreditNote

	if err := r.db.Conn(ctx).
		Where("invoice_id = ?", invoiceID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get credit notes: %w", err)
	}

	return mappers.CreditNoteListToDomain(models), nil
}

func (r *Repository) ListCreditNotes(ctx context.Context, filters map[string]interface{}) ([]*entities.CreditNote, error) {
	var modelsList []*models.CreditNote

	query := r.db.Conn(ctx).Model(&models.CreditNote{})

	// Aplicar filtros
	if businessID, ok := filters["business_id"].(uint); ok {
		query = query.Where("business_id = ?", businessID)
	}

	if status, ok := filters["status"].(string); ok {
		query = query.Where("status = ?", status)
	}

	if invoiceID, ok := filters["invoice_id"].(uint); ok {
		query = query.Where("invoice_id = ?", invoiceID)
	}

	// Ordenar por más reciente
	query = query.Order("created_at DESC")

	// Aplicar limit si existe
	if limit, ok := filters["limit"].(int); ok {
		query = query.Limit(limit)
	}

	if err := query.Find(&modelsList).Error; err != nil {
		return nil, fmt.Errorf("failed to list credit notes: %w", err)
	}

	return mappers.CreditNoteListToDomain(modelsList), nil
}

func (r *Repository) UpdateCreditNote(ctx context.Context, note *entities.CreditNote) error {
	model := mappers.CreditNoteToModel(note)

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update credit note: %w", err)
	}

	return nil
}
