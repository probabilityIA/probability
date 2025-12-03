package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetIntegrationByID obtiene una integración por su ID
func (uc *integrationUseCase) GetIntegrationByID(ctx context.Context, id uint) (*domain.Integration, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationByID")

	integration, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener integración")
		return nil, err
	}

	return integration, nil
}
