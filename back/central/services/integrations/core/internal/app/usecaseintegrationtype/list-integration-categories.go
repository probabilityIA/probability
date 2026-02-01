package usecaseintegrationtype

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

// ListIntegrationCategories lista todas las categorías de integración activas y visibles
func (uc *integrationTypeUseCase) ListIntegrationCategories(ctx context.Context) ([]*domain.IntegrationCategory, error) {
	uc.log.Info(ctx).Msg("Listing integration categories")

	categories, err := uc.repo.ListIntegrationCategories(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error listing integration categories")
		return nil, err
	}

	uc.log.Info(ctx).
		Int("count", len(categories)).
		Msg("Integration categories listed successfully")

	return categories, nil
}
