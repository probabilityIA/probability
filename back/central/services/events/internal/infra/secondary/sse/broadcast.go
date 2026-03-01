package sse

import (
	"context"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

// broadcastToBusinesses envía un evento a todas las conexiones que coincidan
func (m *EventManager) broadcastToBusinesses(event entities.Event) {
	m.mutex.RLock()
	connectionsCopy := make(map[string]*entities.SSEConnection)
	for id, conn := range m.connections {
		connectionsCopy[id] = conn
	}
	m.mutex.RUnlock()

	eventBusinessID := event.BusinessID

	var aliveConnections []string
	sentCount := 0
	filteredCount := 0

	for connectionID, connection := range connectionsCopy {
		shouldReceive := false

		if connection.IsSuperUser() {
			shouldReceive = true
		} else if connection.MatchesBusiness(eventBusinessID) {
			shouldReceive = true
		}

		if shouldReceive {
			if connection.Filter != nil && !connection.Filter.Matches(event) {
				if m.logger != nil {
					m.logger.Debug(context.Background()).
						Str("connection_id", connectionID).
						Str("event_type", event.Type).
						Uint("business_id", connection.BusinessID).
						Msg("Evento filtrado por filtros de conexión SSE")
				}
				filteredCount++
				aliveConnections = append(aliveConnections, connectionID)
				continue
			}

			if err := m.sendSSEMessage(connection.Writer, event); err != nil {
				if m.logger != nil {
					m.logger.Error(context.Background()).
						Err(err).
						Str("connection_id", connectionID).
						Uint("business_id", connection.BusinessID).
						Str("event_type", event.Type).
						Msg("Error enviando evento SSE, removiendo conexión rota")
				}
				continue
			}

			sentCount++
			aliveConnections = append(aliveConnections, connectionID)

			if m.logger != nil {
				m.logger.Info(context.Background()).
					Str("connection_id", connectionID).
					Uint("business_id", connection.BusinessID).
					Str("event_type", event.Type).
					Str("event_id", event.ID).
					Msg("Evento SSE enviado exitosamente a conexión")
			}
		} else {
			aliveConnections = append(aliveConnections, connectionID)
		}
	}

	// Remover conexiones rotas
	m.mutex.Lock()
	for connectionID := range m.connections {
		found := false
		for _, aliveID := range aliveConnections {
			if connectionID == aliveID {
				found = true
				break
			}
		}
		if !found {
			delete(m.connections, connectionID)
		}
	}
	m.mutex.Unlock()

	if m.logger != nil {
		m.logger.Info(context.Background()).
			Uint("event_business_id", eventBusinessID).
			Str("event_type", event.Type).
			Int("sent_count", sentCount).
			Int("filtered_count", filteredCount).
			Int("total_connections", len(connectionsCopy)).
			Msg("Evento broadcast a conexiones SSE")
	}
}
