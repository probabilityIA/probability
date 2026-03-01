package sse

import (
	"context"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

// PublishEvent publica un evento para ser broadcast a conexiones SSE
func (m *EventManager) PublishEvent(event entities.Event) {
	if m.logger != nil {
		m.logger.Info(context.Background()).
			Str("event_id", event.ID).
			Str("event_type", event.Type).
			Uint("integration_id", event.IntegrationID).
			Uint("business_id", event.BusinessID).
			Msg("Publicando evento SSE")
	}

	businessID := event.BusinessID

	m.mutex.Lock()
	if _, ok := m.eventTypeCount[businessID]; !ok {
		m.eventTypeCount[businessID] = make(map[string]int)
	}
	m.eventCount[businessID]++
	m.eventTypeCount[businessID][event.Type]++
	m.mutex.Unlock()

	select {
	case m.eventChan <- event:
		if m.logger != nil {
			m.logger.Debug(context.Background()).
				Str("event_id", event.ID).
				Str("event_type", event.Type).
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
