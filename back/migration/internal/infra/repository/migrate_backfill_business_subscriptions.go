package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateBackfillBusinessSubscriptions(ctx context.Context) error {
	db := r.db.Conn(ctx)

	var basicType models.SubscriptionType
	if err := db.Where("code = ?", "basico").First(&basicType).Error; err != nil {
		return fmt.Errorf("failed to load default subscription type for backfill: %w", err)
	}

	var businesses []models.Business
	if err := db.Where("subscription_type_id IS NULL AND deleted_at IS NULL").Find(&businesses).Error; err != nil {
		return fmt.Errorf("failed to load businesses pending backfill: %w", err)
	}

	endDate := time.Now().AddDate(0, 1, 0)

	for _, b := range businesses {
		if err := db.Model(&models.Business{}).Where("id = ?", b.ID).Updates(map[string]interface{}{
			"subscription_type_id":  basicType.ID,
			"subscription_status":   "active",
			"subscription_end_date": endDate,
		}).Error; err != nil {
			return fmt.Errorf("failed to backfill business %d: %w", b.ID, err)
		}
	}

	return nil
}
