package events

import (
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
)

// GetRecentEventsByBusiness obtiene eventos recientes para un business
func (m *IntegrationEventManager) GetRecentEventsByBusiness(businessID uint, sinceSeq int64) []domain.IntegrationEvent {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	events, exists := m.recentEvents[businessID]
	if !exists {
		return []domain.IntegrationEvent{}
	}

	// Filtrar por secuencia si se especifica
	if sinceSeq > 0 {
		filtered := make([]domain.IntegrationEvent, 0)
		for _, event := range events {
			// Los eventos no tienen seq directamente, pero podemos usar timestamp
			// Por ahora, retornar todos los eventos recientes
			filtered = append(filtered, event)
		}
		return filtered
	}

	return events
}

// HasRecentEvents verifica si hay eventos recientes para un business
func (m *IntegrationEventManager) HasRecentEvents(businessID uint) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	events, exists := m.recentEvents[businessID]
	return exists && len(events) > 0
}

// addToRecentEvents agrega un evento al caché de eventos recientes
func (m *IntegrationEventManager) addToRecentEvents(event domain.IntegrationEvent) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var businessID uint
	if event.BusinessID != nil {
		businessID = *event.BusinessID
	}

	events, exists := m.recentEvents[businessID]
	if !exists {
		events = make([]domain.IntegrationEvent, 0, m.maxRecent)
	}

	// Agregar al inicio
	events = append([]domain.IntegrationEvent{event}, events...)

	// Limitar tamaño del caché
	if len(events) > m.maxRecent {
		events = events[:m.maxRecent]
	}

	m.recentEvents[businessID] = events
}

// updateStats actualiza las estadísticas de eventos
func (m *IntegrationEventManager) updateStats(event domain.IntegrationEvent) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var businessID uint
	if event.BusinessID != nil {
		businessID = *event.BusinessID
	}

	// Incrementar contador total
	m.eventCount[businessID]++

	// Incrementar contador por tipo
	if m.eventTypeCount[businessID] == nil {
		m.eventTypeCount[businessID] = make(map[domain.IntegrationEventType]int)
	}
	m.eventTypeCount[businessID][event.Type]++
}

// Stop detiene el manager de eventos
func (m *IntegrationEventManager) Stop() {
	close(m.stopChan)
	if m.logger != nil {
		m.logger.Info().Msg("Integration event manager stopped")
	}
}
