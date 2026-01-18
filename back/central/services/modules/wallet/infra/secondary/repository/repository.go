package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/wallet/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"gorm.io/gorm"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) domain.IWalletRepository {
	// Auto-migrate tables
	ctx := context.Background()
	err := database.Conn(ctx).AutoMigrate(&domain.Wallet{}, &domain.Transaction{})
	if err != nil {
		// Log error but don't panic, let the application continue
		// The error will surface when trying to use the tables
		println("ERROR: Failed to auto-migrate wallet tables:", err.Error())
	} else {
		println("INFO: Wallet tables migrated successfully")
	}
	return &Repository{
		db: database,
	}
}

func (r *Repository) GetWalletByBusinessID(ctx context.Context, businessID uint) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.Conn(ctx).Where("business_id = ?", businessID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *Repository) GetWalletByID(ctx context.Context, walletID uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.Conn(ctx).Where("id = ?", walletID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *Repository) CreateWallet(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.Conn(ctx).Create(wallet).Error
}

func (r *Repository) UpdateWallet(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.Conn(ctx).Save(wallet).Error
}

func (r *Repository) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	return r.db.Conn(ctx).Create(transaction).Error
}

func (r *Repository) GetTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.Conn(ctx).Where("wallet_id = ?", walletID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *Repository) GetAllWallets(ctx context.Context) ([]domain.Wallet, error) {
	var wallets []domain.Wallet
	err := r.db.Conn(ctx).Find(&wallets).Error
	return wallets, err
}

func (r *Repository) GetPendingTransactions(ctx context.Context) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.Conn(ctx).
		Where("status = ?", domain.TransactionStatusPending).
		Where("type = ?", domain.TransactionTypeRecharge).
		Order("created_at ASC").
		Find(&transactions).Error
	return transactions, err
}

func (r *Repository) GetTransactionByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	var tx domain.Transaction
	err := r.db.Conn(ctx).Where("id = ?", id).First(&tx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &tx, nil
}

func (r *Repository) UpdateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	return r.db.Conn(ctx).Save(transaction).Error
}
