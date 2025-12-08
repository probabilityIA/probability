package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/shared/log"
)

// OrderEventConsumer consume eventos de órdenes desde Redis y los procesa
type OrderEventConsumer struct {
	subscriber   *redis.OrderEventSubscriber
	eventManager domain.IEventPublisher
	configRepo   domain.INotificationConfigRepository
	logger       log.ILogger
}

// IOrderEventConsumer define la interfaz del consumidor
type IOrderEventConsumer interface {
	Start(ctx context.Context) error
	Stop() error
}

// NewOrderEventConsumer crea un nuevo consumidor de eventos de órdenes
func New(
	subscriber *redis.OrderEventSubscriber,
	eventManager domain.IEventPublisher,
	configRepo domain.INotificationConfigRepository,
	logger log.ILogger,
) IOrderEventConsumer {
	return &OrderEventConsumer{
		subscriber:   subscriber,
		eventManager: eventManager,
		configRepo:   configRepo,
		logger:       logger,
	}
}

// Start inicia el consumidor
func (c *OrderEventConsumer) Start(ctx context.Context) error {
	// Iniciar el suscriptor Redis
	if err := c.subscriber.Start(ctx); err != nil {
		return err
	}

	// Iniciar worker para procesar eventos
	go c.processEvents(ctx)

	c.logger.Info(ctx).Msg("Consumidor de eventos de órdenes iniciado")
	return nil
}

// processEvents procesa los eventos recibidos de Redis
func (c *OrderEventConsumer) processEvents(ctx context.Context) {
	eventChan := c.subscriber.GetEventChannel()

	for {
		select {
		case event := <-eventChan:
			if event == nil {
				continue
			}

			// Verificar si el evento debe ser notificado según la configuración
			if c.shouldNotifyEvent(ctx, event) {
				// Convertir evento de orden a evento genérico y publicar
				c.publishOrderEvent(ctx, event)
			} else {
				c.logger.Debug(ctx).
					Str("event_id", event.ID).
					Str("event_type", string(event.Type)).
					Str("order_id", event.OrderID).
					Msg("Evento filtrado por configuración de notificaciones")
			}

		case <-ctx.Done():
			c.logger.Info(ctx).Msg("Context cancelado, deteniendo procesador de eventos")
			return
		}
	}
}

// shouldNotifyEvent verifica si un evento debe ser notificado según la configuración
func (c *OrderEventConsumer) shouldNotifyEvent(ctx context.Context, event *domain.OrderEvent) bool {
	// Si no hay business_id, notificar siempre (eventos globales)
	if event.BusinessID == nil {
		return true
	}

	// Obtener configuración para este negocio y tipo de evento
	config, err := c.configRepo.GetByBusinessAndEventType(ctx, *event.BusinessID, string(event.Type))
	if err != nil {
		// Si no hay configuración, notificar por defecto
		c.logger.Debug(ctx).
			Err(err).
			Uint("business_id", *event.BusinessID).
			Str("event_type", string(event.Type)).
			Msg("No se encontró configuración, notificando por defecto")
		return true
	}

	// Si no existe configuración, notificar por defecto
	if config == nil {
		return true
	}

	// Verificar si está habilitado
	return config.Enabled
}

// publishOrderEvent publica un evento de orden al sistema de eventos
func (c *OrderEventConsumer) publishOrderEvent(ctx context.Context, orderEvent *domain.OrderEvent) {
	// Convertir OrderEvent a Event genérico para compatibilidad
	var integrationID int64
	if orderEvent.IntegrationID != nil {
		integrationID = int64(*orderEvent.IntegrationID)
	}

	var businessIDStr string
	if orderEvent.BusinessID != nil {
		businessIDStr = fmt.Sprintf("%d", *orderEvent.BusinessID)
	}

	// Agregar información adicional al metadata
	metadata := make(map[string]interface{})
	if orderEvent.Metadata != nil {
		for k, v := range orderEvent.Metadata {
			metadata[k] = v
		}
	}
	metadata["order_id"] = orderEvent.OrderID
	if orderEvent.BusinessID != nil {
		metadata["business_id"] = *orderEvent.BusinessID
	}
	if orderEvent.IntegrationID != nil {
		metadata["integration_id"] = *orderEvent.IntegrationID
	}

	// Convertir OrderEventData a map para Data
	dataMap := map[string]interface{}{
		"order_id":        orderEvent.OrderID,
		"order_number":    orderEvent.Data.OrderNumber,
		"internal_number": orderEvent.Data.InternalNumber,
		"external_id":     orderEvent.Data.ExternalID,
		"previous_status": orderEvent.Data.PreviousStatus,
		"current_status":  orderEvent.Data.CurrentStatus,
		"platform":        orderEvent.Data.Platform,
		"customer_email":  orderEvent.Data.CustomerEmail,
		"total_amount":    orderEvent.Data.TotalAmount,
		"currency":        orderEvent.Data.Currency,
	}
	if orderEvent.Data.Extra != nil {
		for k, v := range orderEvent.Data.Extra {
			dataMap[k] = v
		}
	}

	genericEvent := domain.Event{
		ID:            orderEvent.ID,
		Type:          domain.EventType(orderEvent.Type),
		IntegrationID: integrationID,
		BusinessID:    businessIDStr,
		Timestamp:     orderEvent.Timestamp,
		Data:          dataMap,
		Metadata:      metadata,
	}

	// Publicar el evento
	c.eventManager.PublishEvent(genericEvent)

	c.logger.Debug(ctx).
		Str("event_id", orderEvent.ID).
		Str("event_type", string(orderEvent.Type)).
		Str("order_id", orderEvent.OrderID).
		Msg("Evento de orden publicado al sistema de eventos")
}

// Stop detiene el consumidor
func (c *OrderEventConsumer) Stop() error {
	return c.subscriber.Stop()
}
