package usecaseintegrations

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// ActivateIntegration activa una integración
func (uc *IntegrationUseCase) ActivateIntegration(ctx context.Context, id uint) error {
	ctx = log.WithFunctionCtx(ctx, "ActivateIntegration")

	integration, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrIntegrationNotFound, err)
	}

	if integration.IsActive {
		return nil // Ya está activa
	}

	isActive := true
	dto := domain.UpdateIntegrationDTO{
		IsActive: &isActive,
	}

	_, err = uc.UpdateIntegration(ctx, id, dto)
	if err != nil {
		return err
	}

	uc.log.Info(ctx).Uint("id", id).Msg("Integración activada exitosamente")

	return nil
}
