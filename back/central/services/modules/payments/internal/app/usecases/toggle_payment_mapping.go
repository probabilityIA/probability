package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
)

// TogglePaymentMappingActive activa/desactiva un mapeo
func (uc *UseCase) TogglePaymentMappingActive(ctx context.Context, id uint) (*dtos.PaymentMappingResponse, error) {
	mapping, err := uc.repo.TogglePaymentMappingActive(ctx, id)
	if err != nil {
		return nil, err
	}

	updated, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, mapping.ID)
	if err != nil {
		return nil, err
	}

	response := mappers.MappingEntityToResponse(updated)
	return &response, nil
}
