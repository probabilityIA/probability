package events

import (
	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
)

// appendRecentEvent agrega un evento al caché circular por business_id
func (m *EventManager) appendRecentEvent(businessID uint, event domain.Event) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.recentEvents[businessID]; !ok {
		m.recentEvents[businessID] = make([]domain.Event, 0)
	}

	list := m.recentEvents[businessID]
	list = append(list, event)
	if len(list) > m.maxRecent {
		list = list[len(list)-m.maxRecent:]
	}
	m.recentEvents[businessID] = list
}

// GetRecentEventsByBusiness retorna eventos con sse_seq mayor a sinceSeq para un business_id
func (m *EventManager) GetRecentEventsByBusiness(businessID uint, sinceSeq int64) []domain.Event {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	list := m.recentEvents[businessID]
	if len(list) == 0 {
		return nil
	}

	var out []domain.Event
	for _, ev := range list {
		if seqVal, ok := ev.Metadata["sse_seq"]; ok {
			switch v := seqVal.(type) {
			case int:
				if int64(v) <= sinceSeq {
					continue
				}
			case int64:
				if v <= sinceSeq {
					continue
				}
			case float64:
				if int64(v) <= sinceSeq {
					continue
				}
			}
		}
		out = append(out, ev)
	}
	return out
}

// HasRecentEvents indica si hay eventos recientes en caché para el business_id
func (m *EventManager) HasRecentEvents(businessID uint) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	list := m.recentEvents[businessID]
	return len(list) > 0
}
