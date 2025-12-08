package events

import (
	"context"
	"strconv"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
)

// PublishEvent publica un evento para ser broadcast
func (m *EventManager) PublishEvent(event domain.Event) {
	// Logging opcional si hay logger disponible
	if m.logger != nil {
		m.logger.Info(context.Background()).
			Str("event_id", event.ID).
			Str("event_type", string(event.Type)).
			Int64("integration_id", event.IntegrationID).
			Str("business_id", event.BusinessID).
			Msg("Publicando evento SSE")
	}

	// Extraer business_id del evento
	businessID := m.extractBusinessID(event)

	m.mutex.Lock()
	if _, ok := m.eventTypeCount[businessID]; !ok {
		m.eventTypeCount[businessID] = make(map[domain.EventType]int)
	}
	m.eventCount[businessID]++
	m.eventTypeCount[businessID][event.Type]++
	m.mutex.Unlock()

	select {
	case m.eventChan <- event:
		if m.logger != nil {
			m.logger.Debug(context.Background()).
				Str("event_id", event.ID).
				Str("event_type", string(event.Type)).
				Msg("Evento enviado al canal")
		}
	default:
		if m.logger != nil {
			m.logger.Warn(context.Background()).
				Interface("event", event).
				Msg("Canal de eventos lleno, descartando evento")
		}
	}
}

// extractBusinessID extrae el business_id de un evento
func (m *EventManager) extractBusinessID(event domain.Event) uint {
	// Intentar obtener business_id desde metadata
	if businessIDMeta, ok := event.Metadata["business_id"]; ok {
		if businessID, ok := businessIDMeta.(uint); ok {
			return businessID
		}
		if businessID, ok := businessIDMeta.(int); ok {
			return uint(businessID)
		}
		if businessID, ok := businessIDMeta.(int64); ok {
			return uint(businessID)
		}
		if businessID, ok := businessIDMeta.(float64); ok {
			return uint(businessID)
		}
	}

	// Intentar parsear desde BusinessID string
	if event.BusinessID != "" {
		if businessID, err := strconv.ParseUint(event.BusinessID, 10, 32); err == nil {
			return uint(businessID)
		}
	}

	return 0
}
