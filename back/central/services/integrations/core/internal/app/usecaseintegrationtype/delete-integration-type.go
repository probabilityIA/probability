package usecaseintegrationtype

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// DeleteIntegrationType elimina un tipo de integración
func (uc *integrationTypeUseCase) DeleteIntegrationType(ctx context.Context, id uint) error {
	ctx = log.WithFunctionCtx(ctx, "DeleteIntegrationType")

	// Verificar que el tipo de integración exista
	_, err := uc.repo.GetIntegrationTypeByID(ctx, id)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener tipo de integración para eliminar")
		return fmt.Errorf("%w: %w", domain.ErrIntegrationTypeNotFound, err)
	}

	// TODO: Verificar que no haya integraciones usando este tipo antes de eliminar
	// Por ahora, GORM manejará la restricción de foreign key

	if err := uc.repo.DeleteIntegrationType(ctx, id); err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("id", id).
			Msg("Error al eliminar tipo de integración de la base de datos")
		return fmt.Errorf("error al eliminar tipo de integración: %w", err)
	}

	uc.log.Info(ctx).
		Uint("id", id).
		Msg("Tipo de integración eliminado exitosamente")

	return nil
}
