package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// DeleteIntegration elimina una integración
func (uc *integrationUseCase) DeleteIntegration(ctx context.Context, id uint) error {
	ctx = log.WithFunctionCtx(ctx, "DeleteIntegration")

	// Verificar que existe
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("integración no encontrada: %w", err)
	}

	// Validación: No se puede eliminar WhatsApp si es la única integración de ese tipo
	integration, _ := uc.repo.GetByID(ctx, id)
	if integration.Type == domain.IntegrationTypeWhatsApp {
		return fmt.Errorf("no se puede eliminar la integración de WhatsApp. Solo se puede desactivar")
	}

	// Eliminar
	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al eliminar integración")
		return fmt.Errorf("error al eliminar integración: %w", err)
	}

	uc.log.Info(ctx).Uint("id", id).Msg("Integración eliminada exitosamente")

	return nil
}
