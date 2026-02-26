package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/errors"
)

// CreatePaymentMethod crea un nuevo método de pago
func (uc *UseCase) CreatePaymentMethod(ctx context.Context, req *dtos.CreatePaymentMethod) (*dtos.PaymentMethodResponse, error) {
	exists, err := uc.repo.PaymentMethodExists(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrPaymentMethodCodeAlreadyExists
	}

	method := mappers.CreateDTOToEntity(req)

	if err := uc.repo.CreatePaymentMethod(ctx, method); err != nil {
		return nil, err
	}

	response := mappers.EntityToResponse(method)
	return &response, nil
}
