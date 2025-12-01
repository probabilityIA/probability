package usecasebusinesstype

import (
	"context"
	"fmt"
)

// DeleteBusinessType elimina un tipo de negocio
func (uc *BusinessTypeUseCase) DeleteBusinessType(ctx context.Context, id uint) error {
	uc.log.Info().Uint("id", id).Msg("Eliminando tipo de negocio")

	// Verificar que existe
	existing, err := uc.repository.GetBusinessTypeByID(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al obtener tipo de negocio para eliminar")
		return fmt.Errorf("error al obtener tipo de negocio: %w", err)
	}

	if existing == nil {
		uc.log.Warn().Uint("id", id).Msg("Tipo de negocio no encontrado para eliminar")
		return fmt.Errorf("tipo de negocio no encontrado")
	}

	// Eliminar
	_, err = uc.repository.DeleteBusinessType(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al eliminar tipo de negocio")
		return fmt.Errorf("error al eliminar tipo de negocio: %w", err)
	}

	uc.log.Info().Uint("id", id).Msg("Tipo de negocio eliminado exitosamente")
	return nil
}
