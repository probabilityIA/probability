package redis

import (
	"context"
	"encoding/json"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// IntegrationEventSubscriber consume eventos de integración desde Redis Pub/Sub
type IntegrationEventSubscriber struct {
	redisClient redisclient.IRedis
	logger      log.ILogger
	channel     string
	pubsub      *goredis.PubSub
	eventChan   chan *domain.IntegrationEvent
	stopChan    chan struct{}
}

// NewIntegrationEventSubscriber crea un nuevo suscriptor de eventos de integración
func NewIntegrationEventSubscriber(
	redisClient redisclient.IRedis,
	logger log.ILogger,
	channel string,
) *IntegrationEventSubscriber {
	return &IntegrationEventSubscriber{
		redisClient: redisClient,
		logger:      logger,
		channel:     channel,
		eventChan:   make(chan *domain.IntegrationEvent, 100),
		stopChan:    make(chan struct{}),
	}
}

// Start inicia el consumidor de eventos desde Redis
func (s *IntegrationEventSubscriber) Start(ctx context.Context) error {
	client := s.redisClient.Client(ctx)
	if client == nil {
		return fmt.Errorf("redis client no disponible")
	}

	// Suscribirse al canal
	s.pubsub = client.Subscribe(ctx, s.channel)

	// Iniciar goroutine para procesar mensajes
	go s.processMessages(ctx)

	return nil
}

// processMessages procesa los mensajes recibidos de Redis
func (s *IntegrationEventSubscriber) processMessages(ctx context.Context) {
	ch := s.pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Deserializar el mensaje
			var event domain.IntegrationEvent
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				s.logger.Error(ctx).
					Err(err).
					Str("payload", msg.Payload).
					Msg("Error deserializando integration event desde Redis")
				continue
			}

			// Validar el tipo de evento
			if !event.Type.IsValid() {
				s.logger.Warn(ctx).
					Str("event_type", string(event.Type)).
					Str("event_id", event.ID).
					Msg("Tipo de integration event inválido recibido")
				continue
			}

			// Enviar al canal de eventos
			select {
			case s.eventChan <- &event:
				s.logger.Info(ctx).
					Str("event_id", event.ID).
					Str("event_type", string(event.Type)).
					Uint("integration_id", event.IntegrationID).
					Interface("business_id", event.BusinessID).
					Msg("📥 Integration event recibido desde Redis, enviando al consumer...")
			default:
				s.logger.Warn(ctx).
					Str("event_id", event.ID).
					Str("event_type", string(event.Type)).
					Msg("⚠️ Canal de integration events lleno, descartando evento")
			}

		case <-s.stopChan:
			s.logger.Info(ctx).Msg("Deteniendo suscriptor Redis de integration events")
			return
		case <-ctx.Done():
			s.logger.Info(ctx).Msg("Context cancelado, deteniendo suscriptor Redis de integration events")
			return
		}
	}
}

// GetEventChannel retorna el canal de eventos para consumo externo
func (s *IntegrationEventSubscriber) GetEventChannel() <-chan *domain.IntegrationEvent {
	return s.eventChan
}

// Stop detiene el suscriptor
func (s *IntegrationEventSubscriber) Stop() error {
	close(s.stopChan)
	if s.pubsub != nil {
		return s.pubsub.Close()
	}
	return nil
}
