package events

import (
	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
)

// startEventWorker procesa eventos del channel y los envía a las conexiones
func (m *EventManager) startEventWorker() {
	for {
		select {
		case event := <-m.eventChan:
			// Extraer business_id del evento
			businessID := m.extractBusinessID(event)

			// Agregar secuencia al evento
			if event.Metadata == nil {
				event.Metadata = make(map[string]interface{})
			}
			m.mutex.Lock()
			// Usar business_id para el secuenciador
			if _, ok := m.recentEvents[businessID]; !ok {
				m.recentEvents[businessID] = make([]domain.Event, 0)
			}
			seq := len(m.recentEvents[businessID]) + 1
			event.Metadata["sse_seq"] = seq
			m.mutex.Unlock()

			// Broadcast a todas las conexiones que coincidan (por business_id y filtros)
			m.broadcastToBusinesses(event)

			// Agregar al caché de eventos recientes
			m.appendRecentEvent(businessID, event)

		case <-m.stopChan:
			return
		}
	}
}
