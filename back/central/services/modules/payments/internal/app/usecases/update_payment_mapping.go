package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/errors"
)

// UpdatePaymentMapping actualiza un mapeo existente
func (uc *UseCase) UpdatePaymentMapping(ctx context.Context, id uint, req *dtos.UpdatePaymentMapping) (*dtos.PaymentMappingResponse, error) {
	mapping, err := uc.repo.GetPaymentMappingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = uc.repo.GetPaymentMethodByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, errors.ErrPaymentMethodNotFound
	}

	mapping.OriginalMethod = req.OriginalMethod
	mapping.PaymentMethodID = req.PaymentMethodID
	mapping.Priority = req.Priority

	if err := uc.repo.UpdatePaymentMapping(ctx, mapping); err != nil {
		return nil, err
	}

	updated, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, mapping.ID)
	if err != nil {
		return nil, err
	}

	response := mappers.MappingEntityToResponse(updated)
	return &response, nil
}
