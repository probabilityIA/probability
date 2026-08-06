package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

func (uc *walletUseCase) AdminAdjustBalance(ctx context.Context, dto *dtos.AdminAdjustBalanceDTO) error {
	if dto.Amount == 0 {
		return fmt.Errorf("amount cannot be zero")
	}

	wallet, err := uc.GetWallet(ctx, dto.BusinessID)
	if err != nil {
		return fmt.Errorf("error getting wallet for business %d: %w", dto.BusinessID, err)
	}

	txType := entities.WalletTxTypeRecharge
	amount := dto.Amount
	if amount < 0 {
		txType = entities.WalletTxTypeUsage
		amount = -amount
	}

	concept := dto.Concept
	if !entities.ValidWalletTxConcept(concept) {
		concept = entities.WalletTxConceptOther
	}

	tx := &entities.WalletTransaction{
		ID:         uuid.New(),
		WalletID:   wallet.ID,
		Amount:     amount,
		Type:       txType,
		Status:     entities.WalletTxStatusCompleted,
		Concept:    concept,
		Reference:  fmt.Sprintf("ADMIN_ADJ_%s: %s", uuid.New().String()[:8], dto.Reference),
		BusinessID: dto.BusinessID,
		CreatedAt:  time.Now(),
	}

	if err := uc.repo.CreateWalletTransaction(ctx, tx); err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	wallet.Balance += dto.Amount
	if err := uc.repo.UpdateWallet(ctx, wallet); err != nil {
		return fmt.Errorf("error updating wallet balance: %w", err)
	}

	uc.log.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Float64("amount", dto.Amount).
		Float64("new_balance", wallet.Balance).
		Str("reference", dto.Reference).
		Msg("Admin adjusted wallet balance")

	if dto.Amount < 0 {
		uc.CheckLowBalanceForBusiness(ctx, dto.BusinessID)
	}

	return nil
}
