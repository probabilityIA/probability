package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) Save(ctx context.Context, log *domain.WebhookLog) error {
	dbLog := mappers.ToDBWebhookLog(log)
	if err := r.db.Conn(ctx).Create(dbLog).Error; err != nil {
		return err
	}
	log.ID = dbLog.ID
	log.CreatedAt = dbLog.CreatedAt
	return nil
}

func (r *Repository) MarkProcessed(ctx context.Context, id uuid.UUID, errorMessage *string) error {
	now := time.Now()
	updates := map[string]any{
		"status":        domain.WebhookLogStatusProcessed,
		"processed_at":  &now,
		"error_message": errorMessage,
	}
	if errorMessage != nil {
		updates["status"] = domain.WebhookLogStatusFailed
	}
	return r.db.Conn(ctx).Model(&models.WebhookLog{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) TrimOldBySource(ctx context.Context, source string, keepCount int) error {
	subQuery := r.db.Conn(ctx).
		Model(&models.WebhookLog{}).
		Select("id").
		Where("source = ?", source).
		Order("created_at DESC").
		Offset(keepCount)

	return r.db.Conn(ctx).
		Where("id IN (?)", subQuery).
		Delete(&models.WebhookLog{}).Error
}
