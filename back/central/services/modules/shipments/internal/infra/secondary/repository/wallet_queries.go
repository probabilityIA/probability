package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) DebitWalletForGuide(ctx context.Context, businessID uint, amount float64, trackingNumber string, shipmentID *uint) error {
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
