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

type IQueue interface {
	Publish(ctx context.Context, queueName string, message []byte) error

	PublishToExchange(ctx context.Context, exchangeName string, routingKey string, message []byte) error

	Consume(ctx context.Context, queueName string, handler func([]byte) error) error

	ConsumeConcurrent(ctx context.Context, queueName string, handler func([]byte) error, workers int) error

	Close() error

	DeclareQueue(queueName string, durable bool) error

	DeclareExchange(exchangeName string, exchangeType string, durable bool) error

	BindQueue(queueName string, exchangeName string, routingKey string) error

	Ping() error
}

type QueueRegistryCallback func(queueName string)

type consumerRegistration struct {
	queueName string
	handler   func([]byte) error
	ctx       context.Context
	workers   int
}

type rabbitMQ struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	logger        log.ILogger
	config        env.IConfig
	queueRegistry QueueRegistryCallback

	mu        sync.RWMutex
	consumers []consumerRegistration
	done      chan struct{}
}

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

func (r *rabbitMQ) SetQueueRegistry(callback QueueRegistryCallback) {
	r.queueRegistry = callback
}

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

func (r *rabbitMQ) watchConnection() {
	closeChan := make(chan *amqp.Error, 1)
	r.conn.NotifyClose(closeChan)

	go func() {
		select {
		case amqpErr, ok := <-closeChan:
			if !ok {
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
					Msg("RabbitMQ connection lost - starting automatic reconnection")
			} else {
				r.logger.Warn().
					Msg("RabbitMQ connection closed unexpectedly - starting automatic reconnection")
			}
			r.reconnect()

		case <-r.done:
			return
		}
	}()
}

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
			Msg("Attempting RabbitMQ reconnection...")

		time.Sleep(backoff)

		r.mu.Lock()
		err := r.connect()
		if err != nil {
			r.mu.Unlock()
			r.logger.Error().
				Err(err).
				Int("attempt", attempt).
				Dur("next_backoff", backoff*2).
				Msg("RabbitMQ reconnection failed - will retry")

			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		consumerCount := len(r.consumers)
		r.logger.Info().
			Int("attempt", attempt).
			Int("consumers_to_restore", consumerCount).
			Msg("RabbitMQ reconnected successfully - re-registering consumers")

		r.reregisterConsumers()
		r.mu.Unlock()

		r.watchConnection()
		return
	}
}

func (r *rabbitMQ) reregisterConsumers() {
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
		if err := r.startConsumer(c.ctx, c.queueName, c.handler, c.workers); err != nil {
			r.logger.Error().
				Err(err).
				Str("queue", c.queueName).
				Msg("Failed to re-register consumer after reconnection")
		} else {
			r.logger.Info().
				Str("queue", c.queueName).
				Msg("Consumer re-registered successfully")
		}
	}
}

const consumerPrefetchCount = 50

func (r *rabbitMQ) startConsumer(ctx context.Context, queueName string, handler func([]byte) error, workers int) error {
	if workers < 1 {
		workers = 1
	}

	consumerChannel, err := r.conn.Channel()
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Failed to create channel for consumer")
		return fmt.Errorf("failed to create consumer channel: %w", err)
	}

	if err := consumerChannel.Qos(consumerPrefetchCount, 0, false); err != nil {
		consumerChannel.Close()
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Failed to set consumer QoS/prefetch")
		return fmt.Errorf("failed to set consumer QoS: %w", err)
	}

	msgs, err := consumerChannel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		consumerChannel.Close()
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Error al registrar consumer")
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	for i := 0; i < workers; i++ {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					r.logger.Info().
						Str("queue", queueName).
						Int("worker", workerID).
						Msg("Stopping consumer due to context cancellation")
					return
				case msg, ok := <-msgs:
					if !ok {
						r.logger.Warn().
							Str("queue", queueName).
							Int("worker", workerID).
							Msg("Consumer channel closed - will be restored on reconnection")
						return
					}

					r.logger.Debug().
						Str("queue", queueName).
						Int("worker", workerID).
						Int("message_size", len(msg.Body)).
						Msg("Message received from queue - processing")

					if err := handler(msg.Body); err != nil {
						r.logger.Error().
							Err(err).
							Str("queue", queueName).
							Msg("Error processing message")
						r.logger.Debug().
							Err(err).
							Str("queue", queueName).
							Msg("Message processing FAILED - will be requeued")
						msg.Nack(false, true)
					} else {
						r.logger.Debug().
							Str("queue", queueName).
							Msg("Message processed successfully - ACK sent")
						msg.Ack(false)
					}
				}
			}
		}(i)
	}

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
		"",
		queueName,
		false,
		false,
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
	return r.ConsumeConcurrent(ctx, queueName, handler, 1)
}

func (r *rabbitMQ) ConsumeConcurrent(ctx context.Context, queueName string, handler func([]byte) error, workers int) error {
	if workers < 1 {
		workers = 1
	}

	r.mu.RLock()
	if r.conn == nil {
		r.mu.RUnlock()
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	err := r.startConsumer(ctx, queueName, handler, workers)
	r.mu.RUnlock()

	if err != nil {
		return err
	}

	r.mu.Lock()
	r.consumers = append(r.consumers, consumerRegistration{
		queueName: queueName,
		handler:   handler,
		ctx:       ctx,
		workers:   workers,
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
		queueName,
		durable,
		false,
		false,
		false,
		nil,
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
		exchangeName,
		routingKey,
		false,
		false,
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
		exchangeName,
		exchangeType,
		durable,
		false,
		false,
		false,
		nil,
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
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
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
