package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
)

// IntegrationEventService implementa IIntegrationEventService
type IntegrationEventService struct {
	eventPublisher domain.IIntegrationEventPublisher
}

// NewIntegrationEventService crea un nuevo servicio de eventos de integraciones
func NewIntegrationEventService(eventPublisher domain.IIntegrationEventPublisher) domain.IIntegrationEventService {
	return &IntegrationEventService{
		eventPublisher: eventPublisher,
	}
}

// PublishSyncOrderCreated publica un evento de orden creada exitosamente
func (s *IntegrationEventService) PublishSyncOrderCreated(ctx context.Context, integrationID uint, businessID *uint, data domain.SyncOrderCreatedEvent) error {
	event := domain.IntegrationEvent{
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncOrderCreated,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
		Metadata: map[string]interface{}{
			"order_id":     data.OrderID,
			"order_number": data.OrderNumber,
			"external_id":  data.ExternalID,
			"platform":     data.Platform,
			"status":       data.Status,
			"created_at":   data.CreatedAt.Format(time.RFC3339),
		},
	}

	s.eventPublisher.PublishEvent(event)
	return nil
}

// PublishSyncOrderUpdated publica un evento de orden actualizada exitosamente
func (s *IntegrationEventService) PublishSyncOrderUpdated(ctx context.Context, integrationID uint, businessID *uint, data domain.SyncOrderUpdatedEvent) error {
	event := domain.IntegrationEvent{
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncOrderUpdated,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
		Metadata: map[string]interface{}{
			"order_id":     data.OrderID,
			"order_number": data.OrderNumber,
			"external_id":  data.ExternalID,
			"platform":     data.Platform,
			"status":       data.Status,
			"created_at":   data.CreatedAt.Format(time.RFC3339),
		},
	}

	fmt.Printf("[IntegrationEventService] Publicando evento integration.sync.order.updated - IntegrationID: %d, BusinessID: %v, OrderID: %s, OrderNumber: %s\n",
		integrationID, businessID, data.OrderID, data.OrderNumber)
	s.eventPublisher.PublishEvent(event)
	fmt.Printf("[IntegrationEventService] Evento integration.sync.order.updated publicado exitosamente\n")
	return nil
}

// PublishSyncOrderRejected publica un evento de orden rechazada
func (s *IntegrationEventService) PublishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, data domain.SyncOrderRejectedEvent) error {
	event := domain.IntegrationEvent{
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncOrderRejected,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
		Metadata: map[string]interface{}{
			"order_id":     data.OrderID,
			"order_number": data.OrderNumber,
			"external_id":  data.ExternalID,
			"platform":     data.Platform,
			"reason":       data.Reason,
			"error":        data.Error,
		},
	}

	s.eventPublisher.PublishEvent(event)
	return nil
}

// PublishSyncStarted publica un evento de inicio de sincronización
func (s *IntegrationEventService) PublishSyncStarted(ctx context.Context, integrationID uint, businessID *uint, data domain.SyncStartedEvent) error {
	event := domain.IntegrationEvent{
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncStarted,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
		Metadata: map[string]interface{}{
			"integration_type": data.IntegrationType,
		},
	}

	s.eventPublisher.PublishEvent(event)
	return nil
}

// PublishSyncCompleted publica un evento de sincronización completada
func (s *IntegrationEventService) PublishSyncCompleted(ctx context.Context, integrationID uint, businessID *uint, data domain.SyncCompletedEvent) error {
	event := domain.IntegrationEvent{
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncCompleted,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
		Metadata: map[string]interface{}{
			"integration_type": data.IntegrationType,
			"total_orders":     data.TotalOrders,
			"created_orders":   data.CreatedOrders,
			"updated_orders":   data.UpdatedOrders,
			"rejected_orders":  data.RejectedOrders,
		},
	}

	s.eventPublisher.PublishEvent(event)
	return nil
}

// PublishSyncFailed publica un evento de sincronización fallida
func (s *IntegrationEventService) PublishSyncFailed(ctx context.Context, integrationID uint, businessID *uint, data domain.SyncFailedEvent) error {
	event := domain.IntegrationEvent{
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncFailed,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
		Metadata: map[string]interface{}{
			"integration_type": data.IntegrationType,
			"error":            data.Error,
		},
	}

	s.eventPublisher.PublishEvent(event)
	return nil
}

// SetSyncState guarda el estado de una sincronización
// Nota: Implementación stub - el estado se determina desde eventos recientes
func (s *IntegrationEventService) SetSyncState(ctx context.Context, integrationID uint, state domain.SyncState) error {
	// Por ahora, el estado se determina desde los eventos recientes en memoria
	// Esta implementación puede extenderse en el futuro para usar Redis
	return nil
}

// GetSyncState obtiene el estado de una sincronización
// Nota: Implementación stub - el estado se determina desde eventos recientes
func (s *IntegrationEventService) GetSyncState(ctx context.Context, integrationID uint) (*domain.SyncState, error) {
	// Por ahora, el estado se determina desde los eventos recientes en memoria
	// Esta implementación puede extenderse en el futuro para usar Redis
	return nil, nil
}

// DeleteSyncState elimina el estado de una sincronización
// Nota: Implementación stub - el estado se determina desde eventos recientes
func (s *IntegrationEventService) DeleteSyncState(ctx context.Context, integrationID uint) error {
	// Por ahora, el estado se determina desde los eventos recientes en memoria
	// Esta implementación puede extenderse en el futuro para usar Redis
	return nil
}

// IncrementSyncCounter incrementa un contador de sincronización
// Nota: Implementación stub - los contadores se rastrean desde eventos
func (s *IntegrationEventService) IncrementSyncCounter(ctx context.Context, integrationID uint, counterType string) error {
	// Por ahora, los contadores se rastrean desde los eventos recientes en memoria
	// Esta implementación puede extenderse en el futuro para usar Redis
	return nil
}

// generateEventID genera un ID único para un evento
func generateEventID() string {
	return fmt.Sprintf("%d-%s", time.Now().Unix(), fmt.Sprintf("%x", time.Now().UnixNano()%1000000))
}
