package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const lowBalanceThreshold = 50000.0
const eventTypeWalletLowBalance = "wallet.low_balance"

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

	integrationID, err := uc.repo.GetWhatsAppIntegrationID(ctx, c.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", c.BusinessID).Msg("Error resolviendo integracion de WhatsApp")
		return
	}
	if integrationID == nil {
		return
	}

	envelope := rabbitmq.EventEnvelope{
		Type:          eventTypeWalletLowBalance,
		Category:      eventCategoryPay,
		BusinessID:    c.BusinessID,
		IntegrationID: *integrationID,
		Data: map[string]interface{}{
			"business_name": c.BusinessName,
			"phone_number":  c.WarehousePhone,
			"balance":       c.Balance,
		},
	}
	if err := rabbitmq.PublishEvent(ctx, uc.rabbit, envelope); err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", c.BusinessID).Msg("Error publicando evento de saldo bajo")
		return
	}

	uc.log.Info(ctx).
		Uint("business_id", c.BusinessID).
		Float64("balance", c.Balance).
		Msg("Evento de saldo bajo publicado")
}
