package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	models "github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ──────────────────────────────────────────────────────
// PaymentTransactions
// ──────────────────────────────────────────────────────

func (r *Repository) CreatePaymentTransaction(ctx context.Context, tx *entities.PaymentTransaction) error {
	model := toModel(tx)
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create payment transaction: %w", err)
	}
	tx.ID = model.ID
	tx.CreatedAt = model.CreatedAt
	tx.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *Repository) GetPaymentTransactionByID(ctx context.Context, id uint) (*entities.PaymentTransaction, error) {
	var model models.PaymentTransaction
	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment transaction %d not found", id)
		}
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *Repository) GetPaymentTransactionByReference(ctx context.Context, ref string) (*entities.PaymentTransaction, error) {
	var model models.PaymentTransaction
	if err := r.db.Conn(ctx).Where("reference = ?", ref).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment transaction with reference %s not found", ref)
		}
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *Repository) UpdatePaymentTransaction(ctx context.Context, tx *entities.PaymentTransaction) error {
	model := toModel(tx)
	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update payment transaction: %w", err)
	}
	return nil
}

func (r *Repository) ListPaymentTransactions(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.PaymentTransaction, int64, error) {
	var modelList []models.PaymentTransaction
	var total int64

	query := r.db.Conn(ctx).Model(&models.PaymentTransaction{}).Where("business_id = ?", businessID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&modelList).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*entities.PaymentTransaction, len(modelList))
	for i, m := range modelList {
		result[i] = toDomain(&m)
	}
	return result, total, nil
}

// ──────────────────────────────────────────────────────
// PaymentSyncLogs
// ──────────────────────────────────────────────────────

func (r *Repository) CreateSyncLog(ctx context.Context, log *entities.PaymentSyncLog) error {
	model := toSyncLogModel(log)
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create payment sync log: %w", err)
	}
	log.ID = model.ID
	log.CreatedAt = model.CreatedAt
	return nil
}

func (r *Repository) UpdateSyncLog(ctx context.Context, log *entities.PaymentSyncLog) error {
	model := toSyncLogModel(log)
	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update payment sync log: %w", err)
	}
	return nil
}

func (r *Repository) GetPendingSyncLogRetries(ctx context.Context, limit int) ([]*entities.PaymentSyncLog, error) {
	var modelList []models.PaymentSyncLog
	now := time.Now()

	if err := r.db.Conn(ctx).
		Where("status = ?", constants.StatusFailed).
		Where("retry_count < ?", constants.MaxRetries).
		Where("next_retry_at <= ?", now).
		Limit(limit).
		Find(&modelList).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.PaymentSyncLog, len(modelList))
	for i, m := range modelList {
		result[i] = toSyncLogDomain(&m)
	}
	return result, nil
}

func (r *Repository) CancelPendingSyncLogs(ctx context.Context, transactionID uint) error {
	return r.db.Conn(ctx).
		Model(&models.PaymentSyncLog{}).
		Where("payment_transaction_id = ? AND status = ?", transactionID, constants.StatusProcessing).
		Update("status", constants.StatusCancelled).Error
}

func (r *Repository) GetSyncLogsByTransactionID(ctx context.Context, transactionID uint) ([]*entities.PaymentSyncLog, error) {
	var modelList []models.PaymentSyncLog
	if err := r.db.Conn(ctx).
		Where("payment_transaction_id = ?", transactionID).
		Order("created_at DESC").
		Find(&modelList).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.PaymentSyncLog, len(modelList))
	for i, m := range modelList {
		result[i] = toSyncLogDomain(&m)
	}
	return result, nil
}

// ──────────────────────────────────────────────────────
// Mappers
// ──────────────────────────────────────────────────────

func toModel(e *entities.PaymentTransaction) *models.PaymentTransaction {
	m := &models.PaymentTransaction{
		BusinessID:    e.BusinessID,
		Amount:        e.Amount,
		Currency:      e.Currency,
		Status:        string(e.Status),
		GatewayCode:   e.GatewayCode,
		ExternalID:    e.ExternalID,
		Reference:     e.Reference,
		PaymentMethod: e.PaymentMethod,
		Description:   e.Description,
		CallbackURL:   e.CallbackURL,
	}
	if e.ID != 0 {
		m.Model.ID = e.ID
	}
	return m
}

func toDomain(m *models.PaymentTransaction) *entities.PaymentTransaction {
	return &entities.PaymentTransaction{
		ID:            m.ID,
		BusinessID:    m.BusinessID,
		Amount:        m.Amount,
		Currency:      m.Currency,
		Status:        entities.PaymentStatus(m.Status),
		GatewayCode:   m.GatewayCode,
		ExternalID:    m.ExternalID,
		Reference:     m.Reference,
		PaymentMethod: m.PaymentMethod,
		Description:   m.Description,
		CallbackURL:   m.CallbackURL,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func toSyncLogModel(e *entities.PaymentSyncLog) *models.PaymentSyncLog {
	m := &models.PaymentSyncLog{
		PaymentTransactionID: e.PaymentTransactionID,
		Status:               e.Status,
		RetryCount:           e.RetryCount,
		ErrorMessage:         e.ErrorMessage,
		NextRetryAt:          e.NextRetryAt,
	}
	if e.ID != 0 {
		m.Model.ID = e.ID
	}
	return m
}

func toSyncLogDomain(m *models.PaymentSyncLog) *entities.PaymentSyncLog {
	return &entities.PaymentSyncLog{
		ID:                   m.ID,
		PaymentTransactionID: m.PaymentTransactionID,
		Status:               m.Status,
		RetryCount:           m.RetryCount,
		ErrorMessage:         m.ErrorMessage,
		NextRetryAt:          m.NextRetryAt,
		CreatedAt:            m.CreatedAt,
	}
}
