package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IQueue define la interfaz para manejar colas (RabbitMQ, etc.)
type IQueue interface {
	// Publish publica un mensaje en una cola específica (legacy - usar PublishToExchange)
	Publish(ctx context.Context, queueName string, message []byte) error

	// PublishToExchange publica un mensaje a un exchange
	PublishToExchange(ctx context.Context, exchangeName string, routingKey string, message []byte) error

	// Consume consume mensajes de una cola específica
	// El handler se ejecuta para cada mensaje recibido
	Consume(ctx context.Context, queueName string, handler func([]byte) error) error

	// Close cierra la conexión con el sistema de colas
	Close() error

	// DeclareQueue declara/crea una cola si no existe
	DeclareQueue(queueName string, durable bool) error

	// DeclareExchange declara/crea un exchange si no existe
	DeclareExchange(exchangeName string, exchangeType string, durable bool) error

	// BindQueue vincula una cola a un exchange con un routing key
	BindQueue(queueName string, exchangeName string, routingKey string) error

	// Ping verifica que la conexión esté activa
	Ping() error
}

// QueueRegistryCallback es un callback para registrar colas declaradas
type QueueRegistryCallback func(queueName string)

// consumerRegistration almacena la info necesaria para re-registrar un consumer
// después de una reconexión a RabbitMQ.
type consumerRegistration struct {
	queueName string
	handler   func([]byte) error
	ctx       context.Context
}

type rabbitMQ struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	logger        log.ILogger
	config        env.IConfig
	queueRegistry QueueRegistryCallback

	// Reconexión automática
	mu        sync.RWMutex           // Protege conn/channel durante reconexión
	consumers []consumerRegistration // Registro de consumers para re-registrar
	done      chan struct{}           // Señal de cierre intencional (Close())
}

// New crea una nueva instancia de RabbitMQ y conecta automáticamente
func New(logger log.ILogger, config env.IConfig) (IQueue, error) {
	r := &rabbitMQ{
		logger:    logger,
		config:    config,
		consumers: make([]consumerRegistration, 0),
		done:      make(chan struct{}),
	}

	if err := r.connect(); err != nil {
		return nil, err
	}

	r.watchConnection()

	return r, nil
}

// SetQueueRegistry establece un callback para registrar colas declaradas
func (r *rabbitMQ) SetQueueRegistry(callback QueueRegistryCallback) {
	r.queueRegistry = callback
}

// connect establece la conexión AMQP y crea el channel de publish.
// NO adquiere mutex — el caller es responsable de tener el lock apropiado.
func (r *rabbitMQ) connect() error {
	host := r.config.Get("RABBITMQ_HOST")
	port := r.config.Get("RABBITMQ_PORT")
	user := r.config.Get("RABBITMQ_USER")
	pass := r.config.Get("RABBITMQ_PASS")
	vhost := r.config.Get("RABBITMQ_VHOST")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5672"
	}
	if user == "" {
		user = "guest"
	}
	if pass == "" {
		pass = "guest"
	}
	if vhost == "" {
		vhost = "/"
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s%s", user, pass, host, port, vhost)

	var err error
	r.conn, err = amqp.Dial(url)
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to connect to RabbitMQ")
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to open RabbitMQ channel")
		return fmt.Errorf("failed to open channel: %w", err)
	}

	return nil
}

// watchConnection escucha NotifyClose de la conexión AMQP y dispara reconexión automática.
func (r *rabbitMQ) watchConnection() {
	closeChan := make(chan *amqp.Error, 1)
	r.conn.NotifyClose(closeChan)

	go func() {
		select {
		case amqpErr, ok := <-closeChan:
			if !ok {
				// Canal cerrado sin error — verificar si fue intencional
				select {
				case <-r.done:
					return
				default:
				}
			}
			if amqpErr != nil {
				r.logger.Error().
					Int("code", amqpErr.Code).
					Str("reason", amqpErr.Reason).
					Msg("🔴 RabbitMQ connection lost - starting automatic reconnection")
			} else {
				r.logger.Warn().
					Msg("🔴 RabbitMQ connection closed unexpectedly - starting automatic reconnection")
			}
			r.reconnect()

		case <-r.done:
			return
		}
	}()
}

// reconnect intenta reconectar con backoff exponencial y re-registra todos los consumers.
func (r *rabbitMQ) reconnect() {
	backoff := time.Second
	maxBackoff := 30 * time.Second

	for attempt := 1; ; attempt++ {
		select {
		case <-r.done:
			r.logger.Info().Msg("Reconnection cancelled - intentional shutdown")
			return
		default:
		}

		r.logger.Info().
			Int("attempt", attempt).
			Dur("backoff", backoff).
			Msg("⏳ Attempting RabbitMQ reconnection...")

		time.Sleep(backoff)

		r.mu.Lock()
		err := r.connect()
		if err != nil {
			r.mu.Unlock()
			r.logger.Error().
				Err(err).
				Int("attempt", attempt).
				Dur("next_backoff", backoff*2).
				Msg("❌ RabbitMQ reconnection failed - will retry")

			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		// Conexión restaurada — re-registrar consumers
		consumerCount := len(r.consumers)
		r.logger.Info().
			Int("attempt", attempt).
			Int("consumers_to_restore", consumerCount).
			Msg("✅ RabbitMQ reconnected successfully - re-registering consumers")

		r.reregisterConsumers()
		r.mu.Unlock()

		// Iniciar watcher en la nueva conexión
		r.watchConnection()
		return
	}
}

// reregisterConsumers re-crea los channels y goroutines de todos los consumers registrados.
// DEBE ser llamado con r.mu.Lock() adquirido.
func (r *rabbitMQ) reregisterConsumers() {
	// Filtrar consumers cuyo contexto fue cancelado
	active := make([]consumerRegistration, 0, len(r.consumers))
	for _, c := range r.consumers {
		select {
		case <-c.ctx.Done():
			r.logger.Info().
				Str("queue", c.queueName).
				Msg("Skipping consumer re-registration - context cancelled")
			continue
		default:
			active = append(active, c)
		}
	}
	r.consumers = active

	for _, c := range r.consumers {
		if err := r.startConsumer(c.ctx, c.queueName, c.handler); err != nil {
			r.logger.Error().
				Err(err).
				Str("queue", c.queueName).
				Msg("❌ Failed to re-register consumer after reconnection")
		} else {
			r.logger.Info().
				Str("queue", c.queueName).
				Msg("✅ Consumer re-registered successfully")
		}
	}
}

// startConsumer crea un channel dedicado e inicia la goroutine de consumo.
// NO adquiere mutex — el caller debe tener al menos RLock.
func (r *rabbitMQ) startConsumer(ctx context.Context, queueName string, handler func([]byte) error) error {
	consumerChannel, err := r.conn.Channel()
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Failed to create channel for consumer")
		return fmt.Errorf("failed to create consumer channel: %w", err)
	}

	msgs, err := consumerChannel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		consumerChannel.Close()
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Error al registrar consumer")
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				r.logger.Info().
					Str("queue", queueName).
					Msg("Stopping consumer due to context cancellation")
				return
			case msg, ok := <-msgs:
				if !ok {
					r.logger.Warn().
						Str("queue", queueName).
						Msg("Consumer channel closed - will be restored on reconnection")
					return
				}

				r.logger.Debug().
					Str("queue", queueName).
					Int("message_size", len(msg.Body)).
					Msg("📨 Message received from queue - processing")

				if err := handler(msg.Body); err != nil {
					r.logger.Error().
						Err(err).
						Str("queue", queueName).
						Msg("Error processing message")
					r.logger.Debug().
						Err(err).
						Str("queue", queueName).
						Msg("❌ Message processing FAILED - will be requeued")
					msg.Nack(false, true)
				} else {
					r.logger.Debug().
						Str("queue", queueName).
						Msg("✅ Message processed successfully - ACK sent")
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

func (r *rabbitMQ) Publish(ctx context.Context, queueName string, message []byte) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.channel == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}

	err := r.channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	if err != nil {
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Int("message_size", len(message)).
			Msg("Failed to publish message to queue")
		return fmt.Errorf("failed to publish message: %w", err)
	}

	r.logger.Info().
		Str("queue", queueName).
		Int("message_size", len(message)).
		Msg("Message published to queue")

	return nil
}

func (r *rabbitMQ) Consume(ctx context.Context, queueName string, handler func([]byte) error) error {
	r.mu.RLock()
	if r.conn == nil {
		r.mu.RUnlock()
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	err := r.startConsumer(ctx, queueName, handler)
	r.mu.RUnlock()

	if err != nil {
		return err
	}

	// Registrar consumer para re-creación automática tras reconexión
	r.mu.Lock()
	r.consumers = append(r.consumers, consumerRegistration{
		queueName: queueName,
		handler:   handler,
		ctx:       ctx,
	})
	r.mu.Unlock()

	return nil
}

func (r *rabbitMQ) DeclareQueue(queueName string, durable bool) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.conn == nil {
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	ch, err := r.conn.Channel()
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Error al crear channel para declarar cola")
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // name
		durable,   // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Bool("durable", durable).
			Msg("Error al declarar cola")
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if r.queueRegistry != nil {
		r.queueRegistry(queueName)
	}

	return nil
}

func (r *rabbitMQ) Ping() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.conn == nil || r.conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed")
	}
	if r.channel == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}
	return nil
}

func (r *rabbitMQ) PublishToExchange(ctx context.Context, exchangeName string, routingKey string, message []byte) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.channel == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}

	err := r.channel.PublishWithContext(
		ctx,
		exchangeName, // exchange
		routingKey,   // routing key (vacío para fanout)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	if err != nil {
		r.logger.Error().
			Err(err).
			Str("exchange", exchangeName).
			Str("routing_key", routingKey).
			Int("message_size", len(message)).
			Msg("Failed to publish message to exchange")
		return fmt.Errorf("failed to publish message: %w", err)
	}

	r.logger.Info().
		Str("exchange", exchangeName).
		Str("routing_key", routingKey).
		Int("message_size", len(message)).
		Msg("Message published to exchange")

	return nil
}

func (r *rabbitMQ) DeclareExchange(exchangeName string, exchangeType string, durable bool) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.conn == nil {
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	ch, err := r.conn.Channel()
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("exchange", exchangeName).
			Msg("Error al crear channel para declarar exchange")
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type (fanout, direct, topic, headers)
		durable,      // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		r.logger.Error().
			Err(err).
			Str("exchange", exchangeName).
			Str("type", exchangeType).
			Bool("durable", durable).
			Msg("Error al declarar exchange")
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	return nil
}

func (r *rabbitMQ) BindQueue(queueName string, exchangeName string, routingKey string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.conn == nil {
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	ch, err := r.conn.Channel()
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Str("exchange", exchangeName).
			Msg("Failed to create channel for queue binding")
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	err = ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key (vacío para fanout)
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Str("exchange", exchangeName).
			Str("routing_key", routingKey).
			Msg("Error al bindear cola a exchange")
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

func (r *rabbitMQ) Close() error {
	r.logger.Info().Msg("Closing RabbitMQ connection")

	// Señalar cierre intencional para que watchConnection no intente reconectar
	close(r.done)

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			r.logger.Error().
				Err(err).
				Msg("Error closing RabbitMQ channel")
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			r.logger.Error().
				Err(err).
				Msg("Error closing RabbitMQ connection")
			return err
		}
	}

	r.logger.Info().Msg("RabbitMQ connection closed successfully")
	return nil
}
