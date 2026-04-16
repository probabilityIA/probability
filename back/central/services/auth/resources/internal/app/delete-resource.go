package app

import (
	"context"
	"fmt"
)

// DeleteResource elimina un recurso por su ID
func (uc *ResourceUseCase) DeleteResource(ctx context.Context, id uint) (string, error) {
	uc.logger.Info().Uint("resource_id", id).Msg("Iniciando eliminación de recurso")

	// Verificar que el recurso existe antes de eliminarlo
	existingResource, err := uc.repository.GetResourceByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("resource_id", id).Msg("Error al verificar existencia del recurso")
		return "", fmt.Errorf("recurso con ID %d no encontrado", id)
	}

	if existingResource == nil {
		uc.logger.Warn().Uint("resource_id", id).Msg("Recurso no encontrado para eliminar")
		return "", fmt.Errorf("recurso con ID %d no encontrado", id)
	}

	// TODO: Aquí se podría agregar validaciones adicionales, como:
	// - Verificar si el recurso está siendo usado en permisos
	// - Verificar si está asociado a roles activos
	// - Etc.

	// Eliminar recurso del repositorio
	message, err := uc.repository.DeleteResource(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("resource_id", id).Msg("Error al eliminar recurso")
		return "", fmt.Errorf("error al eliminar recurso: %w", err)
	}

	uc.logger.Info().
		Uint("resource_id", id).
		Str("resource_name", existingResource.Name).
		Msg("Recurso eliminado exitosamente")

	return message, nil
}
