package app

import (
	"context"
	"fmt"
)

// DeleteAction elimina un action por su ID
func (uc *ActionUseCase) DeleteAction(ctx context.Context, id uint) (string, error) {
	uc.logger.Info().Uint("action_id", id).Msg("Iniciando eliminación de action")

	// Verificar que el action existe antes de eliminarlo
	existingAction, err := uc.repository.GetActionByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("action_id", id).Msg("Error al verificar existencia del action")
		return "", fmt.Errorf("action con ID %d no encontrado", id)
	}

	if existingAction == nil {
		uc.logger.Warn().Uint("action_id", id).Msg("Action no encontrado para eliminar")
		return "", fmt.Errorf("action con ID %d no encontrado", id)
	}

	// Eliminar action del repositorio (el repositorio ya valida permisos asociados)
	message, err := uc.repository.DeleteAction(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("action_id", id).Msg("Error al eliminar action")
		return "", fmt.Errorf("error al eliminar action: %w", err)
	}

	uc.logger.Info().
		Uint("action_id", id).
		Str("action_name", existingAction.Name).
		Msg("Action eliminado exitosamente")

	return message, nil
}
