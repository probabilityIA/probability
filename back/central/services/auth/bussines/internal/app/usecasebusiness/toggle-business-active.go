package usecasebusiness

import (
	"context"
	"errors"
)

// ToggleBusinessActive activa o desactiva un business
func (uc *BusinessUseCase) ToggleBusinessActive(ctx context.Context, businessID uint, active bool) error {
	uc.log.Info().Uint("business_id", businessID).Bool("active", active).Msg("Cambiando estado de activación del business")

	if _, err := uc.repository.GetBusinessByID(ctx, businessID); err != nil {
		uc.log.Error().Err(err).Uint("business_id", businessID).Msg("Error al verificar business")
		return errors.New("business no encontrado")
	}

	if err := uc.repository.ToggleBusinessActive(ctx, businessID, active); err != nil {
		uc.log.Error().Err(err).Uint("business_id", businessID).Msg("Error al cambiar estado del business")
		return err
	}

	uc.log.Info().Uint("business_id", businessID).Bool("active", active).Msg("Estado del business actualizado exitosamente")
	return nil
}
