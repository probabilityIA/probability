package events

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
)

// eventServiceInstance es la instancia global del servicio de eventos
var eventServiceInstance domain.IIntegrationEventService

// SetEventService establece la instancia global del servicio de eventos
func SetEventService(service domain.IIntegrationEventService) {
	eventServiceInstance = service
}

// GetEventService retorna la instancia global del servicio de eventos
func GetEventService() domain.IIntegrationEventService {
	return eventServiceInstance
}

// PublishSyncOrderCreated publica un evento de orden creada (helper global)
func PublishSyncOrderCreated(ctx context.Context, integrationID uint, businessID *uint, orderID, orderNumber, externalID, platform, customerEmail, currency, status string, createdAt time.Time, totalAmount *float64) {
	if eventServiceInstance == nil {
		fmt.Printf("[PublishSyncOrderCreated] eventServiceInstance es nil, no se puede publicar evento\n")
		return
	}

	var syncedAt time.Time
	if val := ctx.Value("synced_at"); val != nil {
		if t, ok := val.(time.Time); ok {
			syncedAt = t
		}
	}
	if syncedAt.IsZero() {
		syncedAt = time.Now()
	}

	data := domain.SyncOrderCreatedEvent{
		OrderID:       orderID,
		OrderNumber:   orderNumber,
		ExternalID:    externalID,
		Platform:      platform,
		CustomerEmail: customerEmail,
		TotalAmount:   totalAmount,
		Currency:      currency,
		Status:        status,
		CreatedAt:     createdAt,
		SyncedAt:      syncedAt,
	}
	fmt.Printf("[PublishSyncOrderCreated] Publicando evento para integración %d, orden %s (%s)\n", integrationID, orderID, orderNumber)
	if err := eventServiceInstance.PublishSyncOrderCreated(ctx, integrationID, businessID, data); err != nil {
		fmt.Printf("[PublishSyncOrderCreated] Error al publicar evento: %v\n", err)
	} else {
		fmt.Printf("[PublishSyncOrderCreated] Evento publicado exitosamente\n")
	}
}

// PublishSyncOrderUpdated publica un evento de orden actualizada (helper global)
func PublishSyncOrderUpdated(ctx context.Context, integrationID uint, businessID *uint, orderID, orderNumber, externalID, platform, customerEmail, currency, status string, createdAt time.Time, totalAmount *float64) {
	if eventServiceInstance == nil {
		fmt.Printf("[PublishSyncOrderUpdated] eventServiceInstance es nil, no se puede publicar evento\n")
		return
	}

	data := domain.SyncOrderUpdatedEvent{
		OrderID:       orderID,
		OrderNumber:   orderNumber,
		ExternalID:    externalID,
		Platform:      platform,
		CustomerEmail: customerEmail,
		TotalAmount:   totalAmount,
		Currency:      currency,
		Status:        status,
		CreatedAt:     createdAt,
		UpdatedAt:     time.Now(),
	}
	fmt.Printf("[PublishSyncOrderUpdated] Publicando evento para integración %d, orden %s (%s)\n", integrationID, orderID, orderNumber)
	if err := eventServiceInstance.PublishSyncOrderUpdated(ctx, integrationID, businessID, data); err != nil {
		fmt.Printf("[PublishSyncOrderUpdated] Error al publicar evento: %v\n", err)
	} else {
		fmt.Printf("[PublishSyncOrderUpdated] Evento publicado exitosamente\n")
	}
}

// PublishSyncOrderRejected publica un evento de orden rechazada (helper global)
func PublishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, orderID, orderNumber, externalID, platform, reason, errorMsg string) {
	if eventServiceInstance == nil {
		return
	}

	data := domain.SyncOrderRejectedEvent{
		OrderID:     orderID,
		OrderNumber: orderNumber,
		ExternalID:  externalID,
		Platform:    platform,
		Reason:      reason,
		Error:       errorMsg,
		RejectedAt:  time.Now(),
	}
	eventServiceInstance.PublishSyncOrderRejected(ctx, integrationID, businessID, data)
}

// PublishSyncStarted publica un evento de inicio de sincronización (helper global)
func PublishSyncStarted(ctx context.Context, integrationID uint, businessID *uint, integrationType string, createdAtMin, createdAtMax *time.Time, status, financialStatus, fulfillmentStatus string) {
	if eventServiceInstance == nil {
		return
	}

	data := domain.SyncStartedEvent{
		IntegrationID:   integrationID,
		IntegrationType: integrationType,
		Params: domain.SyncParams{
			CreatedAtMin:      createdAtMin,
			CreatedAtMax:      createdAtMax,
			Status:            status,
			FinancialStatus:   financialStatus,
			FulfillmentStatus: fulfillmentStatus,
		},
		StartedAt: time.Now(),
	}
	eventServiceInstance.PublishSyncStarted(ctx, integrationID, businessID, data)
}

// PublishSyncCompleted publica un evento de sincronización completada (helper global)
func PublishSyncCompleted(ctx context.Context, integrationID uint, businessID *uint, integrationType string, totalOrders, createdOrders, updatedOrders, rejectedOrders int, duration time.Duration) {
	if eventServiceInstance == nil {
		return
	}

	data := domain.SyncCompletedEvent{
		IntegrationID:   integrationID,
		IntegrationType: integrationType,
		TotalOrders:     totalOrders,
		CreatedOrders:   createdOrders,
		UpdatedOrders:   updatedOrders,
		RejectedOrders:  rejectedOrders,
		Duration:        duration,
		CompletedAt:     time.Now(),
	}
	eventServiceInstance.PublishSyncCompleted(ctx, integrationID, businessID, data)
}

// PublishSyncFailed publica un evento de sincronización fallida (helper global)
func PublishSyncFailed(ctx context.Context, integrationID uint, businessID *uint, integrationType, errorMsg string) {
	if eventServiceInstance == nil {
		return
	}

	data := domain.SyncFailedEvent{
		IntegrationID:   integrationID,
		IntegrationType: integrationType,
		Error:           errorMsg,
		FailedAt:        time.Now(),
	}
	eventServiceInstance.PublishSyncFailed(ctx, integrationID, businessID, data)
}
