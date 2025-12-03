package queue

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IQueue define la interfaz para manejar colas (RabbitMQ, etc.)
type IQueue interface {
	// Publish publica un mensaje en una cola específica
	Publish(ctx context.Context, queueName string, message []byte) error

	// Consume consume mensajes de una cola específica
	// El handler se ejecuta para cada mensaje recibido
	Consume(ctx context.Context, queueName string, handler func([]byte) error) error

	// Close cierra la conexión con el sistema de colas
	Close() error

	// DeclareQueue declara/crea una cola si no existe
	DeclareQueue(queueName string, durable bool) error

	// Ping verifica que la conexión esté activa
	Ping() error
}

type rabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  log.ILogger
	config  env.IConfig
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

	r.logger.Debug().
		Str("host", host).
		Str("port", port).
		Str("vhost", vhost).
		Msg("Connecting to RabbitMQ")

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

	r.logger.Debug().Msg("Successfully connected to RabbitMQ")
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
	if r.channel == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}

	msgs, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Failed to register consumer")
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	r.logger.Info().
		Str("queue", queueName).
		Msg("Started consuming messages from queue")

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

				r.logger.Debug().
					Str("queue", queueName).
					Int("message_size", len(msg.Body)).
					Msg("Received message from queue")

				if err := handler(msg.Body); err != nil {
					r.logger.Error().
						Err(err).
						Str("queue", queueName).
						Msg("Error processing message")
					// Nack the message so it can be requeued
					msg.Nack(false, true)
				} else {
					// Ack the message
					msg.Ack(false)
					r.logger.Debug().
						Str("queue", queueName).
						Msg("Message processed successfully")
				}
			}
		}
	}()

	return nil
}

func (r *rabbitMQ) DeclareQueue(queueName string, durable bool) error {
	if r.channel == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}

	_, err := r.channel.QueueDeclare(
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
			Msg("Failed to declare queue")
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	r.logger.Info().
		Str("queue", queueName).
		Bool("durable", durable).
		Msg("Queue declared successfully")

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
