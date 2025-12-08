package events

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
)

// AddConnection agrega una nueva conexión SSE por business_id con filtros opcionales
func (m *EventManager) AddConnection(businessID uint, filter *domain.SSEConnectionFilter, conn http.ResponseWriter) string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Generar ID único para la conexión
	connectionID := fmt.Sprintf("conn_%d_%d", businessID, atomic.AddUint64(&m.connectionCounter, 1))

	// Crear conexión con filtros
	connection := &domain.SSEConnection{
		BusinessID:   businessID,
		Filter:      filter,
		Writer:      conn,
		ConnectionID: connectionID,
	}

	m.connections[connectionID] = connection

	if m.logger != nil {
		m.logger.Info(context.Background()).
			Uint("business_id", businessID).
			Str("connection_id", connectionID).
			Interface("filter", filter).
			Msg("Nueva conexión SSE agregada")
	}

	return connectionID
}

// RemoveConnection remueve una conexión SSE por connectionID
func (m *EventManager) RemoveConnection(connectionID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if connection, exists := m.connections[connectionID]; exists {
		delete(m.connections, connectionID)
		if m.logger != nil {
			m.logger.Info(context.Background()).
				Uint("business_id", connection.BusinessID).
				Str("connection_id", connectionID).
				Msg("Conexión SSE removida")
		}
	}
}

// GetConnectionCount retorna el número de conexiones para un business_id
func (m *EventManager) GetConnectionCount(businessID uint) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	count := 0
	for _, conn := range m.connections {
		// Super usuario (businessID = 0) cuenta para todos, o si coincide el business_id
		if conn.BusinessID == 0 || conn.BusinessID == businessID {
			count++
		}
	}
	return count
}

// GetConnectionInfo retorna información sobre las conexiones activas de un business_id
func (m *EventManager) GetConnectionInfo(businessID uint) map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	connections := make([]map[string]interface{}, 0)
	for _, conn := range m.connections {
		if conn.BusinessID == 0 || conn.BusinessID == businessID {
			connections = append(connections, map[string]interface{}{
				"connection_id": conn.ConnectionID,
				"business_id":   conn.BusinessID,
				"is_super_user": conn.IsSuperUser(),
				"filter":        conn.Filter,
			})
		}
	}

	return map[string]interface{}{
		"business_id":      businessID,
		"active_connections": len(connections),
		"connections":      connections,
	}
}

