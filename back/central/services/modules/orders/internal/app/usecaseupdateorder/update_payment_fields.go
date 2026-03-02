package usecaseupdateorder

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// updatePaymentFields actualiza los campos relacionados con el pago
func (uc *UseCaseUpdateOrder) updatePaymentFields(_ context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	changed := false

	if len(dto.Payments) > 0 {
		payment := dto.Payments[0]

		if payment.PaymentMethodID > 0 && order.PaymentMethodID != payment.PaymentMethodID {
			order.PaymentMethodID = payment.PaymentMethodID
			changed = true
		}

		if payment.Status == "completed" && !order.IsPaid {
			order.IsPaid = true
			changed = true
		}

		if payment.PaidAt != nil && (order.PaidAt == nil || !order.PaidAt.Equal(*payment.PaidAt)) {
			order.PaidAt = payment.PaidAt
			changed = true
		}
	}

	return changed
}
