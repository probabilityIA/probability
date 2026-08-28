package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

const scopeCodePlatform = "platform"

func (r *Repository) GetActiveDocuments(ctx context.Context) ([]entities.LegalDocument, error) {
	var docs []models.LegalDocument
	err := r.db.Conn(ctx).
		Where("is_active = ? AND deleted_at IS NULL AND effective_from <= ?", true, time.Now()).
		Order("code ASC").
		Find(&docs).Error
	if err != nil {
		return nil, err
	}
	return mappers.ToDomainDocuments(docs), nil
}

func (r *Repository) GetDocumentsByIDs(ctx context.Context, ids []uint) ([]entities.LegalDocument, error) {
	if len(ids) == 0 {
		return []entities.LegalDocument{}, nil
	}
	var docs []models.LegalDocument
	err := r.db.Conn(ctx).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&docs).Error
	if err != nil {
		return nil, err
	}
	return mappers.ToDomainDocuments(docs), nil
}

func (r *Repository) GetAcceptedDocumentIDs(ctx context.Context, userID uint) (map[uint]bool, error) {
	var ids []uint
	err := r.db.Conn(ctx).
		Model(&models.LegalAcceptance{}).
		Where("user_id = ?", userID).
		Pluck("legal_document_id", &ids).Error
	if err != nil {
		return nil, err
	}

	aceptados := make(map[uint]bool, len(ids))
	for _, id := range ids {
		aceptados[id] = true
	}
	return aceptados, nil
}

func (r *Repository) SaveAcceptances(ctx context.Context, aceptaciones []entities.LegalAcceptance) error {
	if len(aceptaciones) == 0 {
		return nil
	}
	filas := mappers.ToModelAcceptances(aceptaciones)
	return r.db.Conn(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "legal_document_id"}},
			DoNothing: true,
		}).
		Create(&filas).Error
}

func (r *Repository) IsPlatformUser(ctx context.Context, userID uint) (bool, error) {
	var usuario models.User
	err := r.db.Conn(ctx).
		Preload("Scope").
		Where("id = ?", userID).
		First(&usuario).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if usuario.Scope == nil {
		return false, nil
	}
	return usuario.Scope.Code == scopeCodePlatform, nil
}
