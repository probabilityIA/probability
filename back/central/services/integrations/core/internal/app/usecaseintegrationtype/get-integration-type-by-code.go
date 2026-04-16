package usecaseintegrationtype

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetIntegrationTypeByCode obtiene un tipo de integración por código
func (uc *integrationTypeUseCase) GetIntegrationTypeByCode(ctx context.Context, code string) (*domain.IntegrationType, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationTypeByCode")

	integrationType, err := uc.repo.GetIntegrationTypeByCode(ctx, code)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("code", code).Msg("Error al obtener tipo de integración por código")
		return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationTypeNotFound, err)
	}

	return integrationType, nil
}
