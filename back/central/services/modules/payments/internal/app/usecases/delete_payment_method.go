package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/errors"
)

// DeletePaymentMethod elimina un método de pago
func (uc *UseCase) DeletePaymentMethod(ctx context.Context, id uint) error {
	hasActive, err := uc.repo.PaymentMethodHasActiveMappings(ctx, id)
	if err != nil {
		return err
	}
	if hasActive {
		return errors.ErrPaymentMethodHasActiveMappings
	}

	return uc.repo.DeletePaymentMethod(ctx, id)
}
