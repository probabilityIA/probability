package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
)

// UpdatePaymentMethod actualiza un método de pago existente
func (uc *UseCase) UpdatePaymentMethod(ctx context.Context, id uint, req *dtos.UpdatePaymentMethod) (*dtos.PaymentMethodResponse, error) {
	method, err := uc.repo.GetPaymentMethodByID(ctx, id)
	if err != nil {
		return nil, err
	}

	method.Name = req.Name
	method.Description = req.Description
	method.Category = req.Category
	method.Provider = req.Provider
	method.Icon = req.Icon
	method.Color = req.Color

	if err := uc.repo.UpdatePaymentMethod(ctx, method); err != nil {
		return nil, err
	}

	response := mappers.EntityToResponse(method)
	return &response, nil
}
