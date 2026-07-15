package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) consolidatePack(ctx context.Context, accessToken string, packID int64, primary *domain.MeliOrder) (*domain.MeliOrder, error) {
	pack, err := uc.client.GetPack(ctx, accessToken, packID)
	if err != nil {
		return nil, err
	}
	if pack == nil || len(pack.OrderIDs) <= 1 {
		return nil, nil
	}

	merged := *primary
	merged.ID = packID
	merged.OrderItems = nil
	merged.Payments = nil
	merged.TotalAmount = 0
	merged.CouponAmount = 0
	merged.Shipping = nil

	for _, childID := range pack.OrderIDs {
		var child *domain.MeliOrder
		if childID == primary.ID {
			child = primary
		} else {
			c, _, gerr := uc.client.GetOrder(ctx, accessToken, childID)
			if gerr != nil {
				uc.logger.Warn(ctx).Err(gerr).Int64("child_order_id", childID).Msg("Failed to fetch pack child order")
				continue
			}
			child = c
		}

		merged.OrderItems = append(merged.OrderItems, child.OrderItems...)
		merged.Payments = append(merged.Payments, child.Payments...)
		merged.TotalAmount += child.TotalAmount
		merged.CouponAmount += child.CouponAmount
		if merged.Shipping == nil && child.Shipping != nil {
			merged.Shipping = child.Shipping
		}
	}

	return &merged, nil
}
