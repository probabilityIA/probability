package usecaseintegrationtype

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetIntegrationTypeByID obtiene un tipo de integración por ID
func (uc *integrationTypeUseCase) GetIntegrationTypeByID(ctx context.Context, id uint) (*domain.IntegrationType, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationTypeByID")

	integrationType, err := uc.repo.GetIntegrationTypeByID(ctx, id)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener tipo de integración por ID")
		return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationTypeNotFound, err)
	}

	return integrationType, nil
}
