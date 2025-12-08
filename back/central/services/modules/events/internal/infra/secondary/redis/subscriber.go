package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// OrderEventSubscriber consume eventos de órdenes desde Redis Pub/Sub
type OrderEventSubscriber struct {
	redisClient redisclient.IRedis
	logger      log.ILogger
	channel     string
	pubsub      *redis.PubSub
	eventChan   chan *domain.OrderEvent
	stopChan    chan struct{}
}

// NewOrderEventSubscriber crea un nuevo suscriptor de eventos de órdenes
func New(
	redisClient redisclient.IRedis,
	logger log.ILogger,
	channel string,
) *OrderEventSubscriber {
	return &OrderEventSubscriber{
		redisClient: redisClient,
		logger:      logger,
		channel:     channel,
		eventChan:   make(chan *domain.OrderEvent, 100),
		stopChan:    make(chan struct{}),
	}
}

// Start inicia el consumidor de eventos desde Redis
func (s *OrderEventSubscriber) Start(ctx context.Context) error {
	client := s.redisClient.Client(ctx)
	if client == nil {
		return fmt.Errorf("redis client no disponible")
	}

	// Suscribirse al canal
	s.pubsub = client.Subscribe(ctx, s.channel)

	s.logger.Info(ctx).
		Str("channel", s.channel).
		Msg("Suscriptor Redis iniciado para eventos de órdenes")

	// Iniciar goroutine para procesar mensajes
	go s.processMessages(ctx)

	return nil
}

// processMessages procesa los mensajes recibidos de Redis
func (s *OrderEventSubscriber) processMessages(ctx context.Context) {
	ch := s.pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Deserializar el mensaje
			var orderEvent domain.OrderEvent
			if err := json.Unmarshal([]byte(msg.Payload), &orderEvent); err != nil {
				s.logger.Error(ctx).
					Err(err).
					Str("payload", msg.Payload).
					Msg("Error deserializando evento de orden desde Redis")
				continue
			}

			// Validar el evento
			if !orderEvent.Type.IsValid() {
				s.logger.Warn(ctx).
					Str("event_type", string(orderEvent.Type)).
					Str("order_id", orderEvent.OrderID).
					Msg("Tipo de evento de orden inválido recibido")
				continue
			}

			// Enviar al canal de eventos
			select {
			case s.eventChan <- &orderEvent:
				s.logger.Debug(ctx).
					Str("event_id", orderEvent.ID).
					Str("event_type", string(orderEvent.Type)).
					Str("order_id", orderEvent.OrderID).
					Msg("Evento de orden recibido desde Redis")
			default:
				s.logger.Warn(ctx).
					Str("event_id", orderEvent.ID).
					Msg("Canal de eventos lleno, descartando evento")
			}

		case <-s.stopChan:
			s.logger.Info(ctx).Msg("Deteniendo suscriptor Redis")
			return
		case <-ctx.Done():
			s.logger.Info(ctx).Msg("Context cancelado, deteniendo suscriptor Redis")
			return
		}
	}
}

// GetEventChannel retorna el canal de eventos para consumo externo
func (s *OrderEventSubscriber) GetEventChannel() <-chan *domain.OrderEvent {
	return s.eventChan
}

// Stop detiene el suscriptor
func (s *OrderEventSubscriber) Stop() error {
	close(s.stopChan)
	if s.pubsub != nil {
		return s.pubsub.Close()
	}
	return nil
}
