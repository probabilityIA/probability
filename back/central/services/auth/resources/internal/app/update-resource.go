package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/resources/internal/domain"
)

// UpdateResource actualiza un recurso existente
func (uc *ResourceUseCase) UpdateResource(ctx context.Context, id uint, updateDTO domain.UpdateResourceDTO) (*domain.ResourceDTO, error) {
	uc.logger.Info().Uint("resource_id", id).Str("name", updateDTO.Name).Msg("Iniciando actualización de recurso")

	// Validar datos de entrada
	if err := uc.validateUpdateResource(updateDTO); err != nil {
		uc.logger.Error().Err(err).Uint("resource_id", id).Msg("Validación fallida para actualizar recurso")
		return nil, err
	}

	// Verificar que el recurso existe
	existingResource, err := uc.repository.GetResourceByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("resource_id", id).Msg("Error al obtener recurso para actualizar")
		return nil, fmt.Errorf("recurso con ID %d no encontrado", id)
	}

	if existingResource == nil {
		uc.logger.Warn().Uint("resource_id", id).Msg("Recurso no encontrado para actualizar")
		return nil, fmt.Errorf("recurso con ID %d no encontrado", id)
	}

	// Verificar que no existe otro recurso con el mismo nombre (si se está cambiando el nombre)
	if strings.TrimSpace(updateDTO.Name) != existingResource.Name {
		duplicateResource, err := uc.repository.GetResourceByName(ctx, updateDTO.Name)
		if err == nil && duplicateResource != nil && duplicateResource.ID != id {
			uc.logger.Warn().Str("name", updateDTO.Name).Uint("resource_id", id).Msg("Otro recurso ya existe con ese nombre")
			return nil, fmt.Errorf("ya existe otro recurso con el nombre '%s'", updateDTO.Name)
		}
	}

	// Crear entidad de dominio con los datos actualizados
	resource := domain.Resource{
		ID:             id,
		Name:           strings.TrimSpace(updateDTO.Name),
		Description:    strings.TrimSpace(updateDTO.Description),
		BusinessTypeID: existingResource.BusinessTypeID, // Mantener el existente por defecto
	}

	// Actualizar business_type_id si se proporciona
	if updateDTO.BusinessTypeID != nil {
		resource.BusinessTypeID = *updateDTO.BusinessTypeID
	}

	// Actualizar recurso en el repositorio
	_, err = uc.repository.UpdateResource(ctx, id, resource)
	if err != nil {
		uc.logger.Error().Err(err).Uint("resource_id", id).Msg("Error al actualizar recurso")
		return nil, fmt.Errorf("error al actualizar recurso: %w", err)
	}

	// Obtener el recurso actualizado para devolver el DTO completo
	updatedResource, err := uc.repository.GetResourceByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("resource_id", id).Msg("Error al obtener recurso actualizado")
		return nil, fmt.Errorf("error al obtener recurso actualizado: %w", err)
	}

	// Convertir a DTO
	resourceDTO := &domain.ResourceDTO{
		ID:               updatedResource.ID,
		Name:             updatedResource.Name,
		Description:      updatedResource.Description,
		BusinessTypeID:   updatedResource.BusinessTypeID,
		BusinessTypeName: updatedResource.BusinessTypeName,
		CreatedAt:        updatedResource.CreatedAt,
		UpdatedAt:        updatedResource.UpdatedAt,
	}

	uc.logger.Info().
		Uint("resource_id", id).
		Str("name", updateDTO.Name).
		Msg("Recurso actualizado exitosamente")

	return resourceDTO, nil
}

// validateUpdateResource valida los datos para actualizar un recurso
func (uc *ResourceUseCase) validateUpdateResource(updateDTO domain.UpdateResourceDTO) error {
	if strings.TrimSpace(updateDTO.Name) == "" {
		return fmt.Errorf("el nombre del recurso es obligatorio")
	}

	if len(updateDTO.Name) > 100 {
		return fmt.Errorf("el nombre del recurso no puede exceder 100 caracteres")
	}

	if len(updateDTO.Description) > 500 {
		return fmt.Errorf("la descripción del recurso no puede exceder 500 caracteres")
	}

	return nil
}
