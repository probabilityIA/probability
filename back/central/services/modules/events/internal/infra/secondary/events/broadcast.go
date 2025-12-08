package events

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
)

// broadcastToBusinesses envía un evento a todas las conexiones que coincidan con el business_id y filtros
func (m *EventManager) broadcastToBusinesses(event domain.Event) {
	m.mutex.RLock()
	// Crear copia de conexiones para iterar sin bloqueo
	connectionsCopy := make(map[string]*domain.SSEConnection)
	for id, conn := range m.connections {
		connectionsCopy[id] = conn
	}
	m.mutex.RUnlock()

	// Extraer business_id del evento
	eventBusinessID := m.extractBusinessID(event)
	eventBusinessIDStr := event.BusinessID
	if eventBusinessIDStr == "" {
		eventBusinessIDStr = fmt.Sprintf("%d", eventBusinessID)
	}

	var aliveConnections []string
	sentCount := 0
	filteredCount := 0

	for connectionID, connection := range connectionsCopy {
		// Verificar si la conexión debe recibir este evento
		shouldReceive := false

		// Super usuario recibe todos los eventos (pero puede tener filtros)
		if connection.IsSuperUser() {
			shouldReceive = true
		} else {
			// Verificar si el business_id coincide
			if connection.MatchesBusiness(eventBusinessIDStr) {
				shouldReceive = true
			}
		}

		// Si debe recibir, aplicar filtros
		if shouldReceive {
			if connection.Filter != nil && !connection.Filter.Matches(event) {
				filteredCount++
				aliveConnections = append(aliveConnections, connectionID)
				continue
			}

			// Enviar evento a esta conexión
			if err := m.sendSSEMessage(connection.Writer, event); err != nil {
				if m.logger != nil {
					m.logger.Debug(context.Background()).
						Err(err).
						Str("connection_id", connectionID).
						Uint("business_id", connection.BusinessID).
						Msg("Removiendo conexión SSE rota")
				}
				continue
			}

			sentCount++
			aliveConnections = append(aliveConnections, connectionID)

			if m.logger != nil {
				m.logger.Debug(context.Background()).
					Str("connection_id", connectionID).
					Uint("business_id", connection.BusinessID).
					Str("event_type", string(event.Type)).
					Msg("Evento SSE enviado exitosamente")
			}
		} else {
			// Mantener conexión viva aunque no reciba este evento
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
			Str("event_type", string(event.Type)).
			Int("sent_count", sentCount).
			Int("filtered_count", filteredCount).
			Int("total_connections", len(connectionsCopy)).
			Msg("Evento broadcast a conexiones SSE")
	}
}
