package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
)

func (uc *tiktokUseCase) SaveSnapshot(ctx context.Context, snapshot *domain.ShopSnapshot) error {
	if snapshot == nil || snapshot.BusinessID == 0 || snapshot.IntegrationID == 0 {
		return fmt.Errorf("tiktok: invalid snapshot message")
	}

	if err := uc.repo.Save(ctx, snapshot); err != nil {
		uc.logger.Error(ctx).Err(err).
			Uint("integration_id", snapshot.IntegrationID).
			Msg("Error guardando snapshot de TikTok")
		return err
	}

	uc.logger.Info(ctx).
		Uint("integration_id", snapshot.IntegrationID).
		Str("shop_name", snapshot.ShopName).
		Msg("Snapshot de TikTok guardado")

	return nil
}

func (uc *tiktokUseCase) ListSnapshots(ctx context.Context, businessID, integrationID uint, page, pageSize int) ([]domain.ShopSnapshot, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.repo.ListByIntegration(ctx, businessID, integrationID, page, pageSize)
}
