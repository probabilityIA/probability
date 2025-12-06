package usecaseintegrations

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// DeactivateIntegration desactiva una integración
func (uc *IntegrationUseCase) DeactivateIntegration(ctx context.Context, id uint) error {
	ctx = log.WithFunctionCtx(ctx, "DeactivateIntegration")

	integration, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrIntegrationNotFound, err)
	}

	if !integration.IsActive {
		return nil // Ya está desactivada
	}

	isActive := false
	dto := domain.UpdateIntegrationDTO{
		IsActive:    &isActive,
		UpdatedByID: 0, // No actualizar updated_by_id en activar/desactivar
	}

	_, err = uc.UpdateIntegration(ctx, id, dto)
	if err != nil {
		return err
	}

	uc.log.Info(ctx).Uint("id", id).Msg("Integración desactivada exitosamente")

	return nil
}
