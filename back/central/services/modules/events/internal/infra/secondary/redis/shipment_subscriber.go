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

// ShipmentEventSubscriber consume eventos de envíos desde Redis Pub/Sub
type ShipmentEventSubscriber struct {
	redisClient redisclient.IRedis
	logger      log.ILogger
	channel     string
	pubsub      *goredis.PubSub
	eventChan   chan *domain.ShipmentEvent
	stopChan    chan struct{}
}

// NewShipmentEventSubscriber crea un nuevo suscriptor de eventos de envíos
func NewShipmentEventSubscriber(
	redisClient redisclient.IRedis,
	logger log.ILogger,
	channel string,
) *ShipmentEventSubscriber {
	return &ShipmentEventSubscriber{
		redisClient: redisClient,
		logger:      logger,
		channel:     channel,
		eventChan:   make(chan *domain.ShipmentEvent, 100),
		stopChan:    make(chan struct{}),
	}
}

// Start inicia el consumidor de eventos desde Redis
func (s *ShipmentEventSubscriber) Start(ctx context.Context) error {
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
func (s *ShipmentEventSubscriber) processMessages(ctx context.Context) {
	ch := s.pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Deserializar el mensaje
			var shipmentEvent domain.ShipmentEvent
			if err := json.Unmarshal([]byte(msg.Payload), &shipmentEvent); err != nil {
				s.logger.Error(ctx).
					Err(err).
					Str("payload", msg.Payload).
					Msg("Error deserializando evento de envío desde Redis")
				continue
			}

			// Validar el evento
			if !shipmentEvent.Type.IsValid() {
				s.logger.Warn(ctx).
					Str("event_type", string(shipmentEvent.Type)).
					Msg("Tipo de evento de envío inválido recibido")
				continue
			}

			// Enviar al canal de eventos
			select {
			case s.eventChan <- &shipmentEvent:
				s.logger.Info(ctx).
					Str("event_id", shipmentEvent.ID).
					Str("event_type", string(shipmentEvent.Type)).
					Uint("business_id", shipmentEvent.BusinessID).
					Msg("Evento de envío recibido desde Redis, enviando al consumer...")
			default:
				s.logger.Warn(ctx).
					Str("event_id", shipmentEvent.ID).
					Str("event_type", string(shipmentEvent.Type)).
					Msg("Canal de eventos de envío lleno, descartando evento")
			}

		case <-s.stopChan:
			s.logger.Info(ctx).Msg("Deteniendo suscriptor Redis de envíos")
			return
		case <-ctx.Done():
			s.logger.Info(ctx).Msg("Context cancelado, deteniendo suscriptor Redis de envíos")
			return
		}
	}
}

// GetEventChannel retorna el canal de eventos para consumo externo
func (s *ShipmentEventSubscriber) GetEventChannel() <-chan *domain.ShipmentEvent {
	return s.eventChan
}

// Stop detiene el suscriptor
func (s *ShipmentEventSubscriber) Stop() error {
	close(s.stopChan)
	if s.pubsub != nil {
		return s.pubsub.Close()
	}
	return nil
}
