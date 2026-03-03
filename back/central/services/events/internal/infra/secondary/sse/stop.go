package sse

import (
	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

// Stop detiene el manager y limpia conexiones
func (m *EventManager) Stop() {
	close(m.stopChan)
	m.mutex.Lock()
	m.connections = make(map[string]*entities.SSEConnection)
	m.mutex.Unlock()
}
