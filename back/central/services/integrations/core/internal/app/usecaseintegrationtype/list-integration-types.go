package usecaseintegrationtype

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// ListIntegrationTypes obtiene todos los tipos de integración, opcionalmente filtrados por categoría
func (uc *integrationTypeUseCase) ListIntegrationTypes(ctx context.Context, categoryID *uint) ([]*domain.IntegrationType, error) {
	ctx = log.WithFunctionCtx(ctx, "ListIntegrationTypes")

	integrationTypes, err := uc.repo.ListIntegrationTypes(ctx, categoryID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al listar tipos de integración")
		return nil, err
	}

	return integrationTypes, nil
}

// ListActiveIntegrationTypes obtiene solo los tipos de integración activos
func (uc *integrationTypeUseCase) ListActiveIntegrationTypes(ctx context.Context) ([]*domain.IntegrationType, error) {
	ctx = log.WithFunctionCtx(ctx, "ListActiveIntegrationTypes")

	integrationTypes, err := uc.repo.ListActiveIntegrationTypes(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al listar tipos de integración activos")
		return nil, err
	}

	return integrationTypes, nil
}
