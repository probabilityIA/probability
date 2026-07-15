package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) PushOrderStatus(ctx context.Context, integrationID string, shipmentID int64, status string) error {
	mlStatus := mapProbabilityStatusToMeli(status)
	if mlStatus == "" {
		return nil
	}

	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return err
	}

	if err := uc.client.SendShipmentStatus(ctx, accessToken, shipmentID, mlStatus); err != nil {
		if err == domain.ErrTokenExpired {
			newToken, rerr := uc.EnsureValidToken(ctx, integrationID)
			if rerr != nil {
				return rerr
			}
			return uc.client.SendShipmentStatus(ctx, newToken, shipmentID, mlStatus)
		}
		return err
	}
	return nil
}

func mapProbabilityStatusToMeli(status string) string {
	switch status {
	case "shipped", "on_the_way", "in_transit":
		return "shipped"
	case "delivered":
		return "delivered"
	case "cancelled", "canceled":
		return "cancelled"
	default:
		return ""
	}
}
