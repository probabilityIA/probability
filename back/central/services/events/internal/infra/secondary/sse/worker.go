package sse

import (
	"context"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

// startEventWorker procesa eventos del channel y los envía a las conexiones
func (m *EventManager) startEventWorker() {
	for {
		select {
		case event := <-m.eventChan:
			businessID := event.BusinessID

			if event.Metadata == nil {
				event.Metadata = make(map[string]interface{})
			}
			m.mutex.Lock()
			if _, ok := m.recentEvents[businessID]; !ok {
				m.recentEvents[businessID] = make([]entities.Event, 0)
			}
			seq := len(m.recentEvents[businessID]) + 1
			event.Metadata["sse_seq"] = seq
			m.mutex.Unlock()

			m.broadcastToBusinesses(event)
			m.appendRecentEvent(businessID, event)

		case <-m.stopChan:
			if m.logger != nil {
				m.logger.Info(context.Background()).
					Msg("Event worker SSE detenido")
			}
			return
		}
	}
}
