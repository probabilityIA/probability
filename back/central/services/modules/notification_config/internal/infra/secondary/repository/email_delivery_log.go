package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
)

const emailLogModuleNotification = "notification"

type deliveryLogRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func NewDeliveryLogRepository(database db.IDatabase, logger log.ILogger) ports.IDeliveryLogRepository {
	return &deliveryLogRepository{
		db:     database,
		logger: logger.WithModule("delivery_log_repository"),
	}
}

func (r *deliveryLogRepository) CreateEmailLog(ctx context.Context, entry *entities.EmailDeliveryLog) error {
	model := &models.EmailLog{
		BusinessID:    entry.BusinessID,
		IntegrationID: entry.IntegrationID,
		ConfigID:      entry.ConfigID,
		Module:        emailLogModuleNotification,
		To:            entry.To,
		Subject:       entry.Subject,
		EventType:     entry.EventType,
		Status:        entry.Status,
		CreatedAt:     entry.SentAt,
	}

	if entry.ErrorMessage != "" {
		model.ErrorMessage = &entry.ErrorMessage
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("error creando email_log: %w", err)
	}

	return nil
}
