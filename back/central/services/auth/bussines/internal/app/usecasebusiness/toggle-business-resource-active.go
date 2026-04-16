package usecasebusiness

import (
	"context"
	"errors"
)

// ToggleBusinessResourceActive activa o desactiva un recurso para un business específico
func (uc *BusinessUseCase) ToggleBusinessResourceActive(ctx context.Context, businessID uint, resourceID uint, active bool) error {
	uc.log.Info().Uint("business_id", businessID).Uint("resource_id", resourceID).Bool("active", active).Msg("Cambiando estado de activación del recurso")

	// Verificar que el business existe
	business, err := uc.repository.GetBusinessByID(ctx, businessID)
	if err != nil {
		uc.log.Error().Err(err).Uint("business_id", businessID).Msg("Error al verificar business")
		return errors.New("business no encontrado")
	}

	// Verificar que el recurso existe
	resource, err := uc.repository.GetResourceByID(ctx, resourceID)
	if err != nil {
		uc.log.Error().Err(err).Uint("resource_id", resourceID).Msg("Error al verificar recurso")
		return errors.New("recurso no encontrado")
	}

	// Verificar que el recurso está permitido para el tipo de business del business
	permittedResources, err := uc.repository.GetBusinessTypeResourcesPermitted(ctx, business.BusinessTypeID)
	if err != nil {
		uc.log.Error().Err(err).Uint("business_type_id", business.BusinessTypeID).Msg("Error al obtener recursos permitidos")
		return errors.New("error al verificar recursos permitidos")
	}

	// Verificar que el recurso está permitido
	resourcePermitted := false
	for _, permitted := range permittedResources {
		if permitted.ResourceID == resource.ID {
			resourcePermitted = true
			break
		}
	}

	if !resourcePermitted {
		uc.log.Error().Uint("resource_id", resourceID).Uint("business_type_id", business.BusinessTypeID).Msg("Recurso no permitido para este tipo de business")
		return errors.New("el recurso no está permitido para este tipo de business")
	}

	// Toggle el estado activo
	if err := uc.repository.ToggleBusinessResourceActive(ctx, businessID, resourceID, active); err != nil {
		uc.log.Error().Err(err).Uint("business_id", businessID).Uint("resource_id", resourceID).Msg("Error al cambiar estado del recurso")
		return err
	}

	uc.log.Info().Uint("business_id", businessID).Uint("resource_id", resourceID).Bool("active", active).Msg("Estado del recurso actualizado exitosamente")
	return nil
}
