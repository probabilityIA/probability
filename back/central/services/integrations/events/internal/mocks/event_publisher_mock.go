package mocks

import (
	"net/http"
	"sync"

	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
)

// IntegrationEventPublisherMock es el mock de domain.IIntegrationEventPublisher.
// Captura todos los eventos publicados para que los tests puedan inspeccionarlos.
type IntegrationEventPublisherMock struct {
	mu     sync.Mutex
	events []domain.IntegrationEvent

	// Funciones inyectables para personalizar el comportamiento en cada test.
	AddConnectionFn            func(businessID uint, filter *domain.IntegrationSSEFilter, conn http.ResponseWriter) string
	RemoveConnectionFn         func(connectionID string)
	PublishEventFn             func(event domain.IntegrationEvent)
	GetConnectionCountFn       func(businessID uint) int
	GetConnectionInfoFn        func(businessID uint) map[string]interface{}
	GetRecentEventsByBusinessFn func(businessID uint, sinceSeq int64) []domain.IntegrationEvent
	HasRecentEventsFn          func(businessID uint) bool
	StopFn                     func()
}

// AddConnection implementa IIntegrationEventPublisher.
func (m *IntegrationEventPublisherMock) AddConnection(businessID uint, filter *domain.IntegrationSSEFilter, conn http.ResponseWriter) string {
	if m.AddConnectionFn != nil {
		return m.AddConnectionFn(businessID, filter, conn)
	}
	return "mock-connection-id"
}

// RemoveConnection implementa IIntegrationEventPublisher.
func (m *IntegrationEventPublisherMock) RemoveConnection(connectionID string) {
	if m.RemoveConnectionFn != nil {
		m.RemoveConnectionFn(connectionID)
	}
}

// PublishEvent implementa IIntegrationEventPublisher.
// Por defecto acumula cada evento recibido para que los tests puedan inspeccionarlos.
func (m *IntegrationEventPublisherMock) PublishEvent(event domain.IntegrationEvent) {
	if m.PublishEventFn != nil {
		m.PublishEventFn(event)
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
}

// GetConnectionCount implementa IIntegrationEventPublisher.
func (m *IntegrationEventPublisherMock) GetConnectionCount(businessID uint) int {
	if m.GetConnectionCountFn != nil {
		return m.GetConnectionCountFn(businessID)
	}
	return 0
}

// GetConnectionInfo implementa IIntegrationEventPublisher.
func (m *IntegrationEventPublisherMock) GetConnectionInfo(businessID uint) map[string]interface{} {
	if m.GetConnectionInfoFn != nil {
		return m.GetConnectionInfoFn(businessID)
	}
	return map[string]interface{}{}
}

// GetRecentEventsByBusiness implementa IIntegrationEventPublisher.
func (m *IntegrationEventPublisherMock) GetRecentEventsByBusiness(businessID uint, sinceSeq int64) []domain.IntegrationEvent {
	if m.GetRecentEventsByBusinessFn != nil {
		return m.GetRecentEventsByBusinessFn(businessID, sinceSeq)
	}
	return nil
}

// HasRecentEvents implementa IIntegrationEventPublisher.
func (m *IntegrationEventPublisherMock) HasRecentEvents(businessID uint) bool {
	if m.HasRecentEventsFn != nil {
		return m.HasRecentEventsFn(businessID)
	}
	return false
}

// Stop implementa IIntegrationEventPublisher.
func (m *IntegrationEventPublisherMock) Stop() {
	if m.StopFn != nil {
		m.StopFn()
	}
}

// PublishedEvents retorna una copia segura de todos los eventos capturados.
func (m *IntegrationEventPublisherMock) PublishedEvents() []domain.IntegrationEvent {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]domain.IntegrationEvent, len(m.events))
	copy(result, m.events)
	return result
}

// Reset limpia los eventos acumulados (util para reutilizar el mock entre sub-tests).
func (m *IntegrationEventPublisherMock) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = nil
}
