package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
)

func (r *Repository) SaveSyncLog(ctx context.Context, log *domain.SyncLog) error {
	row := &models.ShipmentSyncLog{
		ShipmentID:     log.ShipmentID,
		OperationType:  log.OperationType,
		Provider:       log.Provider,
		Status:         log.Status,
		RequestURL:     log.RequestURL,
		RequestMethod:  log.RequestMethod,
		RequestPayload: datatypes.JSON(log.RequestPayload),
		ResponseStatus: log.ResponseStatus,
		ResponseBody:   datatypes.JSON(log.ResponseBody),
		ErrorMessage:   log.ErrorMessage,
		ErrorCode:      log.ErrorCode,
		CorrelationID:  log.CorrelationID,
		TriggeredBy:    log.TriggeredBy,
		UserID:         log.UserID,
		StartedAt:      log.StartedAt,
		CompletedAt:    log.CompletedAt,
		Duration:       log.Duration,
	}

	if row.StartedAt.IsZero() {
		row.StartedAt = time.Now()
	}

	if err := r.db.Conn(ctx).Create(row).Error; err != nil {
		return err
	}
	log.ID = row.ID
	return nil
}
