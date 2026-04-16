package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

// ListPayments lista las transacciones de pago de un negocio con paginación
func (uc *useCase) ListPayments(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.PaymentTransaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return uc.repo.ListPaymentTransactions(ctx, businessID, page, pageSize)
}
