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

// InvoiceEventSubscriber consume eventos de facturación desde Redis Pub/Sub
type InvoiceEventSubscriber struct {
	redisClient redisclient.IRedis
	logger      log.ILogger
	channel     string
	pubsub      *goredis.PubSub
	eventChan   chan *domain.InvoiceEvent
	stopChan    chan struct{}
}

// NewInvoiceEventSubscriber crea un nuevo suscriptor de eventos de facturación
func NewInvoiceEventSubscriber(
	redisClient redisclient.IRedis,
	logger log.ILogger,
	channel string,
) *InvoiceEventSubscriber {
	return &InvoiceEventSubscriber{
		redisClient: redisClient,
		logger:      logger,
		channel:     channel,
		eventChan:   make(chan *domain.InvoiceEvent, 100),
		stopChan:    make(chan struct{}),
	}
}

// Start inicia el consumidor de eventos desde Redis
func (s *InvoiceEventSubscriber) Start(ctx context.Context) error {
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
func (s *InvoiceEventSubscriber) processMessages(ctx context.Context) {
	ch := s.pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Deserializar el mensaje
			var invoiceEvent domain.InvoiceEvent
			if err := json.Unmarshal([]byte(msg.Payload), &invoiceEvent); err != nil {
				s.logger.Error(ctx).
					Err(err).
					Str("payload", msg.Payload).
					Msg("Error deserializando evento de facturación desde Redis")
				continue
			}

			// Validar el evento
			if !invoiceEvent.Type.IsValid() {
				s.logger.Warn(ctx).
					Str("event_type", string(invoiceEvent.Type)).
					Msg("Tipo de evento de facturación inválido recibido")
				continue
			}

			// Enviar al canal de eventos
			select {
			case s.eventChan <- &invoiceEvent:
				s.logger.Info(ctx).
					Str("event_id", invoiceEvent.ID).
					Str("event_type", string(invoiceEvent.Type)).
					Uint("business_id", invoiceEvent.BusinessID).
					Msg("Evento de facturación recibido desde Redis, enviando al consumer...")
			default:
				s.logger.Warn(ctx).
					Str("event_id", invoiceEvent.ID).
					Str("event_type", string(invoiceEvent.Type)).
					Msg("Canal de eventos de facturación lleno, descartando evento")
			}

		case <-s.stopChan:
			s.logger.Info(ctx).Msg("Deteniendo suscriptor Redis de facturación")
			return
		case <-ctx.Done():
			s.logger.Info(ctx).Msg("Context cancelado, deteniendo suscriptor Redis de facturación")
			return
		}
	}
}

// GetEventChannel retorna el canal de eventos para consumo externo
func (s *InvoiceEventSubscriber) GetEventChannel() <-chan *domain.InvoiceEvent {
	return s.eventChan
}

// Stop detiene el suscriptor
func (s *InvoiceEventSubscriber) Stop() error {
	close(s.stopChan)
	if s.pubsub != nil {
		return s.pubsub.Close()
	}
	return nil
}
