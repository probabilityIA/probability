package usecases

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type billingRetryMessage struct {
	IntegrationID string `json:"integration_id"`
	OrderID       int64  `json:"order_id"`
	Attempts      int    `json:"attempts"`
}

func (uc *meliUseCase) enrichBillingInfo(ctx context.Context, accessToken string, order *domain.MeliOrder) bool {
	if order.Buyer.BillingInfo != nil && order.Buyer.BillingInfo.DocNumber != "" {
		return false
	}
	info, err := uc.client.GetBillingInfo(ctx, accessToken, order.ID)
	if err != nil {
		if err == domain.ErrBillingInfoNotFound {
			return true
		}
		uc.logger.Warn(ctx).Err(err).Int64("order_id", order.ID).Msg("Failed to fetch billing_info")
		return false
	}
	order.Buyer.BillingInfo = info
	return false
}

func (uc *meliUseCase) enqueueBillingRetry(ctx context.Context, integrationID string, orderID int64, attempts int) {
	if uc.rabbit == nil {
		return
	}
	body, err := json.Marshal(billingRetryMessage{
		IntegrationID: integrationID,
		OrderID:       orderID,
		Attempts:      attempts,
	})
	if err != nil {
		return
	}
	if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueMeliBillingRetry, body); perr != nil {
		uc.logger.Warn(ctx).Err(perr).Int64("order_id", orderID).Msg("Failed to enqueue billing retry")
	}
}

func (uc *meliUseCase) RetryBilling(ctx context.Context, integrationID string, orderID int64) (bool, error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return false, err
	}
	if integration == nil {
		return false, domain.ErrIntegrationNotFound
	}
	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return false, err
	}
	_, err = uc.client.GetBillingInfo(ctx, accessToken, orderID)
	if err == domain.ErrBillingInfoNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if perr := uc.publishOrder(ctx, integration, accessToken, orderID); perr != nil {
		return true, perr
	}
	return true, nil
}
