package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/errors"
)

// CreatePaymentMapping crea un nuevo mapeo
func (uc *UseCase) CreatePaymentMapping(ctx context.Context, req *dtos.CreatePaymentMapping) (*dtos.PaymentMappingResponse, error) {
	exists, err := uc.repo.PaymentMappingExists(ctx, req.IntegrationType, req.OriginalMethod)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrPaymentMappingAlreadyExists
	}

	_, err = uc.repo.GetPaymentMethodByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, errors.ErrPaymentMethodNotFound
	}

	mapping := mappers.CreateMappingDTOToEntity(req)

	if err := uc.repo.CreatePaymentMapping(ctx, mapping); err != nil {
		return nil, err
	}

	created, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, mapping.ID)
	if err != nil {
		return nil, err
	}

	response := mappers.MappingEntityToResponse(created)
	return &response, nil
}
