package usecaseintegrations

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// ListIntegrations lista integraciones con filtros
func (uc *IntegrationUseCase) ListIntegrations(ctx context.Context, filters domain.IntegrationFilters) ([]*domain.Integration, int64, error) {
	ctx = log.WithFunctionCtx(ctx, "ListIntegrations")

	// Aplicar valores por defecto
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10
	}
	if filters.PageSize > 100 {
		filters.PageSize = 100
	}

	integrations, total, err := uc.repo.ListIntegrations(ctx, filters)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al listar integraciones")
		return nil, 0, err
	}

	uc.log.Info(ctx).
		Int("count", len(integrations)).
		Int64("total", total).
		Msg("Integraciones listadas exitosamente")

	return integrations, total, nil
}
