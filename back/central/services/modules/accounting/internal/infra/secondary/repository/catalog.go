package repository

import (
	"context"
	stderrors "errors"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) taxIndex(ctx context.Context) (map[uint]models.AccountingTax, error) {
	var taxes []models.AccountingTax
	if err := r.db.Conn(ctx).Find(&taxes).Error; err != nil {
		return nil, err
	}
	index := make(map[uint]models.AccountingTax, len(taxes))
	for _, t := range taxes {
		index[t.ID] = t
	}
	return index, nil
}

func (r *Repository) ListConcepts(ctx context.Context) ([]entities.Concept, error) {
	var conceptModels []models.AccountingConcept
	if err := r.db.Conn(ctx).Preload("Taxes").Order("kind ASC, code ASC").Find(&conceptModels).Error; err != nil {
		return nil, err
	}
	index, err := r.taxIndex(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entities.Concept, len(conceptModels))
	for i := range conceptModels {
		out[i] = *mappers.ConceptToEntity(&conceptModels[i], index)
	}
	return out, nil
}

func (r *Repository) GetConceptByID(ctx context.Context, id uint) (*entities.Concept, error) {
	var model models.AccountingConcept
	err := r.db.Conn(ctx).Preload("Taxes").First(&model, id).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrConceptNotFound
		}
		return nil, err
	}
	index, err := r.taxIndex(ctx)
	if err != nil {
		return nil, err
	}
	return mappers.ConceptToEntity(&model, index), nil
}

func (r *Repository) ConceptCodeExists(ctx context.Context, code string, excludeID *uint) (bool, error) {
	var count int64
	q := r.db.Conn(ctx).Model(&models.AccountingConcept{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateConcept(ctx context.Context, dto dtos.CreateConceptDTO) (*entities.Concept, error) {
	model := &models.AccountingConcept{
		Code:         dto.Code,
		Name:         dto.Name,
		Description:  dto.Description,
		Kind:         dto.Kind,
		IsRealIncome: dto.IsRealIncome,
		IsAutomatic:  false,
		SourceType:   entities.SourceManual,
		IsActive:     true,
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return r.GetConceptByID(ctx, model.ID)
}

func (r *Repository) UpdateConcept(ctx context.Context, dto dtos.UpdateConceptDTO) (*entities.Concept, error) {
	updates := map[string]interface{}{
		"name":           dto.Name,
		"description":    dto.Description,
		"is_real_income": dto.IsRealIncome,
		"is_active":      dto.IsActive,
	}
	res := r.db.Conn(ctx).Model(&models.AccountingConcept{}).Where("id = ?", dto.ID).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrConceptNotFound
	}
	return r.GetConceptByID(ctx, dto.ID)
}

func (r *Repository) ListTaxes(ctx context.Context) ([]entities.Tax, error) {
	var taxModels []models.AccountingTax
	if err := r.db.Conn(ctx).Order("code ASC").Find(&taxModels).Error; err != nil {
		return nil, err
	}
	out := make([]entities.Tax, len(taxModels))
	for i := range taxModels {
		out[i] = *mappers.TaxToEntity(&taxModels[i])
	}
	return out, nil
}

func (r *Repository) GetTaxByID(ctx context.Context, id uint) (*entities.Tax, error) {
	var model models.AccountingTax
	err := r.db.Conn(ctx).First(&model, id).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrTaxNotFound
		}
		return nil, err
	}
	return mappers.TaxToEntity(&model), nil
}

func (r *Repository) TaxCodeExists(ctx context.Context, code string, excludeID *uint) (bool, error) {
	var count int64
	q := r.db.Conn(ctx).Model(&models.AccountingTax{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateTax(ctx context.Context, dto dtos.CreateTaxDTO) (*entities.Tax, error) {
	kind := dto.Kind
	if kind == "" {
		kind = entities.TaxKindCharge
	}
	model := &models.AccountingTax{
		Code:        dto.Code,
		Name:        dto.Name,
		Description: dto.Description,
		RatePercent: dto.RatePercent,
		Kind:        kind,
		IsActive:    true,
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return mappers.TaxToEntity(model), nil
}

func (r *Repository) UpdateTax(ctx context.Context, dto dtos.UpdateTaxDTO) (*entities.Tax, error) {
	updates := map[string]interface{}{
		"name":         dto.Name,
		"description":  dto.Description,
		"rate_percent": dto.RatePercent,
		"is_active":    dto.IsActive,
	}
	if dto.Kind != "" {
		updates["kind"] = dto.Kind
	}
	res := r.db.Conn(ctx).Model(&models.AccountingTax{}).Where("id = ?", dto.ID).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrTaxNotFound
	}
	return r.GetTaxByID(ctx, dto.ID)
}

func (r *Repository) SetConceptTax(ctx context.Context, dto dtos.SetConceptTaxDTO) error {
	var existing models.AccountingConceptTax
	err := r.db.Conn(ctx).
		Where("concept_id = ? AND tax_id = ?", dto.ConceptID, dto.TaxID).
		First(&existing).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Conn(ctx).Create(&models.AccountingConceptTax{
				ConceptID: dto.ConceptID,
				TaxID:     dto.TaxID,
				IsActive:  dto.IsActive,
			}).Error
		}
		return err
	}
	return r.db.Conn(ctx).Model(&models.AccountingConceptTax{}).
		Where("id = ?", existing.ID).
		Update("is_active", dto.IsActive).Error
}
