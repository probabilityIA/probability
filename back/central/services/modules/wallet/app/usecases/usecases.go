package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/wallet/domain"
)

type WalletUsecases struct {
	repo         domain.IWalletRepository
	nequiService domain.INequiService
}

func New(repo domain.IWalletRepository, nequiService domain.INequiService) *WalletUsecases {
	return &WalletUsecases{
		repo:         repo,
		nequiService: nequiService,
	}
}

func (u *WalletUsecases) GetWallet(ctx context.Context, businessID uint) (*domain.Wallet, error) {
	wallet, err := u.repo.GetWalletByBusinessID(ctx, businessID)
	if err != nil {
		return nil, err
	}
	// Create if not exists
	if wallet == nil {
		wallet = &domain.Wallet{
			BusinessID: businessID,
			Balance:    0,
		}
		if err = u.repo.CreateWallet(ctx, wallet); err != nil {
			return nil, err
		}
	}
	return wallet, nil
}

func (u *WalletUsecases) Recharge(ctx context.Context, businessID uint, amount float64) (string, error) {
	// 1. Get Wallet
	wallet, err := u.GetWallet(ctx, businessID)
	if err != nil {
		return "", err
	}

	// 2. Generate QR (SKIPPED - Static QR used on frontend)
	// User requested to remove Nequi API usage completely.
	// We generate a dummy reference for the transaction.
	qr := "STATIC_QR"
	txID := "MANUAL_" + uuid.New().String()

	// 3. Create Transaction (PENDING)
	tx := &domain.Transaction{
		WalletID:  wallet.ID,
		Amount:    amount,
		Type:      domain.TransactionTypeRecharge,
		Status:    domain.TransactionStatusPending,
		Reference: txID,
		QrCode:    qr,
	}

	if err := u.repo.CreateTransaction(ctx, tx); err != nil {
		return "", err
	}

	// DO NOT Update Balance automatically anymore

	// Return QR code (frontend might ignore it if it uses static, but good to have)
	return qr, nil
}

func (u *WalletUsecases) GetAllWallets(ctx context.Context) ([]domain.Wallet, error) {
	return u.repo.GetAllWallets(ctx)
}

func (u *WalletUsecases) GetPendingTransactions(ctx context.Context) ([]domain.Transaction, error) {
	// We need a repository method for this.
	// Assuming GetAllWallets exists, I need to check if GetTransactionsByWalletID can be used or if I need a new repo method.
	// Ideally I need GetPendingTransactions() on repo.
	// I'll assume I need to add it to the interface first.
	// But let's check ports.go again.
	// ports.go has: GetTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) ([]Transaction, error)
	// It does NOT have GetPendingTransactions().
	// I'll add the method to the repo requirement.
	// For now, I will write the usecase assuming the repo has it, and then update ports and repo implementation.
	return u.repo.GetPendingTransactions(ctx)
}

func (u *WalletUsecases) ApproveTransaction(ctx context.Context, transactionID string) error {
	// 1. Get Transaction
	// Need repo method GetTransactionByID
	tx, err := u.repo.GetTransactionByID(ctx, uuid.MustParse(transactionID))
	if err != nil {
		return err
	}

	if tx.Status != domain.TransactionStatusPending {
		// Already processed
		return nil // or error
	}

	// 2. Update Status
	tx.Status = domain.TransactionStatusCompleted
	if err := u.repo.UpdateTransaction(ctx, tx); err != nil {
		return err
	}

	// 3. Update Wallet Balance
	wallet, err := u.repo.GetWalletByID(ctx, tx.WalletID)
	if err != nil {
		return err
	}
	wallet.Balance += tx.Amount
	return u.repo.UpdateWallet(ctx, wallet)
}

func (u *WalletUsecases) RejectTransaction(ctx context.Context, transactionID string) error {
	tx, err := u.repo.GetTransactionByID(ctx, uuid.MustParse(transactionID))
	if err != nil {
		return err
	}

	if tx.Status != domain.TransactionStatusPending {
		return nil
	}

	tx.Status = domain.TransactionStatusFailed
	return u.repo.UpdateTransaction(ctx, tx)
}

func (u *WalletUsecases) GetProcessedTransactions(ctx context.Context) ([]domain.Transaction, error) {
	return u.repo.GetProcessedTransactions(ctx)
}

func (u *WalletUsecases) GetTransactionsByBusinessID(ctx context.Context, businessID uint) ([]domain.Transaction, error) {
	wallet, err := u.GetWallet(ctx, businessID)
	if err != nil {
		return nil, err
	}
	return u.repo.GetTransactionsByWalletID(ctx, wallet.ID)
}

func (u *WalletUsecases) ManualDebit(ctx context.Context, businessID uint, amount float64, reference string) error {
	wallet, err := u.GetWallet(ctx, businessID)
	if err != nil {
		return err
	}

	// Create a USAGE transaction
	tx := &domain.Transaction{
		WalletID:  wallet.ID,
		Amount:    amount,
		Type:      domain.TransactionTypeUsage,
		Status:    domain.TransactionStatusCompleted,
		Reference: "MAN_DEB_" + uuid.New().String()[:8] + ": " + reference,
		CreatedAt: time.Now(),
	}

	if err := u.repo.CreateTransaction(ctx, tx); err != nil {
		return err
	}

	// Update Wallet Balance
	wallet.Balance -= amount
	return u.repo.UpdateWallet(ctx, wallet)
}

func (u *WalletUsecases) ClearRechargeHistory(ctx context.Context, businessID uint) error {
	wallet, err := u.GetWallet(ctx, businessID)
	if err != nil {
		return err
	}

	return u.repo.DeleteTransactionsByWalletIDAndType(ctx, wallet.ID, domain.TransactionTypeRecharge)
}
