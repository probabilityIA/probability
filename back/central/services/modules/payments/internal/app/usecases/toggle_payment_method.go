package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
)

// TogglePaymentMethodActive activa/desactiva un método de pago
func (uc *UseCase) TogglePaymentMethodActive(ctx context.Context, id uint) (*dtos.PaymentMethodResponse, error) {
	method, err := uc.repo.TogglePaymentMethodActive(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mappers.EntityToResponse(method)
	return &response, nil
}
