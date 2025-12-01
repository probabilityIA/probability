package usecasebusiness

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

// GetBusinessesConfiguredResources obtiene businesses con sus recursos configurados con paginación
func (uc *BusinessUseCase) GetBusinessesConfiguredResources(ctx context.Context, page, perPage int, businessID *uint, businessTypeID *uint) ([]domain.BusinessWithConfiguredResourcesResponse, int64, error) {
	uc.log.Info().Int("page", page).Int("per_page", perPage).Msg("Obteniendo businesses con recursos configurados")

	businesses, total, err := uc.repository.GetBusinessesWithConfiguredResourcesPaginated(ctx, page, perPage, businessID, businessTypeID)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error al obtener businesses con recursos configurados")
		return nil, 0, err
	}

	uc.log.Info().Int64("total", total).Int("returned", len(businesses)).Msg("Businesses con recursos configurados obtenidos exitosamente")
	return businesses, total, nil
}
