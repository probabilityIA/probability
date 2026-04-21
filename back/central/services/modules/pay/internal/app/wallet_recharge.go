package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

// RechargeWallet crea una solicitud de recarga pendiente
func (uc *walletUseCase) RechargeWallet(ctx context.Context, dto *dtos.RechargeWalletDTO) (*entities.WalletTransaction, error) {
	uc.log.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Float64("amount", dto.Amount).
		Msg("Processing wallet recharge request")

	if dto.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	wallet, err := uc.GetWallet(ctx, dto.BusinessID)
	if err != nil {
		return nil, err
	}

	// Crear transacción PENDING con referencia proporcionada
	reference := dto.Reference
	if reference == "" {
		reference = "MANUAL_" + uuid.New().String()
	}

	// Generar llave dinámica para Nequi (simulación)
	// En producción, esta llave vendría de la API de Nequi
	nequiKey := fmt.Sprintf("%010d", dto.BusinessID*1000000+uint(dto.Amount)%1000000)

	tx := &entities.WalletTransaction{
		WalletID:  wallet.ID,
		Amount:    dto.Amount,
		Type:      entities.WalletTxTypeRecharge,
		Status:    entities.WalletTxStatusPending,
		Reference: reference,
		QrCode:    nequiKey,
	}

	if err := uc.repo.CreateWalletTransaction(ctx, tx); err != nil {
		return nil, err
	}

	uc.log.Info(ctx).
		Str("tx_id", tx.ID.String()).
		Float64("amount", dto.Amount).
		Msg("Wallet recharge request created - awaiting payment confirmation from Nequi webhook or manual admin approval")

	return tx, nil
}
