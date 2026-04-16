package usecaseupdatestatus

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// publishStatusChangeEvents publica eventos después de un cambio de estado exitoso
func (uc *UseCaseUpdateStatus) publishStatusChangeEvents(ctx context.Context, order *entities.ProbabilityOrder, previousStatus string) {
	if uc.rabbitEventPublisher == nil {
		return
	}

	// 1. Publicar evento de orden actualizada
	uc.publishOrderUpdatedEvent(ctx, order)

	// 2. Publicar evento de cambio de estado
	uc.publishOrderStatusChangedEvent(ctx, order, previousStatus)
}

// publishOrderUpdatedEvent publica el evento de orden actualizada al exchange fanout
func (uc *UseCaseUpdateStatus) publishOrderUpdatedEvent(ctx context.Context, order *entities.ProbabilityOrder) {
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
				Msg("Error al publicar evento de actualización a RabbitMQ")
		}
	}()
}

// publishOrderStatusChangedEvent publica el evento de cambio de estado al exchange fanout
func (uc *UseCaseUpdateStatus) publishOrderStatusChangedEvent(_ context.Context, order *entities.ProbabilityOrder, previousStatus string) {
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
