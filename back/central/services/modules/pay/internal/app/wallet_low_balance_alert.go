package app

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const lowBalanceThreshold = 50000.0

type walletBalanceAlertMessage struct {
	BusinessID   uint    `json:"business_id"`
	BusinessName string  `json:"business_name"`
	PhoneNumber  string  `json:"phone_number"`
	Balance      float64 `json:"balance"`
}

func (uc *walletUseCase) CheckLowBalances(ctx context.Context) error {
	if uc.rabbit == nil {
		uc.log.Warn(ctx).Msg("RabbitMQ no disponible - omitiendo chequeo de saldo bajo")
		return nil
	}

	candidates, err := uc.repo.ListLowBalanceAlertCandidates(ctx, lowBalanceThreshold)
	if err != nil {
		return err
	}

	for _, c := range candidates {
		uc.publishLowBalanceAlert(ctx, c)
	}

	return nil
}

func (uc *walletUseCase) CheckLowBalanceForBusiness(ctx context.Context, businessID uint) {
	if uc.rabbit == nil {
		return
	}

	candidate, err := uc.repo.GetLowBalanceAlertCandidate(ctx, businessID, lowBalanceThreshold)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("Error verificando saldo bajo tras transaccion")
		return
	}
	if candidate == nil {
		return
	}

	uc.publishLowBalanceAlert(ctx, *candidate)
}

func (uc *walletUseCase) publishLowBalanceAlert(ctx context.Context, c entities.LowBalanceAlertCandidate) {
	if c.WarehousePhone == "" {
		uc.log.Warn(ctx).
			Uint("business_id", c.BusinessID).
			Msg("Negocio con saldo bajo sin telefono de bodega principal - omitiendo aviso")
		return
	}

	alreadyAlerted, err := uc.repo.HasRecentSystemAlert(ctx, c.WarehousePhone)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", c.BusinessID).Msg("Error verificando aviso reciente de saldo")
		return
	}
	if alreadyAlerted {
		return
	}

	msg := walletBalanceAlertMessage{
		BusinessID:   c.BusinessID,
		BusinessName: c.BusinessName,
		PhoneNumber:  c.WarehousePhone,
		Balance:      c.Balance,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	if err := uc.rabbit.DeclareQueue(rabbitmq.QueueWalletBalanceAlertRequested, true); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error declarando cola de aviso de saldo bajo")
		return
	}
	if err := uc.rabbit.Publish(ctx, rabbitmq.QueueWalletBalanceAlertRequested, data); err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", c.BusinessID).Msg("Error publicando aviso de saldo bajo")
		return
	}

	uc.log.Info(ctx).
		Uint("business_id", c.BusinessID).
		Float64("balance", c.Balance).
		Msg("Aviso de saldo bajo publicado")
}
