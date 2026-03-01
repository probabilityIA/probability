package app

import (
	"context"
	"fmt"
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
		d.logger.Debug(ctx).
			Str("event_id", event.ID).
			Str("event_type", event.Type).
			Msg("Sin configs para este evento, broadcast SSE por defecto")
		d.ssePublisher.PublishEvent(event)
		return nil
	}

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
			d.handleEmailNotification(ctx, event, config)

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

// handleEmailNotification envía un email de notificación
func (d *EventDispatcher) handleEmailNotification(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) {
	if d.emailService == nil {
		d.logger.Warn(ctx).
			Uint("config_id", config.ID).
			Msg("Email service no disponible, saltando notificación email")
		return
	}

	// Extraer email del destinatario desde los datos del evento
	customerEmail := ""
	if email, ok := event.Data["customer_email"]; ok {
		if emailStr, ok := email.(string); ok {
			customerEmail = emailStr
		}
	}

	if customerEmail == "" {
		d.logger.Warn(ctx).
			Str("event_id", event.ID).
			Msg("No se encontró email del cliente en el evento, saltando notificación email")
		return
	}

	subject := fmt.Sprintf("Notificación: %s", event.Type)
	html := fmt.Sprintf("<h2>Evento: %s</h2><p>Se ha producido un evento en tu cuenta.</p>", event.Type)

	if err := d.emailService.SendHTML(ctx, customerEmail, subject, html); err != nil {
		d.logger.Error(ctx).
			Err(err).
			Str("to", customerEmail).
			Uint("config_id", config.ID).
			Msg("Error enviando email de notificación")
	} else {
		d.logger.Info(ctx).
			Str("to", customerEmail).
			Uint("config_id", config.ID).
			Msg("Email de notificación enviado")
	}
}

// validateConditions valida si un evento cumple las condiciones de una config
func (d *EventDispatcher) validateConditions(event entities.Event, config entities.CachedNotificationConfig) bool {
	// Si hay OrderStatusCodes configurados, verificar que el status actual esté en la lista
	if len(config.OrderStatusCodes) > 0 {
		currentStatus := ""
		if status, ok := event.Data["current_status"]; ok {
			if statusStr, ok := status.(string); ok {
				currentStatus = statusStr
			}
		}

		if currentStatus == "" {
			return true // Sin status en el evento → no filtrar
		}

		if !slices.Contains(config.OrderStatusCodes, currentStatus) {
			return false
		}
	}

	return true
}
