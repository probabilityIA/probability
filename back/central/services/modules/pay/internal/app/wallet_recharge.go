package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	payerrs "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/errors"
)

// RechargeWallet crea una solicitud de recarga pendiente
func (uc *walletUseCase) RechargeWallet(ctx context.Context, dto *dtos.RechargeWalletDTO) (*entities.WalletTransaction, error) {
	uc.log.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Float64("amount", dto.Amount).
		Msg("Processing wallet recharge request")

	if dto.Amount < constants.WalletMinRechargeAmount {
		return nil, fmt.Errorf("%w: minimum is %d", payerrs.ErrMinimumRechargeAmount, constants.WalletMinRechargeAmount)
	}

	wallet, err := uc.GetWallet(ctx, dto.BusinessID)
	if err != nil {
		return nil, err
	}

	// Crear transacción PENDING con referencia manual
	tx := &entities.WalletTransaction{
		WalletID:  wallet.ID,
		Amount:    dto.Amount,
		Type:      entities.WalletTxTypeRecharge,
		Status:    entities.WalletTxStatusPending,
		Reference: "MANUAL_" + uuid.New().String(),
		QrCode:    "STATIC_QR",
	}

	if err := uc.repo.CreateWalletTransaction(ctx, tx); err != nil {
		return nil, err
	}

	uc.log.Info(ctx).
		Str("tx_id", tx.ID.String()).
		Float64("amount", dto.Amount).
		Msg("Wallet recharge request created (pending admin approval)")

	return tx, nil
}
