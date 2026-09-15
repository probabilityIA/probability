package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type templateOverageBillingRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func (r *templateOverageBillingRepository) GetActiveTemplatePlanLimits(ctx context.Context, businessID uint) (time.Time, time.Time, *int, *float64, bool, error) {
	var sub models.BusinessSubscription
	err := r.db.Conn(ctx).
		Preload("SubscriptionType").
		Where("business_id = ? AND status = ?", businessID, "paid").
		Order("created_at desc").
		First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return time.Time{}, time.Time{}, nil, nil, false, nil
		}
		return time.Time{}, time.Time{}, nil, nil, false, err
	}
	if sub.SubscriptionType == nil {
		return time.Time{}, time.Time{}, nil, nil, false, nil
	}

	return sub.StartDate, sub.EndDate, sub.SubscriptionType.IncludedTemplates, sub.SubscriptionType.TemplateOveragePrice, true, nil
}

func (r *templateOverageBillingRepository) CountTemplatesInCycle(ctx context.Context, businessID uint, cycleStart, cycleEnd time.Time) (int64, error) {
	var count int64
	err := r.db.Conn(ctx).Table("whatsapp_templates").
		Where("business_id = ? AND deleted_at IS NULL AND created_at BETWEEN ? AND ?", businessID, cycleStart, cycleEnd).
		Count(&count).Error
	return count, err
}

func (r *templateOverageBillingRepository) DebitWalletForTemplateOverage(ctx context.Context, businessID uint, amount float64, templateID uint) error {
	reference := fmt.Sprintf("TEMPLATE_OVERAGE_%d", templateID)

	var existing int64
	if err := r.db.Conn(ctx).Model(&models.WalletTransaction{}).
		Where("business_id = ? AND concept = ? AND reference = ?", businessID, "TEMPLATE_OVERAGE", reference).
		Count(&existing).Error; err != nil {
		return fmt.Errorf("failed to check existing wallet transaction: %w", err)
	}
	if existing > 0 {
		return nil
	}

	var wallet models.Wallet
	if err := r.db.Conn(ctx).Where("business_id = ?", businessID).First(&wallet).Error; err != nil {
		return fmt.Errorf("wallet not found for business %d: %w", businessID, err)
	}

	tx := r.db.Conn(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	txRecord := &models.WalletTransaction{
		ID:         uuid.New(),
		WalletID:   wallet.ID,
		Amount:     amount,
		Type:       "USAGE",
		Status:     "COMPLETED",
		Concept:    "TEMPLATE_OVERAGE",
		Reference:  reference,
		BusinessID: businessID,
		CreatedAt:  time.Now(),
	}

	if err := tx.Create(txRecord).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create wallet transaction: %w", err)
	}

	wallet.Balance -= amount
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	return tx.Commit().Error
}
