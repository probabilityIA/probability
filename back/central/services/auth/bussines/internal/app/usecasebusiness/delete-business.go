package usecasebusiness

import (
	"context"
	"fmt"
)

// DeleteBusiness elimina un negocio
func (uc *BusinessUseCase) DeleteBusiness(ctx context.Context, id uint) error {
	uc.log.Info().Uint("id", id).Msg("Eliminando negocio")

	// Verificar que existe
	existing, err := uc.repository.GetBusinessByID(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al obtener negocio para eliminar")
		return fmt.Errorf("error al obtener negocio: %w", err)
	}

	if existing == nil {
		uc.log.Warn().Uint("id", id).Msg("Negocio no encontrado para eliminar")
		return fmt.Errorf("negocio no encontrado")
	}

	// Eliminar
	_, err = uc.repository.DeleteBusiness(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al eliminar negocio")
		return fmt.Errorf("error al eliminar negocio: %w", err)
	}

	uc.log.Info().Uint("id", id).Msg("Negocio eliminado exitosamente")
	return nil
}
