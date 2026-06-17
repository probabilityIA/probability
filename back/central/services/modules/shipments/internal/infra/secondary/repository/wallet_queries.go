package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) DebitWalletForGuide(ctx context.Context, businessID uint, amount float64, trackingNumber string, shipmentID *uint) error {
	var wallet models.Wallet
	if err := r.db.Conn(ctx).Where("business_id = ?", businessID).First(&wallet).Error; err != nil {
		return fmt.Errorf("wallet not found for business %d: %w", businessID, err)
	}

	if shipmentID != nil {
		var existing int64
		if err := r.db.Conn(ctx).Model(&models.WalletTransaction{}).
			Where("shipment_id = ? AND type = ?", *shipmentID, "USAGE").
			Count(&existing).Error; err != nil {
			return fmt.Errorf("failed to check existing wallet transaction: %w", err)
		}
		if existing > 0 {
			return nil
		}
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
		Reference:  fmt.Sprintf("MAN_DEB_%s: Guide generation: %s", uuid.New().String()[:8], trackingNumber),
		ShipmentID: shipmentID,
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

func (r *Repository) FindUnchargedGuides(ctx context.Context, createdAfter, createdBefore time.Time, limit int) ([]domain.UnchargedGuide, error) {
	var rows []domain.UnchargedGuide
	err := r.db.Conn(ctx).
		Table("shipments s").
		Select("s.id AS shipment_id, o.business_id AS business_id, s.total_cost AS total_cost, s.tracking_number AS tracking_number").
		Joins("INNER JOIN orders o ON o.id = s.order_id").
		Where("s.tracking_number IS NOT NULL AND s.tracking_number <> ''").
		Where("s.total_cost IS NOT NULL AND s.total_cost > 0").
		Where("s.deleted_at IS NULL").
		Where("s.is_test = ?", false).
		Where("s.created_at >= ?", createdAfter).
		Where("s.created_at <= ?", createdBefore).
		Where("o.business_id > 0").
		Where("EXISTS (SELECT 1 FROM wallet w WHERE w.business_id = o.business_id)").
		Where("NOT EXISTS (SELECT 1 FROM transaction t WHERE t.shipment_id = s.id AND t.type = ?)", "USAGE").
		Order("s.created_at ASC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find uncharged guides: %w", err)
	}
	return rows, nil
}
