package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IntegrationEventConsumer consume eventos de integración desde Redis y los publica al EventManager
type IntegrationEventConsumer struct {
	subscriber   *redis.IntegrationEventSubscriber
	eventManager domain.IEventPublisher
	configRepo   domain.INotificationConfigRepository
	logger       log.ILogger
}

// IIntegrationEventConsumer define la interfaz del consumidor de integration events
type IIntegrationEventConsumer interface {
	Start(ctx context.Context) error
	Stop() error
}

// NewIntegrationEventConsumer crea un nuevo consumidor de eventos de integración
func NewIntegrationEventConsumer(
	subscriber *redis.IntegrationEventSubscriber,
	eventManager domain.IEventPublisher,
	configRepo domain.INotificationConfigRepository,
	logger log.ILogger,
) IIntegrationEventConsumer {
	return &IntegrationEventConsumer{
		subscriber:   subscriber,
		eventManager: eventManager,
		configRepo:   configRepo,
		logger:       logger,
	}
}

// Start inicia el consumidor
func (c *IntegrationEventConsumer) Start(ctx context.Context) error {
	if err := c.subscriber.Start(ctx); err != nil {
		return err
	}

	go c.processEvents(ctx)

	return nil
}

// processEvents procesa los eventos recibidos de Redis
func (c *IntegrationEventConsumer) processEvents(ctx context.Context) {
	eventChan := c.subscriber.GetEventChannel()

	for {
		select {
		case event := <-eventChan:
			if event == nil {
				continue
			}

			// Verificar si el evento debe ser notificado según la configuración
			if c.shouldNotifyEvent(ctx, event) {
				c.logger.Info(ctx).
					Str("event_id", event.ID).
					Str("event_type", string(event.Type)).
					Uint("integration_id", event.IntegrationID).
					Interface("business_id", event.BusinessID).
					Msg("Integration event aprobado para notificación, publicando...")

				c.publishIntegrationEvent(ctx, event)
			} else {
				c.logger.Debug(ctx).
					Str("event_id", event.ID).
					Str("event_type", string(event.Type)).
					Msg("Integration event filtrado por configuración de notificaciones")
			}

		case <-ctx.Done():
			c.logger.Info(ctx).Msg("Context cancelado, deteniendo procesador de integration events")
			return
		}
	}
}

// shouldNotifyEvent verifica si un evento debe ser notificado según la configuración
func (c *IntegrationEventConsumer) shouldNotifyEvent(ctx context.Context, event *domain.IntegrationEvent) bool {
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
			Msg("No se encontró configuración de notificación, notificando por defecto")
		return true
	}

	if config == nil {
		return true
	}

	return config.Enabled
}

// publishIntegrationEvent convierte un IntegrationEvent a Event genérico y lo publica al EventManager
func (c *IntegrationEventConsumer) publishIntegrationEvent(ctx context.Context, event *domain.IntegrationEvent) {
	var integrationID int64 = int64(event.IntegrationID)

	var businessIDStr string
	if event.BusinessID != nil {
		businessIDStr = fmt.Sprintf("%d", *event.BusinessID)
	}

	// Construir metadata
	metadata := make(map[string]interface{})
	if event.Metadata != nil {
		for k, v := range event.Metadata {
			metadata[k] = v
		}
	}
	metadata["integration_id"] = event.IntegrationID
	if event.BusinessID != nil {
		metadata["business_id"] = *event.BusinessID
	}

	// Usar event.Data como Data del evento genérico
	var dataMap interface{} = event.Data
	if event.Data == nil {
		dataMap = map[string]interface{}{}
	}

	genericEvent := domain.Event{
		ID:            event.ID,
		Type:          domain.EventType(event.Type),
		IntegrationID: integrationID,
		BusinessID:    businessIDStr,
		Timestamp:     event.Timestamp,
		Data:          dataMap,
		Metadata:      metadata,
	}

	c.eventManager.PublishEvent(genericEvent)

	c.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", string(event.Type)).
		Uint("integration_id", event.IntegrationID).
		Str("business_id", businessIDStr).
		Msg("Integration event publicado al EventManager (SSE)")
}

// Stop detiene el consumidor
func (c *IntegrationEventConsumer) Stop() error {
	return c.subscriber.Stop()
}
