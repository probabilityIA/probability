package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/shared/log"
)

// ShipmentEventConsumer consume eventos de envíos desde Redis y los publica al EventManager
type ShipmentEventConsumer struct {
	subscriber   *redis.ShipmentEventSubscriber
	eventManager domain.IEventPublisher
	logger       log.ILogger
}

// IShipmentEventConsumer define la interfaz del consumidor de eventos de envíos
type IShipmentEventConsumer interface {
	Start(ctx context.Context) error
	Stop() error
}

// NewShipmentEventConsumer crea un nuevo consumidor de eventos de envíos
func NewShipmentEventConsumer(
	subscriber *redis.ShipmentEventSubscriber,
	eventManager domain.IEventPublisher,
	logger log.ILogger,
) IShipmentEventConsumer {
	return &ShipmentEventConsumer{
		subscriber:   subscriber,
		eventManager: eventManager,
		logger:       logger,
	}
}

// Start inicia el consumidor
func (c *ShipmentEventConsumer) Start(ctx context.Context) error {
	// Iniciar el suscriptor Redis
	if err := c.subscriber.Start(ctx); err != nil {
		return err
	}

	// Iniciar worker para procesar eventos
	go c.processEvents(ctx)

	return nil
}

// processEvents procesa los eventos recibidos de Redis
func (c *ShipmentEventConsumer) processEvents(ctx context.Context) {
	eventChan := c.subscriber.GetEventChannel()

	for {
		select {
		case event := <-eventChan:
			if event == nil {
				continue
			}

			c.logger.Info(ctx).
				Str("event_id", event.ID).
				Str("event_type", string(event.Type)).
				Uint("business_id", event.BusinessID).
				Msg("Evento de envío aprobado para notificación, publicando...")

			c.publishShipmentEvent(ctx, event)

		case <-ctx.Done():
			c.logger.Info(ctx).Msg("Context cancelado, deteniendo procesador de eventos de envíos")
			return
		}
	}
}

// publishShipmentEvent convierte un ShipmentEvent a Event genérico y lo publica al EventManager
func (c *ShipmentEventConsumer) publishShipmentEvent(ctx context.Context, shipmentEvent *domain.ShipmentEvent) {
	businessIDStr := fmt.Sprintf("%d", shipmentEvent.BusinessID)

	// Construir metadata
	metadata := map[string]interface{}{
		"business_id": shipmentEvent.BusinessID,
		"event_type":  string(shipmentEvent.Type),
	}

	// Agregar campos relevantes del data al metadata para filtros
	if shipmentID, ok := shipmentEvent.Data["shipment_id"]; ok {
		metadata["shipment_id"] = shipmentID
	}
	if correlationID, ok := shipmentEvent.Data["correlation_id"]; ok {
		metadata["correlation_id"] = correlationID
	}

	genericEvent := domain.Event{
		ID:         shipmentEvent.ID,
		Type:       domain.EventType(shipmentEvent.Type),
		BusinessID: businessIDStr,
		Timestamp:  shipmentEvent.Timestamp,
		Data:       shipmentEvent.Data,
		Metadata:   metadata,
	}

	// Publicar el evento al EventManager (SSE broadcast)
	c.eventManager.PublishEvent(genericEvent)

	c.logger.Info(ctx).
		Str("event_id", shipmentEvent.ID).
		Str("event_type", string(shipmentEvent.Type)).
		Str("business_id", businessIDStr).
		Msg("Evento de envío publicado al EventManager (SSE)")
}

// Stop detiene el consumidor
func (c *ShipmentEventConsumer) Stop() error {
	return c.subscriber.Stop()
}
