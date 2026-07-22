package usecaseintegrations

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

func (uc *IntegrationUseCase) GetIntegrationStats(ctx context.Context, businessID uint) ([]domain.IntegrationStats, error) {
	if cached, err := uc.cache.GetIntegrationStats(ctx, businessID); err == nil && cached != nil {
		return cached, nil
	}

	stats, err := uc.repo.GetIntegrationStats(ctx, businessID)
	if err != nil {
		return nil, err
	}

	if err := uc.cache.SetIntegrationStats(ctx, businessID, stats); err != nil {
		uc.log.Warn(ctx).Err(err).Uint("business_id", businessID).Msg("No se pudo cachear stats de integraciones")
	}

	return stats, nil
}
