package app

import (
	"context"
	"slices"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

// HandleEvent procesa un evento: consulta configs en cache, rutea por canal
func (d *EventDispatcher) HandleEvent(ctx context.Context, event entities.Event) error {
	d.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", event.Type).
		Uint("business_id", event.BusinessID).
		Uint("integration_id", event.IntegrationID).
		Msg("Procesando evento en dispatcher")

	// Lookup configs en Redis cache
	configs, err := d.configCache.GetActiveConfigsByIntegrationAndTrigger(ctx, event.IntegrationID, event.Type)
	if err != nil {
		d.logger.Warn(ctx).
			Err(err).
			Uint("integration_id", event.IntegrationID).
			Str("event_type", event.Type).
			Msg("Error consultando configs, broadcast SSE por defecto")
		d.ssePublisher.PublishEvent(event)
		return nil
	}

	// Si no hay configs → broadcast SSE por defecto (backward compatible)
	if len(configs) == 0 {
		d.logger.Info(ctx).
			Str("event_id", event.ID).
			Str("event_type", event.Type).
			Uint("integration_id", event.IntegrationID).
			Msg("Sin configs de notificación para este evento, broadcast SSE por defecto")
		d.ssePublisher.PublishEvent(event)
		return nil
	}

	d.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", event.Type).
		Uint("integration_id", event.IntegrationID).
		Int("configs_count", len(configs)).
		Msg("Configs de notificación encontradas, ruteando por canal")

	// Para cada config habilitada → validar condiciones → rutear por canal
	ssePublished := false
	for _, config := range configs {
		// Validar condiciones (OrderStatusCodes)
		if !d.validateConditions(event, config) {
			d.logger.Debug(ctx).
				Uint("config_id", config.ID).
				Str("event_type", event.Type).
				Msg("Config no cumple condiciones, saltando")
			continue
		}

		switch config.NotificationTypeID {
		case dtos.NotificationTypeSSE:
			if !ssePublished {
				d.ssePublisher.PublishEvent(event)
				ssePublished = true
			}
			d.logger.Info(ctx).
				Uint("config_id", config.ID).
				Msg("Evento ruteado a SSE")

		case dtos.NotificationTypeWhatsApp:
			if err := d.channelPublisher.PublishToWhatsApp(ctx, event, config); err != nil {
				d.logger.Error(ctx).
					Err(err).
					Uint("config_id", config.ID).
					Msg("Error publicando a WhatsApp")
			} else {
				d.logger.Info(ctx).
					Uint("config_id", config.ID).
					Msg("Evento ruteado a WhatsApp")
			}

		case dtos.NotificationTypeEmail:
			if err := d.channelPublisher.PublishToEmail(ctx, event, config); err != nil {
				d.logger.Error(ctx).
					Err(err).
					Uint("config_id", config.ID).
					Msg("Error publicando a Email")
			} else {
				d.logger.Info(ctx).
					Uint("config_id", config.ID).
					Msg("Evento ruteado a Email")
			}

		default:
			d.logger.Warn(ctx).
				Uint("notification_type_id", config.NotificationTypeID).
				Msg("Tipo de notificación desconocido")
		}
	}

	// Si ninguna config era SSE, broadcast SSE por defecto
	if !ssePublished {
		d.ssePublisher.PublishEvent(event)
	}

	return nil
}

// validateConditions valida si un evento cumple las condiciones de una config
func (d *EventDispatcher) validateConditions(event entities.Event, config entities.CachedNotificationConfig) bool {
	// Si no hay filtros de estado configurados → aceptar todo
	if len(config.OrderStatusCodes) == 0 && len(config.OrderStatusIDs) == 0 {
		return true
	}

	// 1. Intentar validar por código de estado (current_status string de Changes)
	if len(config.OrderStatusCodes) > 0 {
		if status, ok := event.Data["current_status"]; ok {
			if statusStr, ok := status.(string); ok && statusStr != "" {
				return slices.Contains(config.OrderStatusCodes, statusStr)
			}
		}
	}

	// 2. Fallback: validar por ID de estado (order_status_id del snapshot)
	//    Esto cubre eventos como order.created donde current_status no está en Changes
	if len(config.OrderStatusIDs) > 0 {
		if statusID, ok := event.Data["order_status_id"]; ok {
			var orderStatusID uint
			switch v := statusID.(type) {
			case float64:
				orderStatusID = uint(v)
			case uint:
				orderStatusID = v
			case int:
				orderStatusID = uint(v)
			}

			if orderStatusID > 0 {
				for _, allowedID := range config.OrderStatusIDs {
					if orderStatusID == allowedID {
						return true
					}
				}
				return false
			}
		}
	}

	// Si no hay información de estado en el evento → no filtrar (backward compatible)
	return true
}
