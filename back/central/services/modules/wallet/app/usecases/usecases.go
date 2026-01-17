package usecases

import (
	"context"

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

	// 2. Generate QR
	qr, txID, err := u.nequiService.GenerateQR(ctx, amount)
	if err != nil {
		return "", err
	}

	// 3. Create Transaction
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

	// 4. Update Balance (Instant credit as requested)
	wallet.Balance += amount
	if err := u.repo.UpdateWallet(ctx, wallet); err != nil {
		return "", err
	}

	return qr, nil
}

func (u *WalletUsecases) GetAllWallets(ctx context.Context) ([]domain.Wallet, error) {
	return u.repo.GetAllWallets(ctx)
}
