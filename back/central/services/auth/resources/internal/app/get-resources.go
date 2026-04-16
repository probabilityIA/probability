package app

import (
	"context"
	"math"

	"github.com/secamc93/probability/back/central/services/auth/resources/internal/domain"
)

// GetResources obtiene todos los recursos con filtros y paginación
func (uc *ResourceUseCase) GetResources(ctx context.Context, filters domain.ResourceFilters) (*domain.ResourceListDTO, error) {
	uc.logger.Info().Interface("filters", filters).Msg("Iniciando obtención de recursos")

	// Configurar valores por defecto para paginación
	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}

	// Obtener recursos del repositorio
	resources, total, err := uc.repository.GetResources(ctx, filters)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener recursos")
		return nil, err
	}

	// Convertir a DTOs
	var resourceDTOs []domain.ResourceDTO
	for _, resource := range resources {
		resourceDTOs = append(resourceDTOs, domain.ResourceDTO{
			ID:               resource.ID,
			Name:             resource.Name,
			Description:      resource.Description,
			BusinessTypeID:   resource.BusinessTypeID,
			BusinessTypeName: resource.BusinessTypeName,
			CreatedAt:        resource.CreatedAt,
			UpdatedAt:        resource.UpdatedAt,
		})
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(total) / float64(filters.PageSize)))

	result := &domain.ResourceListDTO{
		Resources:  resourceDTOs,
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	uc.logger.Info().
		Int64("total", total).
		Int("returned", len(resourceDTOs)).
		Int("page", filters.Page).
		Int("total_pages", totalPages).
		Msg("Recursos obtenidos exitosamente")

	return result, nil
}
