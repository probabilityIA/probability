package usecasebusiness

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

// GetBusinessConfiguredResourcesByID obtiene un business por ID con sus recursos configurados
func (uc *BusinessUseCase) GetBusinessConfiguredResourcesByID(ctx context.Context, businessID uint) (*domain.BusinessWithConfiguredResourcesResponse, error) {
	uc.log.Info().Uint("business_id", businessID).Msg("Obteniendo business con recursos configurados por ID")

	business, err := uc.repository.GetBusinessByIDWithConfiguredResources(ctx, businessID)
	if err != nil {
		uc.log.Error().Err(err).Uint("business_id", businessID).Msg("Error al obtener business con recursos configurados")
		return nil, err
	}

	uc.log.Info().Uint("business_id", businessID).Int("resources_count", len(business.Resources)).Msg("Business con recursos configurados obtenido exitosamente")
	return business, nil
}
