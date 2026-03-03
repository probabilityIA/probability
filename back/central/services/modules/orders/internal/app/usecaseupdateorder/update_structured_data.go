package usecaseupdateorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/helpers/statusmapper"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// updateStructuredData actualiza los campos JSONB de la orden
func (uc *UseCaseUpdateOrder) updateStructuredData(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	changed := false

	// Actualizar Metadata si está presente
	if len(dto.Metadata) > 0 {
		if len(order.Metadata) == 0 || !statusmapper.EqualJSON(order.Metadata, dto.Metadata) {
			order.Metadata = dto.Metadata
			changed = true
		}
	}

	// Actualizar FinancialDetails si está presente
	if len(dto.FinancialDetails) > 0 {
		if len(order.FinancialDetails) == 0 || !statusmapper.EqualJSON(order.FinancialDetails, dto.FinancialDetails) {
			order.FinancialDetails = dto.FinancialDetails
			changed = true
		}
	}

	// Actualizar ShippingDetails si está presente
	if len(dto.ShippingDetails) > 0 {
		if len(order.ShippingDetails) == 0 || !statusmapper.EqualJSON(order.ShippingDetails, dto.ShippingDetails) {
			order.ShippingDetails = dto.ShippingDetails
			changed = true
		}
	}

	// Actualizar PaymentDetails si está presente
	if len(dto.PaymentDetails) > 0 {
		if len(order.PaymentDetails) == 0 || !statusmapper.EqualJSON(order.PaymentDetails, dto.PaymentDetails) {
			order.PaymentDetails = dto.PaymentDetails
			changed = true
		}
	}

	// Actualizar FulfillmentDetails si está presente
	if len(dto.FulfillmentDetails) > 0 {
		if len(order.FulfillmentDetails) == 0 || !statusmapper.EqualJSON(order.FulfillmentDetails, dto.FulfillmentDetails) {
			order.FulfillmentDetails = dto.FulfillmentDetails
			changed = true
		}
	}

	return changed
}
