package rabbitmq

import (
	"context"
	"fmt"

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

type rabbitMQ struct {
	conn           *amqp.Connection
	channel        *amqp.Channel
	logger         log.ILogger
	config         env.IConfig
	queueRegistry  QueueRegistryCallback
}

// New crea una nueva instancia de RabbitMQ y conecta automáticamente
func New(logger log.ILogger, config env.IConfig) (IQueue, error) {
	r := &rabbitMQ{
		logger: logger,
		config: config,
	}

	if err := r.connect(); err != nil {
		return nil, err
	}

	return r, nil
}

// SetQueueRegistry establece un callback para registrar colas declaradas
func (r *rabbitMQ) SetQueueRegistry(callback QueueRegistryCallback) {
	r.queueRegistry = callback
}

func (r *rabbitMQ) connect() error {
	// Construir URL de conexión desde variables de entorno
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

	// Conexión silenciosa - info mostrada en LogStartupInfo()
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

	// Conexión exitosa - info mostrada en LogStartupInfo()
	return nil
}

func (r *rabbitMQ) Publish(ctx context.Context, queueName string, message []byte) error {
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
	if r.conn == nil {
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	// Crear un channel SEPARADO para este consumer
	// Esto evita el error "unexpected command received" cuando múltiples consumers
	// intentan usar el mismo channel
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
						Msg("Consumer channel closed")
					return
				}

				// Log del mensaje recibido
				r.logger.Debug().
					Str("queue", queueName).
					Int("message_size", len(msg.Body)).
					Msg("📨 Message received from queue - processing")

				// Procesamos el mensaje
				if err := handler(msg.Body); err != nil {
					r.logger.Error().
						Err(err).
						Str("queue", queueName).
						Msg("Error processing message")
					r.logger.Debug().
						Err(err).
						Str("queue", queueName).
						Msg("❌ Message processing FAILED - will be requeued")
					// Nack the message so it can be requeued
					msg.Nack(false, true)
				} else {
					r.logger.Debug().
						Str("queue", queueName).
						Msg("✅ Message processed successfully - ACK sent")
					// Ack the message
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

func (r *rabbitMQ) DeclareQueue(queueName string, durable bool) error {
	if r.conn == nil {
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	// Crear un channel temporal para declarar la cola
	// Esto evita conflictos con otros channels que puedan estar en uso
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

	// Registrar cola declarada exitosamente
	if r.queueRegistry != nil {
		r.queueRegistry(queueName)
	}

	return nil
}

func (r *rabbitMQ) Ping() error {
	if r.conn == nil || r.conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed")
	}
	if r.channel == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}
	return nil
}

func (r *rabbitMQ) PublishToExchange(ctx context.Context, exchangeName string, routingKey string, message []byte) error {
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
	if r.conn == nil {
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	// Crear un channel temporal para declarar el exchange
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

	// Silencioso - solo loguear errores
	return nil
}

func (r *rabbitMQ) BindQueue(queueName string, exchangeName string, routingKey string) error {
	if r.conn == nil {
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	// Crear un channel temporal para hacer el binding
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

	// Silencioso - solo loguear errores
	return nil
}

func (r *rabbitMQ) Close() error {
	r.logger.Info().Msg("Closing RabbitMQ connection")

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
