package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type WebhookLogRepository struct {
	db  db.IDatabase
	log log.ILogger
}

func NewWebhookLogRepository(database db.IDatabase, logger log.ILogger) domain.IWebhookLogRepository {
	return &WebhookLogRepository{
		db:  database,
		log: logger.WithModule("woocommerce.webhook_log_repository"),
	}
}

func (r *WebhookLogRepository) Save(ctx context.Context, entry *domain.WebhookLog) error {
	dbLog := &models.WebhookLog{
		Source:         domain.WebhookSourceWooCommerce,
		EventType:      entry.EventType,
		URL:            entry.URL,
		Method:         entry.Method,
		Headers:        entry.Headers,
		RequestBody:    entry.RequestBody,
		RemoteIP:       entry.RemoteIP,
		Status:         entry.Status,
		ResponseCode:   entry.ResponseCode,
		ErrorMessage:   entry.ErrorMessage,
		BusinessID:     entry.BusinessID,
		RawStatus:      entry.RawStatus,
		MappedStatus:   entry.MappedStatus,
		SignatureValid: entry.SignatureValid,
		RetentionUntil: entry.RetentionUntil,
	}

	if entry.IntegrationID != "" {
		if id, err := strconv.ParseUint(entry.IntegrationID, 10, 64); err == nil {
			intID := uint(id)
			dbLog.IntegrationID = &intID
		}
	}

	if err := r.db.Conn(ctx).Create(dbLog).Error; err != nil {
		return err
	}

	entry.ID = dbLog.ID
	entry.CreatedAt = dbLog.CreatedAt
	return nil
}

func (r *WebhookLogRepository) DeleteExpired(ctx context.Context) error {
	return r.db.Conn(ctx).
		Where("source = ? AND retention_until IS NOT NULL AND retention_until < ?", domain.WebhookSourceWooCommerce, time.Now()).
		Delete(&models.WebhookLog{}).Error
}
