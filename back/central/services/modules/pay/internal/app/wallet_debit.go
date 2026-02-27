package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	payerrs "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/errors"
)

// ManualDebit realiza un débito manual (admin)
func (uc *walletUseCase) ManualDebit(ctx context.Context, dto *dtos.ManualDebitDTO) error {
	if dto.Amount <= 0 {
		return payerrs.ErrInvalidAmount
	}

	wallet, err := uc.GetWallet(ctx, dto.BusinessID)
	if err != nil {
		return err
	}

	tx := &entities.WalletTransaction{
		WalletID:  wallet.ID,
		Amount:    dto.Amount,
		Type:      entities.WalletTxTypeUsage,
		Status:    entities.WalletTxStatusCompleted,
		Reference: "MAN_DEB_" + uuid.New().String()[:8] + ": " + dto.Reference,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.CreateWalletTransaction(ctx, tx); err != nil {
		return err
	}

	wallet.Balance -= dto.Amount
	if err := uc.repo.UpdateWallet(ctx, wallet); err != nil {
		return err
	}

	uc.log.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Float64("amount", dto.Amount).
		Float64("new_balance", wallet.Balance).
		Msg("Manual debit applied")

	return nil
}

// DebitForGuide realiza un débito para generación de guía
func (uc *walletUseCase) DebitForGuide(ctx context.Context, dto *dtos.DebitForGuideDTO) error {
	if dto.Amount <= 0 {
		return payerrs.ErrInvalidAmount
	}

	return uc.ManualDebit(ctx, &dtos.ManualDebitDTO{
		BusinessID: dto.BusinessID,
		Amount:     dto.Amount,
		Reference:  fmt.Sprintf("Guide generation: %s", dto.TrackingNumber),
	})
}
