package usecaseupdateorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// updateFinancialFields actualiza los campos financieros de la orden
func (uc *UseCaseUpdateOrder) updateFinancialFields(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	changed := false

	if dto.Subtotal > 0 && order.Subtotal != dto.Subtotal {
		order.Subtotal = dto.Subtotal
		changed = true
	}

	if dto.Tax >= 0 && order.Tax != dto.Tax {
		order.Tax = dto.Tax
		changed = true
	}

	if dto.Discount >= 0 && order.Discount != dto.Discount {
		order.Discount = dto.Discount
		changed = true
	}

	if dto.ShippingCost >= 0 && order.ShippingCost != dto.ShippingCost {
		order.ShippingCost = dto.ShippingCost
		changed = true
	}

	if dto.TotalAmount > 0 && order.TotalAmount != dto.TotalAmount {
		order.TotalAmount = dto.TotalAmount
		changed = true
	}

	if dto.Currency != "" && order.Currency != dto.Currency {
		order.Currency = dto.Currency
		changed = true
	}

	if dto.CodTotal != nil && (order.CodTotal == nil || *order.CodTotal != *dto.CodTotal) {
		order.CodTotal = dto.CodTotal
		changed = true
	}

	return changed
}
