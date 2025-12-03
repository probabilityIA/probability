package usecaseintegrations

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// DeleteIntegration elimina una integración
func (uc *IntegrationUseCase) DeleteIntegration(ctx context.Context, id uint) error {
	ctx = log.WithFunctionCtx(ctx, "DeleteIntegration")

	// Verificar que existe
	integration, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrIntegrationNotFound, err)
	}

	// Validación: No se puede eliminar WhatsApp si es la única integración de ese tipo
	if integration.IntegrationType != nil && integration.IntegrationType.Code == "whatsapp" {
		return domain.ErrIntegrationCannotDeleteWhatsApp
	}

	// Eliminar
	if err := uc.repo.DeleteIntegration(ctx, id); err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al eliminar integración")
		return fmt.Errorf("error al eliminar integración: %w", err)
	}

	uc.log.Info(ctx).Uint("id", id).Msg("Integración eliminada exitosamente")

	return nil
}
