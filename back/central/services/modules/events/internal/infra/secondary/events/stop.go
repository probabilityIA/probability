package events

import (
	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
)

// Stop detiene el manager y limpia conexiones
func (m *EventManager) Stop() {
	close(m.stopChan)
	m.mutex.Lock()
	m.connections = make(map[string]*domain.SSEConnection)
	m.mutex.Unlock()
}
