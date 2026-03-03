package usecasecreateorder

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// publishOrderEvents es el punto de entrada ÚNICO para publicar todos los eventos
// después de crear una orden exitosamente.
// Orquesta la publicación a todos los canales necesarios:
//   - Integration sync (notifica al módulo de integraciones)
//   - RabbitMQ fanout (invoicing, inventory, score, whatsapp, events consumers)
//   - Score (cálculo directo)
func (uc *UseCaseCreateOrder) publishOrderEvents(ctx context.Context, order *entities.ProbabilityOrder, isManualOrder bool) {
	// 1. Notificar sincronización exitosa a integraciones (solo órdenes de integración)
	if !isManualOrder {
		uc.publishSyncOrderCreated(ctx, order)
	}

	// 2. Publicar evento de orden creada a RabbitMQ fanout
	// (invoicing, inventory, score, whatsapp, events consumers)
	uc.publishOrderCreatedEvent(ctx, order)

	// 3. Calcular score directamente
	uc.calculateOrderScore(ctx, order)
}

// ───────────────────────────────────────────
//
//	INTEGRATION SYNC EVENTS
//
// ───────────────────────────────────────────

// publishSyncOrderCreated notifica al módulo de integraciones que la orden se creó exitosamente
func (uc *UseCaseCreateOrder) publishSyncOrderCreated(ctx context.Context, order *entities.ProbabilityOrder) {
	if uc.integrationEventPublisher == nil {
		return
	}
	uc.integrationEventPublisher.PublishSyncOrderCreated(ctx, order.IntegrationID, order.BusinessID, map[string]interface{}{
		"order_id":       order.ID,
		"order_number":   order.OrderNumber,
		"external_id":    order.ExternalID,
		"platform":       order.Platform,
		"customer_email": order.CustomerEmail,
		"currency":       order.Currency,
		"status":         order.Status,
		"created_at":     order.CreatedAt,
		"total_amount":   order.TotalAmount,
		"synced_at":      time.Now(),
	})
}

// ───────────────────────────────────────────
//
//	RABBITMQ FANOUT EVENTS
//
// ───────────────────────────────────────────

// publishOrderCreatedEvent publica el evento de orden creada al exchange fanout de RabbitMQ
// (invoicing, inventory, score, whatsapp consumers)
func (uc *UseCaseCreateOrder) publishOrderCreatedEvent(_ context.Context, order *entities.ProbabilityOrder) {
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

	event := entities.NewOrderEvent(entities.OrderEventTypeCreated, order.ID, eventData)
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

// ───────────────────────────────────────────
//
//	SCORE
//
// ───────────────────────────────────────────

// calculateOrderScore calcula el score de la orden en background (goroutine).
func (uc *UseCaseCreateOrder) calculateOrderScore(_ context.Context, order *entities.ProbabilityOrder) {
	orderID := order.ID
	orderNumber := order.OrderNumber

	go func() {
		bgCtx := context.Background()

		if err := uc.scoreUseCase.CalculateAndUpdateOrderScore(bgCtx, orderID); err != nil {
			uc.logger.Error(bgCtx).
				Err(err).
				Str("order_id", orderID).
				Msg("Error al calcular score de la orden")
			return
		}

		uc.logger.Info(bgCtx).
			Str("order_id", orderID).
			Str("order_number", orderNumber).
			Msg("Score calculado exitosamente para la orden")
	}()
}

