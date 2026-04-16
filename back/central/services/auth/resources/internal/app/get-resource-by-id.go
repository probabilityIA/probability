package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/resources/internal/domain"
)

// GetResourceByID obtiene un recurso por su ID
func (uc *ResourceUseCase) GetResourceByID(ctx context.Context, id uint) (*domain.ResourceDTO, error) {
	uc.logger.Info().Uint("resource_id", id).Msg("Iniciando obtención de recurso por ID")

	// Obtener recurso del repositorio
	resource, err := uc.repository.GetResourceByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("resource_id", id).Msg("Error al obtener recurso por ID")
		return nil, fmt.Errorf("error al obtener recurso: %w", err)
	}

	if resource == nil {
		uc.logger.Warn().Uint("resource_id", id).Msg("Recurso no encontrado")
		return nil, fmt.Errorf("recurso con ID %d no encontrado", id)
	}

	// Convertir a DTO
	resourceDTO := &domain.ResourceDTO{
		ID:               resource.ID,
		Name:             resource.Name,
		Description:      resource.Description,
		BusinessTypeID:   resource.BusinessTypeID,
		BusinessTypeName: resource.BusinessTypeName,
		CreatedAt:        resource.CreatedAt,
		UpdatedAt:        resource.UpdatedAt,
	}

	uc.logger.Info().
		Uint("resource_id", id).
		Str("name", resource.Name).
		Msg("Recurso obtenido exitosamente por ID")

	return resourceDTO, nil
}
