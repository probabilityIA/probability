package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) CreateInvoiceSyncLog(ctx context.Context, log *entities.InvoiceSyncLog) error {
	model := mappers.SyncLogToModel(log)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create sync log: %w", err)
	}

	log.ID = model.ID
	return nil
}

func (r *Repository) GetSyncLogsByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.InvoiceSyncLog, error) {
	var models []*models.InvoiceSyncLog

	if err := r.db.Conn(ctx).
		Where("invoice_id = ?", invoiceID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get sync logs: %w", err)
	}

	return mappers.SyncLogListToDomain(models), nil
}

func (r *Repository) GetPendingSyncLogRetries(ctx context.Context, limit int) ([]*entities.InvoiceSyncLog, error) {
	var models []*models.InvoiceSyncLog

	now := time.Now()

	if err := r.db.Conn(ctx).
		Where("status = ? AND next_retry_at IS NOT NULL AND next_retry_at <= ? AND retry_count < max_retries", "failed", now).
		Order("next_retry_at ASC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending retries: %w", err)
	}

	return mappers.SyncLogListToDomain(models), nil
}

func (r *Repository) UpdateInvoiceSyncLog(ctx context.Context, log *entities.InvoiceSyncLog) error {
	model := mappers.SyncLogToModel(log)

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update sync log: %w", err)
	}

	return nil
}
