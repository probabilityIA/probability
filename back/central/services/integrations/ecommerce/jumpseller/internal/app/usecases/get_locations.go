package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

func (uc *jumpsellerUseCase) GetLocations(ctx context.Context, integrationID string, businessID uint) (*domain.LocationsInfo, error) {
	_, cred, err := uc.resolveIntegrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	storeInfo, err := uc.client.GetStoreInfo(ctx, cred)
	if err != nil {
		return nil, err
	}

	locations, err := uc.client.GetLocations(ctx, cred)
	if err != nil {
		return nil, err
	}

	info := &domain.LocationsInfo{
		Locations:        locations,
		SubscriptionPlan: storeInfo.SubscriptionPlan,
		MultiLocation:    len(locations) > 1,
	}

	for _, location := range locations {
		if location.IsStockOrigin && location.Main {
			name := location.Name
			info.StockOriginName = name
			break
		}
	}
	if info.StockOriginName == "" && len(locations) == 1 {
		info.StockOriginName = locations[0].Name
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Int("locations", len(locations)).
		Str("plan", storeInfo.SubscriptionPlan).
		Bool("multi_location", info.MultiLocation).
		Msg("Bodegas de Jumpseller consultadas")

	return info, nil
}
