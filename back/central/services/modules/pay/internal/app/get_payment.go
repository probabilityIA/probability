package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

// GetPayment obtiene una transacción de pago por ID
func (uc *useCase) GetPayment(ctx context.Context, id uint) (*entities.PaymentTransaction, error) {
	return uc.repo.GetPaymentTransactionByID(ctx, id)
}
