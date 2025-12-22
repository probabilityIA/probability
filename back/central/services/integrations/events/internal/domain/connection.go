package domain

import (
	"net/http"
	"sync"
	"time"
)

// ───────────────────────────────────────────
//
//	SSE CONNECTION
//
// ───────────────────────────────────────────

// IntegrationSSEConnection representa una conexión SSE para eventos de integraciones
type IntegrationSSEConnection struct {
	ID           string
	BusinessID   uint
	Filter       *IntegrationSSEFilter
	Writer       http.ResponseWriter
	CreatedAt    time.Time
	LastEventSeq int64
	EventCount   int
}

// IntegrationSSEFilter define los filtros para una conexión SSE
type IntegrationSSEFilter struct {
	IntegrationID *uint
	EventTypes    []IntegrationEventType
	OrderIDs      []string
}

// ConnectionManager gestiona las conexiones SSE
type ConnectionManager struct {
	connections map[string]*IntegrationSSEConnection
	mutex       sync.RWMutex
	seqCounter  int64
}

// NewConnectionManager crea un nuevo gestor de conexiones
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*IntegrationSSEConnection),
		seqCounter:  0,
	}
}

// AddConnection agrega una nueva conexión
func (cm *ConnectionManager) AddConnection(conn *IntegrationSSEConnection) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connections[conn.ID] = conn
}

// RemoveConnection elimina una conexión
func (cm *ConnectionManager) RemoveConnection(connectionID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.connections, connectionID)
}

// GetConnection obtiene una conexión por ID
func (cm *ConnectionManager) GetConnection(connectionID string) (*IntegrationSSEConnection, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	conn, exists := cm.connections[connectionID]
	return conn, exists
}

// GetConnectionsByBusiness obtiene todas las conexiones de un business
func (cm *ConnectionManager) GetConnectionsByBusiness(businessID uint) []*IntegrationSSEConnection {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	var connections []*IntegrationSSEConnection
	for _, conn := range cm.connections {
		if conn.BusinessID == businessID {
			connections = append(connections, conn)
		}
	}
	return connections
}

// GetAllConnections obtiene todas las conexiones
func (cm *ConnectionManager) GetAllConnections() []*IntegrationSSEConnection {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	connections := make([]*IntegrationSSEConnection, 0, len(cm.connections))
	for _, conn := range cm.connections {
		connections = append(connections, conn)
	}
	return connections
}

// NextSeq incrementa y retorna el siguiente número de secuencia
func (cm *ConnectionManager) NextSeq() int64 {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.seqCounter++
	return cm.seqCounter
}

// GetConnectionCount retorna el número de conexiones
func (cm *ConnectionManager) GetConnectionCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.connections)
}
