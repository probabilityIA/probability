package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// DeactivateIntegration desactiva una integración
func (uc *integrationUseCase) DeactivateIntegration(ctx context.Context, id uint) error {
	ctx = log.WithFunctionCtx(ctx, "DeactivateIntegration")

	integration, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("integración no encontrada: %w", err)
	}

	if !integration.IsActive {
		return nil // Ya está desactivada
	}

	isActive := false
	dto := domain.UpdateIntegrationDTO{
		IsActive: &isActive,
	}

	_, err = uc.UpdateIntegration(ctx, id, dto)
	if err != nil {
		return err
	}

	uc.log.Info(ctx).Uint("id", id).Msg("Integración desactivada exitosamente")

	return nil
}
