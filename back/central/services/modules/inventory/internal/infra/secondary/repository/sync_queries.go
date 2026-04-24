package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateSyncLog(ctx context.Context, log *entities.InventorySyncLog) (*entities.InventorySyncLog, error) {
	m := &models.InventorySyncLog{
		BusinessID:    log.BusinessID,
		IntegrationID: log.IntegrationID,
		Direction:     log.Direction,
		DirectionKey:  log.Direction,
		PayloadHash:   log.PayloadHash,
		Status:        log.Status,
		Error:         log.Error,
		SyncedAt:      log.SyncedAt,
	}
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.SyncLogModelToEntity(m), nil
}

func (r *Repository) GetSyncLogByHash(ctx context.Context, businessID uint, direction, hash string) (*entities.InventorySyncLog, error) {
	var m models.InventorySyncLog
	err := r.db.Conn(ctx).
		Where("business_id = ? AND direction_key = ? AND payload_hash = ?", businessID, direction, hash).
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mappers.SyncLogModelToEntity(&m), nil
}

func (r *Repository) UpdateSyncLogStatus(ctx context.Context, id uint, status, errorMsg string) error {
	return r.db.Conn(ctx).Model(&models.InventorySyncLog{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status": status,
			"error":  errorMsg,
		}).Error
}

func (r *Repository) ListSyncLogs(ctx context.Context, params dtos.ListSyncLogsParams) ([]entities.InventorySyncLog, int64, error) {
	var ml []models.InventorySyncLog
	var total int64
	q := r.db.Conn(ctx).Model(&models.InventorySyncLog{}).Where("business_id = ?", params.BusinessID)
	if params.IntegrationID != nil {
		q = q.Where("integration_id = ?", *params.IntegrationID)
	}
	if params.Direction != "" {
		q = q.Where("direction_key = ?", params.Direction)
	}
	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(params.Offset()).Limit(params.PageSize).Order("id DESC").Find(&ml).Error; err != nil {
		return nil, 0, err
	}
	out := make([]entities.InventorySyncLog, len(ml))
	for i := range ml {
		out[i] = *mappers.SyncLogModelToEntity(&ml[i])
	}
	return out, total, nil
}
