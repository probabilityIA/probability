package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
)

func (uc *tiktokUseCase) SyncShopInfo(ctx context.Context, integrationID string, businessID uint) (*domain.ShopSnapshot, error) {
	integration, cred, err := uc.resolveIntegrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	shopName, _ := integration.Config["shop_name"].(string)
	if shopName == "" {
		shopName = integration.StoreID
	}
	if shopName == "" {
		return nil, domain.ErrMissingShopName
	}

	shops, err := uc.client.QueryShopInfo(ctx, cred, shopName, 1)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Str("shop_name", shopName).
			Msg("Error consultando shop info en TikTok")
		return nil, fmt.Errorf("querying tiktok shop info: %w", err)
	}
	if len(shops) == 0 {
		return nil, domain.ErrShopNotFound
	}

	shop := shops[0]
	snapshot := &domain.ShopSnapshot{
		BusinessID:           businessID,
		IntegrationID:        integration.ID,
		ShopID:               shop.ShopID,
		ShopName:             shop.ShopName,
		ShopRating:           shop.ShopRating,
		ShopReviewCount:      shop.ShopReviewCount,
		ItemSoldCount:        shop.ItemSoldCount,
		ShopPerformanceValue: shop.ShopPerformanceValue,
	}

	if err := uc.publisher.PublishSnapshot(ctx, snapshot); err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Msg("Error publicando snapshot de TikTok a la cola")
		return nil, fmt.Errorf("publishing tiktok snapshot: %w", err)
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("shop_name", shop.ShopName).
		Int64("shop_id", shop.ShopID).
		Msg("Snapshot de TikTok publicado a la cola")

	return snapshot, nil
}
