package consumerevent

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/consumer/consumerevent/request"
)

// handleOrderEvent procesa un evento de orden
func (c *consumer) handleOrderEvent(ctx context.Context, payload string) {
	// 1. Parse event
	var event request.OrderEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		c.logger.Error().
			Err(err).
			Str("payload", payload).
			Msg("Error parsing order event")
		return
	}

	c.logger.Info().
		Str("event_type", string(event.Type)).
		Str("order_id", event.OrderID).
		Interface("business_id", event.BusinessID).
		Msg("Received order event")

	// 2. Validar que el evento tenga business_id
	if event.BusinessID == nil {
		c.logger.Warn().
			Str("event_id", event.ID).
			Msg("Event has no business_id, skipping")
		return
	}

	// 3. Obtener integración WhatsApp del business
	whatsappIntegration, err := c.integrationRepo.GetWhatsAppByBusinessID(ctx, *event.BusinessID)
	if err != nil {
		c.logger.Debug().
			Err(err).
			Uint("business_id", *event.BusinessID).
			Msg("No WhatsApp integration found for business, skipping")
		return
	}

	if !whatsappIntegration.IsActive {
		c.logger.Debug().
			Uint("integration_id", whatsappIntegration.ID).
			Msg("WhatsApp integration is not active, skipping")
		return
	}

	// 4. Obtener configs activas para este trigger (usando adaptador directo)
	configs, err := c.notificationConfigRepo.GetActiveConfigsByIntegrationAndTrigger(
		ctx,
		whatsappIntegration.ID,
		string(event.Type), // "order.created", "order.updated", "order.status_changed"
	)

	if err != nil {
		c.logger.Error().
			Err(err).
			Uint("integration_id", whatsappIntegration.ID).
			Str("trigger", string(event.Type)).
			Msg("Error getting notification configs")
		return
	}

	if len(configs) == 0 {
		c.logger.Debug().
			Uint("integration_id", whatsappIntegration.ID).
			Str("trigger", string(event.Type)).
			Msg("No active notification configs found for trigger")
		return
	}

	// 5. Obtener orden completa (necesitamos PaymentMethodID y Status)
	order, err := c.orderRepo.GetByID(ctx, event.OrderID)
	if err != nil {
		c.logger.Error().
			Err(err).
			Str("order_id", event.OrderID).
			Msg("Error getting order data")
		return
	}

	// 6. Validar contra cada config (por prioridad)
	for _, config := range configs {
		// Validar condiciones - PASAR integration_id de la orden
		if c.notificationConfigRepo.ValidateConditions(&config, order.Status, order.PaymentMethodID, order.IntegrationID) {
			c.logger.Info().
				Uint("config_id", config.ID).
				Str("order_id", order.ID).
				Str("template", config.TemplateName).
				Uint("source_integration_id", order.IntegrationID).
				Msg("Order matches notification config, publishing to RabbitMQ")

			// Publicar a RabbitMQ
			if err := c.publishConfirmationRequest(ctx, order, &config); err != nil {
				c.logger.Error().
					Err(err).
					Str("order_id", order.ID).
					Msg("Error publishing confirmation request to RabbitMQ")
			}

			// Solo primera config que coincida
			break
		}
	}
}
