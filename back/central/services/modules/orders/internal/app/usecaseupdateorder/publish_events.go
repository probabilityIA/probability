package usecaseupdateorder

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// publishUpdateEvents es el punto de entrada ÚNICO para publicar todos los eventos
// después de actualizar una orden exitosamente.
// Orquesta la publicación a todos los canales necesarios:
//   - Integration sync (notifica al módulo de integraciones)
//   - RabbitMQ fanout (invoicing, inventory, score, whatsapp, events consumers)
//   - Status changed (si cambió el estado)
func (uc *UseCaseUpdateOrder) publishUpdateEvents(ctx context.Context, order *entities.ProbabilityOrder, previousStatus string, isManualOrder bool) {
	// 1. Notificar sincronización exitosa a integraciones (solo órdenes de integración)
	if !isManualOrder && order.IntegrationID > 0 {
		uc.publishSyncOrderUpdated(ctx, order)
	}

	// 2. Publicar evento de orden actualizada a RabbitMQ fanout
	// (invoicing, inventory, score, whatsapp, events consumers)
	uc.publishOrderUpdatedEvent(ctx, order)

	// 3. Si cambió el estado, publicar evento de cambio de estado
	if previousStatus != order.Status {
		uc.publishOrderStatusChangedEvent(ctx, order, previousStatus)
	}

	// Score recalculation handled by probability module via QueueOrdersToScore
}

// ───────────────────────────────────────────
//
//	INTEGRATION SYNC EVENTS
//
// ───────────────────────────────────────────

// publishSyncOrderUpdated notifica al módulo de integraciones que la orden se actualizó
func (uc *UseCaseUpdateOrder) publishSyncOrderUpdated(ctx context.Context, order *entities.ProbabilityOrder) {
	if uc.integrationEventPublisher == nil {
		return
	}
	uc.integrationEventPublisher.PublishSyncOrderUpdated(ctx, order.IntegrationID, order.BusinessID, map[string]interface{}{
		"order_id":       order.ID,
		"order_number":   order.OrderNumber,
		"external_id":    order.ExternalID,
		"platform":       order.Platform,
		"customer_email": order.CustomerEmail,
		"currency":       order.Currency,
		"status":         order.Status,
		"created_at":     order.OccurredAt,
		"total_amount":   order.TotalAmount,
		"updated_at":     time.Now(),
	})
}

// ───────────────────────────────────────────
//
//	RABBITMQ FANOUT EVENTS
//
// ───────────────────────────────────────────

// publishOrderUpdatedEvent publica el evento de orden actualizada al exchange fanout de RabbitMQ
func (uc *UseCaseUpdateOrder) publishOrderUpdatedEvent(ctx context.Context, order *entities.ProbabilityOrder) {
	if uc.rabbitEventPublisher == nil {
		return
	}

	eventData := entities.OrderEventData{
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,
		ExternalID:     order.ExternalID,
		CurrentStatus:  order.Status,
		CustomerEmail:  order.CustomerEmail,
		TotalAmount:    &order.TotalAmount,
		Currency:       order.Currency,
		Platform:       order.Platform,
	}

	event := entities.NewOrderEvent(entities.OrderEventTypeUpdated, order.ID, eventData)
	event.BusinessID = order.BusinessID
	if order.IntegrationID > 0 {
		integrationID := order.IntegrationID
		event.IntegrationID = &integrationID
	}

	go func() {
		if err := uc.rabbitEventPublisher.PublishOrderEvent(context.Background(), event, order); err != nil {
			uc.logger.Error(context.Background()).
				Err(err).
				Str("event_type", string(event.Type)).
				Str("order_id", event.OrderID).
				Msg("Error al publicar evento a RabbitMQ")
		}
	}()
}

// publishOrderStatusChangedEvent publica el evento de cambio de estado al exchange fanout
func (uc *UseCaseUpdateOrder) publishOrderStatusChangedEvent(_ context.Context, order *entities.ProbabilityOrder, previousStatus string) {
	if uc.rabbitEventPublisher == nil {
		return
	}

	eventData := entities.OrderEventData{
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,
		ExternalID:     order.ExternalID,
		PreviousStatus: previousStatus,
		CurrentStatus:  order.Status,
		CustomerEmail:  order.CustomerEmail,
		TotalAmount:    &order.TotalAmount,
		Currency:       order.Currency,
		Platform:       order.Platform,
	}

	event := entities.NewOrderEvent(entities.OrderEventTypeStatusChanged, order.ID, eventData)
	event.BusinessID = order.BusinessID
	if order.IntegrationID > 0 {
		integrationID := order.IntegrationID
		event.IntegrationID = &integrationID
	}

	go func() {
		if err := uc.rabbitEventPublisher.PublishOrderEvent(context.Background(), event, order); err != nil {
			uc.logger.Error(context.Background()).
				Err(err).
				Str("event_type", string(event.Type)).
				Str("order_id", event.OrderID).
				Msg("Error al publicar evento de cambio de estado a RabbitMQ")
		}
	}()
}


