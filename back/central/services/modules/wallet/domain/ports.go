package domain

import (
	"context"

	"github.com/google/uuid"
)

type IWalletRepository interface {
	GetWalletByBusinessID(ctx context.Context, businessID uint) (*Wallet, error)
	GetWalletByID(ctx context.Context, walletID uuid.UUID) (*Wallet, error)
	CreateWallet(ctx context.Context, wallet *Wallet) error
	UpdateWallet(ctx context.Context, wallet *Wallet) error
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	GetTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) ([]Transaction, error)
	GetPendingTransactions(ctx context.Context) ([]Transaction, error)
	GetProcessedTransactions(ctx context.Context) ([]Transaction, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
	GetAllWallets(ctx context.Context) ([]Wallet, error)
}

type INequiService interface {
	GenerateQR(ctx context.Context, amount float64) (string, string, error)
}
