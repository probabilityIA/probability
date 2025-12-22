package events

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
)

// broadcastEvent envía un evento a todas las conexiones que coincidan con los filtros
func (m *IntegrationEventManager) broadcastEvent(event domain.IntegrationEvent) {
	// Actualizar estadísticas
	m.updateStats(event)

	// Guardar en caché de eventos recientes
	m.addToRecentEvents(event)

	// Determinar a qué businesses enviar el evento
	var businessIDs []uint
	if event.BusinessID != nil {
		// Enviar al business específico Y a conexiones de super usuario (business_id: 0)
		businessIDs = []uint{*event.BusinessID, 0}
	} else {
		// Si no hay business_id, enviar a todas las conexiones (super usuario)
		businessIDs = []uint{0}
	}

	// Broadcast a cada business
	for _, businessID := range businessIDs {
		m.broadcastToBusiness(businessID, event)
	}
}

// broadcastToBusiness envía un evento a todas las conexiones de un business
func (m *IntegrationEventManager) broadcastToBusiness(businessID uint, event domain.IntegrationEvent) {
	connections := m.connectionManager.GetConnectionsByBusiness(businessID)
	if businessID == 0 {
		// Para super usuario, obtener todas las conexiones
		connections = m.connectionManager.GetAllConnections()
	}

	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"location": "broadcast.go:broadcastToBusiness",
			"message":  "Broadcasting integration event",
			"data": map[string]interface{}{
				"business_id":       businessID,
				"event_type":        string(event.Type),
				"event_id":          event.ID,
				"integration_id":    event.IntegrationID,
				"total_connections": len(connections),
			},
			"timestamp":    time.Now().UnixMilli(),
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "B",
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	filteredCount := 0
	sentCount := 0

	for _, conn := range connections {
		// Aplicar filtros
		if !m.matchesFilter(conn, event) {
			continue
		}

		filteredCount++

		// Enviar evento a la conexión
		if m.sendEventToConnection(conn, event) {
			sentCount++
		}
	}

	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"location": "broadcast.go:broadcastToBusiness",
			"message":  "Broadcast completed",
			"data": map[string]interface{}{
				"business_id":    businessID,
				"event_type":     string(event.Type),
				"filtered_count": filteredCount,
				"sent_count":     sentCount,
			},
			"timestamp":    time.Now().UnixMilli(),
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "B",
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	if m.logger != nil {
		m.logger.Debug().
			Uint("business_id", businessID).
			Str("event_type", string(event.Type)).
			Int("total_connections", len(connections)).
			Int("filtered_count", filteredCount).
			Int("sent_count", sentCount).
			Msg("Event broadcast to business connections")
	}
}

// matchesFilter verifica si un evento coincide con los filtros de una conexión
func (m *IntegrationEventManager) matchesFilter(conn *domain.IntegrationSSEConnection, event domain.IntegrationEvent) bool {
	if conn.Filter == nil {
		return true // Sin filtros, aceptar todo
	}

	filter := conn.Filter

	// Filtro por integration_id
	if filter.IntegrationID != nil {
		if event.IntegrationID != *filter.IntegrationID {
			return false
		}
	}

	// Filtro por event_types
	if len(filter.EventTypes) > 0 {
		matched := false
		for _, eventType := range filter.EventTypes {
			if event.Type == eventType {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Filtro por order_ids (si el evento tiene order_id en metadata)
	if len(filter.OrderIDs) > 0 {
		if orderID, ok := event.Metadata["order_id"].(string); ok {
			matched := false
			for _, filterOrderID := range filter.OrderIDs {
				if orderID == filterOrderID {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		} else {
			return false // Si hay filtro de order_ids pero el evento no tiene order_id, rechazar
		}
	}

	return true
}

// sendEventToConnection envía un evento a una conexión SSE
func (m *IntegrationEventManager) sendEventToConnection(conn *domain.IntegrationSSEConnection, event domain.IntegrationEvent) bool {
	// Incrementar secuencia
	seq := m.connectionManager.NextSeq()
	conn.LastEventSeq = seq
	conn.EventCount++

	// Convertir evento a JSON
	eventJSON, err := json.Marshal(map[string]interface{}{
		"id":             event.ID,
		"type":           event.Type,
		"integration_id": event.IntegrationID,
		"business_id":    event.BusinessID,
		"timestamp":      event.Timestamp.Format(time.RFC3339),
		"data":           event.Data,
		"metadata":       event.Metadata,
		"seq":            seq,
	})
	if err != nil {
		if m.logger != nil {
			m.logger.Error().
				Err(err).
				Str("connection_id", conn.ID).
				Str("event_id", event.ID).
				Msg("Error serializing event for SSE")
		}
		return false
	}

	// Formatear mensaje SSE
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", event.Type, string(eventJSON))

	// Escribir al ResponseWriter
	if _, err := conn.Writer.Write([]byte(message)); err != nil {
		if m.logger != nil {
			m.logger.Warn().
				Err(err).
				Str("connection_id", conn.ID).
				Msg("Error writing event to SSE connection, removing connection")
		}
		m.connectionManager.RemoveConnection(conn.ID)
		return false
	}

	// Flush si es posible
	if flusher, ok := conn.Writer.(interface{ Flush() }); ok {
		flusher.Flush()
	}

	return true
}
