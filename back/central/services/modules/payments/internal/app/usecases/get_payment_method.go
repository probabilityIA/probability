package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
)

// GetPaymentMethodByID obtiene un método de pago por ID
func (uc *UseCase) GetPaymentMethodByID(ctx context.Context, id uint) (*dtos.PaymentMethodResponse, error) {
	method, err := uc.repo.GetPaymentMethodByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mappers.EntityToResponse(method)
	return &response, nil
}

// GetPaymentMethodByCode obtiene un método de pago por código
func (uc *UseCase) GetPaymentMethodByCode(ctx context.Context, code string) (*dtos.PaymentMethodResponse, error) {
	method, err := uc.repo.GetPaymentMethodByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	response := mappers.EntityToResponse(method)
	return &response, nil
}
