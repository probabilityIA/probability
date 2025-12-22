package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
)

// AddConnection agrega una nueva conexión SSE
func (m *IntegrationEventManager) AddConnection(businessID uint, filter *domain.IntegrationSSEFilter, conn http.ResponseWriter) string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.connectionCounter++
	connectionID := fmt.Sprintf("conn_%d_%d", time.Now().Unix(), m.connectionCounter)

	connection := &domain.IntegrationSSEConnection{
		ID:           connectionID,
		BusinessID:   businessID,
		Filter:       filter,
		Writer:       conn,
		CreatedAt:    time.Now(),
		LastEventSeq: 0,
		EventCount:   0,
	}

	m.connectionManager.AddConnection(connection)

	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"location": "connections.go:AddConnection",
			"message":  "SSE connection added",
			"data": map[string]interface{}{
				"connection_id": connectionID,
				"business_id":   businessID,
				"filter":        filter,
			},
			"timestamp":    time.Now().UnixMilli(),
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A",
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	if m.logger != nil {
		m.logger.Info().
			Str("connection_id", connectionID).
			Uint("business_id", businessID).
			Msg("New SSE connection added for integration events")
	}

	return connectionID
}

// RemoveConnection elimina una conexión SSE
func (m *IntegrationEventManager) RemoveConnection(connectionID string) {
	m.connectionManager.RemoveConnection(connectionID)

	if m.logger != nil {
		m.logger.Info().
			Str("connection_id", connectionID).
			Msg("SSE connection removed")
	}
}

// GetConnectionCount retorna el número de conexiones activas para un business
func (m *IntegrationEventManager) GetConnectionCount(businessID uint) int {
	if businessID == 0 {
		return m.connectionManager.GetConnectionCount()
	}
	return len(m.connectionManager.GetConnectionsByBusiness(businessID))
}

// GetConnectionInfo retorna información sobre las conexiones de un business
func (m *IntegrationEventManager) GetConnectionInfo(businessID uint) map[string]interface{} {
	connections := m.connectionManager.GetConnectionsByBusiness(businessID)
	if businessID == 0 {
		connections = m.connectionManager.GetAllConnections()
	}

	info := map[string]interface{}{
		"business_id": businessID,
		"total":       len(connections),
		"connections": make([]map[string]interface{}, 0, len(connections)),
	}

	for _, conn := range connections {
		connInfo := map[string]interface{}{
			"id":             conn.ID,
			"created_at":     conn.CreatedAt,
			"last_event_seq": conn.LastEventSeq,
			"event_count":    conn.EventCount,
		}
		if conn.Filter != nil {
			connInfo["filter"] = map[string]interface{}{
				"integration_id": conn.Filter.IntegrationID,
				"event_types":    conn.Filter.EventTypes,
				"order_ids":      conn.Filter.OrderIDs,
			}
		}
		info["connections"] = append(info["connections"].([]map[string]interface{}), connInfo)
	}

	return info
}
