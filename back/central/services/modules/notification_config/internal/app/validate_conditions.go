package app

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// ValidateConditions valida si una orden cumple las condiciones de una configuración
// NUEVA ESTRUCTURA: Usa OrderStatusIDs en lugar de nested Conditions
func (uc *useCase) ValidateConditions(
	config *entities.IntegrationNotificationConfig,
	orderStatusID uint,
	paymentMethodID uint,
) bool {
	// 1. Validar order statuses (si hay filtro configurado)
	// Si OrderStatusIDs está vacío, se aceptan todos los estados
	if len(config.OrderStatusIDs) > 0 {
		statusMatch := false
		for _, allowedStatusID := range config.OrderStatusIDs {
			if orderStatusID == allowedStatusID {
				statusMatch = true
				break
			}
		}
		if !statusMatch {
			uc.logger.Debug().
				Uint("order_status_id", orderStatusID).
				Uints("allowed_status_ids", config.OrderStatusIDs).
				Msg("Order status does not match allowed statuses")
			return false
		}
	}

	// 2. Validar payment_methods - TODO: Implementar cuando se migre payment_methods a nueva estructura
	// Por ahora, aceptamos todos los métodos de pago

	uc.logger.Debug().
		Uint("config_id", config.ID).
		Uint("order_status_id", orderStatusID).
		Uint("payment_method_id", paymentMethodID).
		Msg("Order matches notification config conditions")

	return true
}
