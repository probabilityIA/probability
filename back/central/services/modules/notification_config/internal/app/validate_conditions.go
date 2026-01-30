package app

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// ValidateConditions valida si una orden cumple las condiciones de una configuración
func (uc *useCase) ValidateConditions(
	config *entities.IntegrationNotificationConfig,
	orderStatus string,
	paymentMethodID uint,
) bool {
	// 1. Validar statuses (si hay filtro configurado)
	if len(config.Conditions.Statuses) > 0 {
		statusMatch := false
		for _, allowedStatus := range config.Conditions.Statuses {
			if orderStatus == allowedStatus {
				statusMatch = true
				break
			}
		}
		if !statusMatch {
			uc.logger.Debug().
				Str("order_status", orderStatus).
				Strs("allowed_statuses", config.Conditions.Statuses).
				Msg("Order status does not match allowed statuses")
			return false
		}
	}

	// 2. Validar payment_methods (si hay filtro configurado)
	if len(config.Conditions.PaymentMethods) > 0 {
		paymentMatch := false
		for _, pmID := range config.Conditions.PaymentMethods {
			if paymentMethodID == pmID {
				paymentMatch = true
				break
			}
		}
		if !paymentMatch {
			uc.logger.Debug().
				Uint("payment_method_id", paymentMethodID).
				Uints("allowed_payment_methods", config.Conditions.PaymentMethods).
				Msg("Payment method does not match allowed payment methods")
			return false
		}
	}

	uc.logger.Debug().
		Uint("config_id", config.ID).
		Str("order_status", orderStatus).
		Uint("payment_method_id", paymentMethodID).
		Msg("Order matches notification config conditions")

	return true
}
