package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
)

// ListPaymentMappings obtiene una lista de mapeos
func (uc *UseCase) ListPaymentMappings(ctx context.Context, filters map[string]interface{}) (*dtos.PaymentMappingsListResponse, error) {
	mappings, total, err := uc.repo.ListPaymentMappingsWithMethods(ctx, filters)
	if err != nil {
		return nil, err
	}

	data := mappers.MappingEntitiesToResponses(mappings)

	return &dtos.PaymentMappingsListResponse{
		Data:  data,
		Total: total,
	}, nil
}
