package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

func (uc *walletUseCase) BoldSimulatePayment(ctx context.Context, dto *dtos.BoldSimulateDTO) (*dtos.BoldSimulateResponse, error) {
	if dto == nil || dto.Amount <= 0 || dto.BusinessID == 0 || dto.OrderID == "" {
		return nil, fmt.Errorf("business_id, order_id and amount are required")
	}

	creds, err := uc.repo.GetBoldCredentials(ctx)
	if err != nil {
		return nil, err
	}
	if creds.Environment != "sandbox" {
		return nil, fmt.Errorf("simulate only allowed when bold integration is in sandbox mode")
	}

	tx, err := uc.repo.GetWalletTransactionByReference(ctx, dto.OrderID)
	if err != nil {
		return nil, fmt.Errorf("lookup wallet transaction: %w", err)
	}
	if tx == nil {
		wallet, walletErr := uc.GetWallet(ctx, dto.BusinessID)
		if walletErr != nil {
			return nil, walletErr
		}
		bizIntegration, _ := uc.repo.GetBoldIntegrationForBusiness(ctx, dto.BusinessID)
		var integrationTypeID, integrationID *uint
		if bizIntegration != nil {
			if bizIntegration.IntegrationTypeID != 0 {
				id := bizIntegration.IntegrationTypeID
				integrationTypeID = &id
			}
			if bizIntegration.IntegrationID != 0 {
				id := bizIntegration.IntegrationID
				integrationID = &id
			}
		}
		if integrationTypeID == nil && creds.IntegrationTypeID != 0 {
			id := creds.IntegrationTypeID
			integrationTypeID = &id
		}
		tx = &entities.WalletTransaction{
			WalletID:          wallet.ID,
			Amount:            dto.Amount,
			Type:              entities.WalletTxTypeRecharge,
			Status:            entities.WalletTxStatusPending,
			Reference:         dto.OrderID,
			IntegrationTypeID: integrationTypeID,
			IntegrationID:     integrationID,
		}
		if err := uc.repo.CreateWalletTransaction(ctx, tx); err != nil {
			return nil, err
		}
	}

	if tx.Status == entities.WalletTxStatusCompleted {
		updated, _ := uc.repo.GetWalletByID(ctx, tx.WalletID)
		return &dtos.BoldSimulateResponse{
			Success:       true,
			OrderID:       dto.OrderID,
			TransactionID: tx.ID.String(),
			Amount:        tx.Amount,
			NewBalance:    walletBalance(updated),
			Status:        "ALREADY_APPROVED",
		}, nil
	}

	if err := uc.approveTransactionInternal(ctx, tx); err != nil {
		return nil, err
	}

	updated, err := uc.repo.GetWalletByID(ctx, tx.WalletID)
	if err != nil {
		return nil, err
	}

	uc.log.Info(ctx).
		Str("order_id", dto.OrderID).
		Str("tx_id", tx.ID.String()).
		Float64("amount", dto.Amount).
		Float64("new_balance", updated.Balance).
		Msg("Bold sandbox payment simulated and credited")

	return &dtos.BoldSimulateResponse{
		Success:       true,
		OrderID:       dto.OrderID,
		TransactionID: tx.ID.String(),
		Amount:        tx.Amount,
		NewBalance:    updated.Balance,
		Status:        "APPROVED",
	}, nil
}

func walletBalance(w *entities.Wallet) float64 {
	if w == nil {
		return 0
	}
	return w.Balance
}
