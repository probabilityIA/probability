package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	models "github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ──────────────────────────────────────────────────────
// Wallet
// ──────────────────────────────────────────────────────

func (r *Repository) GetWalletByBusinessID(ctx context.Context, businessID uint) (*entities.Wallet, error) {
	var m models.Wallet
	err := r.db.Conn(ctx).Where("business_id = ?", businessID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return walletToDomain(&m), nil
}

func (r *Repository) GetWalletByID(ctx context.Context, walletID uuid.UUID) (*entities.Wallet, error) {
	var m models.Wallet
	err := r.db.Conn(ctx).Where("id = ?", walletID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, err
	}
	return walletToDomain(&m), nil
}

func (r *Repository) CreateWallet(ctx context.Context, wallet *entities.Wallet) error {
	m := walletToModel(wallet)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}
	wallet.ID = m.ID
	wallet.CreatedAt = m.CreatedAt
	wallet.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *Repository) UpdateWallet(ctx context.Context, wallet *entities.Wallet) error {
	m := walletToModel(wallet)
	if err := r.db.Conn(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}
	return nil
}

func (r *Repository) GetAllWallets(ctx context.Context) ([]*entities.Wallet, error) {
	var list []models.Wallet
	if err := r.db.Conn(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	result := make([]*entities.Wallet, len(list))
	for i, m := range list {
		result[i] = walletToDomain(&m)
	}
	return result, nil
}

// ──────────────────────────────────────────────────────
// WalletTransactions
// ──────────────────────────────────────────────────────

func (r *Repository) CreateWalletTransaction(ctx context.Context, tx *entities.WalletTransaction) error {
	m := walletTxToModel(tx)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("failed to create wallet transaction: %w", err)
	}
	tx.ID = m.ID
	tx.CreatedAt = m.CreatedAt
	return nil
}

func (r *Repository) GetWalletTransactionByID(ctx context.Context, id uuid.UUID) (*entities.WalletTransaction, error) {
	var m models.WalletTransaction
	err := r.db.Conn(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wallet transaction not found")
		}
		return nil, err
	}
	return walletTxToDomain(&m), nil
}

func (r *Repository) UpdateWalletTransaction(ctx context.Context, tx *entities.WalletTransaction) error {
	m := walletTxToModel(tx)
	if err := r.db.Conn(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("failed to update wallet transaction: %w", err)
	}
	return nil
}

func (r *Repository) GetTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) ([]*entities.WalletTransaction, error) {
	var list []models.WalletTransaction
	err := r.db.Conn(ctx).
		Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return walletTxListToDomain(list), nil
}

func (r *Repository) GetPendingRechargeTransactions(ctx context.Context) ([]*entities.WalletTransaction, error) {
	var list []models.WalletTransaction
	err := r.db.Conn(ctx).
		Where("status = ? AND type = ?", entities.WalletTxStatusPending, entities.WalletTxTypeRecharge).
		Order("created_at ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return walletTxListToDomain(list), nil
}

func (r *Repository) GetProcessedTransactions(ctx context.Context) ([]*entities.WalletTransaction, error) {
	var list []models.WalletTransaction
	err := r.db.Conn(ctx).
		Where("status IN ?", []string{entities.WalletTxStatusCompleted, entities.WalletTxStatusFailed}).
		Order("created_at DESC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return walletTxListToDomain(list), nil
}

func (r *Repository) DeleteTransactionsByWalletIDAndType(ctx context.Context, walletID uuid.UUID, txType string) error {
	return r.db.Conn(ctx).
		Where("wallet_id = ? AND type = ?", walletID, txType).
		Delete(&models.WalletTransaction{}).Error
}

// ──────────────────────────────────────────────────────
// Wallet Mappers
// ──────────────────────────────────────────────────────

func walletToModel(e *entities.Wallet) *models.Wallet {
	return &models.Wallet{
		ID:         e.ID,
		BusinessID: e.BusinessID,
		Balance:    e.Balance,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

func walletToDomain(m *models.Wallet) *entities.Wallet {
	return &entities.Wallet{
		ID:         m.ID,
		BusinessID: m.BusinessID,
		Balance:    m.Balance,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func walletTxToModel(e *entities.WalletTransaction) *models.WalletTransaction {
	return &models.WalletTransaction{
		ID:                   e.ID,
		WalletID:             e.WalletID,
		Amount:               e.Amount,
		Type:                 e.Type,
		Status:               e.Status,
		Reference:            e.Reference,
		QrCode:               e.QrCode,
		PaymentTransactionID: e.PaymentTransactionID,
		CreatedAt:            e.CreatedAt,
	}
}

func walletTxToDomain(m *models.WalletTransaction) *entities.WalletTransaction {
	return &entities.WalletTransaction{
		ID:                   m.ID,
		WalletID:             m.WalletID,
		Amount:               m.Amount,
		Type:                 m.Type,
		Status:               m.Status,
		Reference:            m.Reference,
		QrCode:               m.QrCode,
		PaymentTransactionID: m.PaymentTransactionID,
		CreatedAt:            m.CreatedAt,
	}
}

func walletTxListToDomain(list []models.WalletTransaction) []*entities.WalletTransaction {
	result := make([]*entities.WalletTransaction, len(list))
	for i, m := range list {
		result[i] = walletTxToDomain(&m)
	}
	return result
}
