package app

import (
	"context"
	"slices"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

func (d *EventDispatcher) HandleEvent(ctx context.Context, event entities.Event) error {
	d.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", event.Type).
		Uint("business_id", event.BusinessID).
		Uint("integration_id", event.IntegrationID).
		Msg("Procesando evento en dispatcher")

	configs, err := d.configCache.GetActiveConfigsByIntegrationAndTrigger(ctx, event.IntegrationID, event.Type)
	if err == nil {
		configs = d.applyVariants(ctx, event, configs)
	}
	if err != nil {
		d.logger.Warn(ctx).
			Err(err).
			Uint("integration_id", event.IntegrationID).
			Str("event_type", event.Type).
			Msg("Error consultando configs, broadcast SSE por defecto")
		d.ssePublisher.PublishEvent(event)
		return nil
	}

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

	ssePublished := false
	for _, config := range configs {
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

		case dtos.NotificationTypeAssistant:
			d.forwardToAssistant(ctx, event, config)

		case dtos.NotificationTypePush:
			if err := d.channelPublisher.PublishToPush(ctx, event, config); err != nil {
				d.logger.Error(ctx).
					Err(err).
					Uint("config_id", config.ID).
					Msg("Error publicando a Push")
			} else {
				d.logger.Info(ctx).
					Uint("config_id", config.ID).
					Msg("Evento ruteado a Push")
			}

		default:
			d.logger.Warn(ctx).
				Uint("notification_type_id", config.NotificationTypeID).
				Msg("Tipo de notificación desconocido")
		}
	}

	if !ssePublished {
		d.ssePublisher.PublishEvent(event)
	}

	return nil
}

func (d *EventDispatcher) validateConditions(event entities.Event, config entities.CachedNotificationConfig) bool {
	if isOrderConfirmation(config) && !event.IsCOD() {
		return false
	}

	if len(config.OrderStatusCodes) == 0 && len(config.OrderStatusIDs) == 0 {
		return true
	}

	if len(config.OrderStatusCodes) > 0 {
		if status, ok := event.Data["current_status"]; ok {
			if statusStr, ok := status.(string); ok && statusStr != "" {
				return slices.Contains(config.OrderStatusCodes, statusStr)
			}
		}
	}

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

	return true
}

func isOrderConfirmation(config entities.CachedNotificationConfig) bool {
	return config.NotificationTypeID == dtos.NotificationTypeWhatsApp &&
		config.EventCode == dtos.EventCodeOrderCreated
}

var eventVariants = map[string]string{
	dtos.EventCodeOrderCreated: "order.created_with_map",
}

func (d *EventDispatcher) applyVariants(ctx context.Context, event entities.Event, configs []entities.CachedNotificationConfig) []entities.CachedNotificationConfig {
	variantCode, tieneVariante := eventVariants[event.Type]
	if !tieneVariante {
		return configs
	}

	variantes, err := d.configCache.GetActiveConfigsByIntegrationAndTrigger(ctx, event.IntegrationID, variantCode)
	if err != nil || len(variantes) == 0 {
		return configs
	}

	canalesConVariante := make(map[uint]bool, len(variantes))
	for _, v := range variantes {
		canalesConVariante[v.NotificationTypeID] = true
	}

	resultado := make([]entities.CachedNotificationConfig, 0, len(configs)+len(variantes))
	for _, c := range configs {
		if canalesConVariante[c.NotificationTypeID] {
			d.logger.Info(ctx).
				Uint("config_id", c.ID).
				Str("variante", variantCode).
				Msg("La variante del evento esta activa, se omite la config base de ese canal")
			continue
		}
		resultado = append(resultado, c)
	}

	return append(resultado, variantes...)
}

func (d *EventDispatcher) forwardToAssistant(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) {
	if d.alerts == nil {
		return
	}
	if err := d.alerts.PublishToAssistant(ctx, event); err != nil {
		d.logger.Error(ctx).Err(err).Uint("config_id", config.ID).Str("event_type", event.Type).
			Msg("Error publicando la alerta al asistente")
		return
	}
	d.logger.Info(ctx).Uint("config_id", config.ID).Msg("Evento ruteado al asistente")
}
