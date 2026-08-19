package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) CreateAuditLog(ctx context.Context, log *entities.SubscriptionAuditLog) error {
	logDB := &models.SubscriptionAuditLog{
		BusinessID:  log.BusinessID,
		ActorUserID: log.ActorUserID,
		ActorLabel:  log.ActorLabel,
		Action:      log.Action,
		Description: log.Description,
	}
	if err := r.db.Conn(ctx).Create(logDB).Error; err != nil {
		return err
	}
	log.ID = logDB.ID
	log.CreatedAt = logDB.CreatedAt
	return nil
}

func (r *Repository) ListAuditLogsByBusiness(ctx context.Context, businessID uint, limit int) ([]entities.SubscriptionAuditLog, error) {
	var logsDB []models.SubscriptionAuditLog
	err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Order("created_at desc").
		Limit(limit).
		Find(&logsDB).Error
	if err != nil {
		return nil, err
	}

	logs := make([]entities.SubscriptionAuditLog, len(logsDB))
	for i, l := range logsDB {
		logs[i] = *auditLogToEntity(&l)
	}
	return logs, nil
}
