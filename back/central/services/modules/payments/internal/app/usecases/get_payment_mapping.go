package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
)

// GetPaymentMappingByID obtiene un mapeo por ID
func (uc *UseCase) GetPaymentMappingByID(ctx context.Context, id uint) (*dtos.PaymentMappingResponse, error) {
	mapping, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mappers.MappingEntityToResponse(mapping)
	return &response, nil
}

// GetPaymentMappingsByIntegrationType obtiene mapeos por tipo de integración
func (uc *UseCase) GetPaymentMappingsByIntegrationType(ctx context.Context, integrationType string) ([]dtos.PaymentMappingResponse, error) {
	mappings, err := uc.repo.GetPaymentMappingsByIntegrationTypeWithMethods(ctx, integrationType)
	if err != nil {
		return nil, err
	}

	return mappers.MappingEntitiesToResponses(mappings), nil
}

// GetAllPaymentMappingsGroupedByIntegration obtiene todos los mapeos agrupados por tipo de integración
func (uc *UseCase) GetAllPaymentMappingsGroupedByIntegration(ctx context.Context) ([]dtos.PaymentMappingsByIntegrationResponse, error) {
	mappings, _, err := uc.repo.ListPaymentMappingsWithMethods(ctx, nil)
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]dtos.PaymentMappingResponse)
	for _, mapping := range mappings {
		response := mappers.MappingEntityToResponse(&mapping)
		grouped[mapping.IntegrationType] = append(grouped[mapping.IntegrationType], response)
	}

	result := make([]dtos.PaymentMappingsByIntegrationResponse, 0, len(grouped))
	for integrationType, mappings := range grouped {
		result = append(result, dtos.PaymentMappingsByIntegrationResponse{
			IntegrationType: integrationType,
			Mappings:        mappings,
		})
	}

	return result, nil
}
