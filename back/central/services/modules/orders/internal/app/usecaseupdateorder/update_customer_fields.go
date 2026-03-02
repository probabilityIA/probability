package usecaseupdateorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// updateCustomerFields actualiza la información del cliente
func (uc *UseCaseUpdateOrder) updateCustomerFields(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	changed := false

	if dto.CustomerName != "" && order.CustomerName != dto.CustomerName {
		order.CustomerName = dto.CustomerName
		changed = true
	}

	if dto.CustomerEmail != "" && order.CustomerEmail != dto.CustomerEmail {
		order.CustomerEmail = dto.CustomerEmail
		changed = true
	}

	if dto.CustomerPhone != "" && order.CustomerPhone != dto.CustomerPhone {
		order.CustomerPhone = dto.CustomerPhone
		changed = true
	}

	if dto.CustomerDNI != "" && order.CustomerDNI != dto.CustomerDNI {
		order.CustomerDNI = dto.CustomerDNI
		changed = true
	}

	if dto.CustomerOrderCount != nil && order.CustomerOrderCount != *dto.CustomerOrderCount {
		order.CustomerOrderCount = *dto.CustomerOrderCount
		changed = true
	}

	if dto.CustomerTotalSpent != nil && order.CustomerTotalSpent != *dto.CustomerTotalSpent {
		order.CustomerTotalSpent = *dto.CustomerTotalSpent
		changed = true
	}

	return changed
}
