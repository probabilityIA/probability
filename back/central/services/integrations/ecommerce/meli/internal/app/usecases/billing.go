package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) enrichBillingInfo(ctx context.Context, accessToken string, order *domain.MeliOrder) {
	if order.Buyer.BillingInfo != nil && order.Buyer.BillingInfo.DocNumber != "" {
		return
	}
	info, err := uc.client.GetBillingInfo(ctx, accessToken, order.ID)
	if err != nil {
		if err != domain.ErrBillingInfoNotFound {
			uc.logger.Warn(ctx).Err(err).Int64("order_id", order.ID).Msg("Failed to fetch billing_info")
		}
		return
	}
	order.Buyer.BillingInfo = info
}
