package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseordermapping"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	// OrdersCanonicalQueueName es el nombre de la cola donde se reciben órdenes canónicas
	OrdersCanonicalQueueName = "probability.orders.canonical"
)

// OrderConsumer consume órdenes canónicas de RabbitMQ y las procesa
// Implementa domain.IOrderConsumer
type OrderConsumer struct {
	queue          rabbitmq.IQueue
	logger         log.ILogger
	orderMappingUC usecaseordermapping.IOrderMappingUseCase
}

// New crea una nueva instancia del consumidor de órdenes
func New(
	queue rabbitmq.IQueue,
	logger log.ILogger,
	orderMappingUC usecaseordermapping.IOrderMappingUseCase,
) domain.IOrderConsumer {
	return &OrderConsumer{
		queue:          queue,
		logger:         logger,
		orderMappingUC: orderMappingUC,
	}
}

// Start inicia el consumidor de órdenes
func (c *OrderConsumer) Start(ctx context.Context) error {
	// Declarar la cola si no existe (durable para persistencia)
	if err := c.queue.DeclareQueue(OrdersCanonicalQueueName, true); err != nil {
		c.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Msg("Failed to declare queue")
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Iniciar el consumo de mensajes
	if err := c.queue.Consume(ctx, OrdersCanonicalQueueName, c.handleMessage); err != nil {
		c.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Msg("Failed to start consumer")
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	c.logger.Info().
		Str("queue", OrdersCanonicalQueueName).
		Msg("Order consumer started successfully")

	return nil
}

// handleMessage procesa cada mensaje recibido de la cola
func (c *OrderConsumer) handleMessage(messageBody []byte) error {
	ctx := context.Background()

	c.logger.Debug().
		Str("queue", OrdersCanonicalQueueName).
		Int("message_size", len(messageBody)).
		Msg("Processing order message from queue")

	// Deserializar el mensaje a CanonicalOrderDTO
	var orderDTO domain.CanonicalOrderDTO
	if err := json.Unmarshal(messageBody, &orderDTO); err != nil {
		c.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Str("message_body", string(messageBody)).
			Msg("Failed to unmarshal order message")
		return fmt.Errorf("failed to unmarshal order message: %w", err)
	}

	// Validar que la orden tenga los campos mínimos requeridos
	if orderDTO.ExternalID == "" {
		c.logger.Error().
			Str("queue", OrdersCanonicalQueueName).
			Msg("Order message missing external_id")
		return fmt.Errorf("order message missing external_id")
	}

	if orderDTO.IntegrationID == 0 {
		c.logger.Error().
			Str("queue", OrdersCanonicalQueueName).
			Str("external_id", orderDTO.ExternalID).
			Msg("Order message missing integration_id")
		return fmt.Errorf("order message missing integration_id")
	}

	// Llamar al caso de uso para mapear y guardar la orden
	orderResponse, err := c.orderMappingUC.MapAndSaveOrder(ctx, &orderDTO)
	if err != nil {
		if errors.Is(err, domain.ErrOrderAlreadyExists) {
			c.logger.Info().
				Str("queue", OrdersCanonicalQueueName).
				Str("external_id", orderDTO.ExternalID).
				Msg("Order already exists, skipping")
			return nil
		}

		c.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Str("external_id", orderDTO.ExternalID).
			Uint("integration_id", orderDTO.IntegrationID).
			Str("platform", orderDTO.Platform).
			Msg("Failed to map and save order")
		return fmt.Errorf("failed to map and save order: %w", err)
	}

	c.logger.Info().
		Str("queue", OrdersCanonicalQueueName).
		Str("order_id", orderResponse.ID).
		Str("external_id", orderDTO.ExternalID).
		Str("platform", orderDTO.Platform).
		Uint("integration_id", orderDTO.IntegrationID).
		Int("items_count", len(orderDTO.OrderItems)).
		Int("addresses_count", len(orderDTO.Addresses)).
		Int("payments_count", len(orderDTO.Payments)).
		Int("shipments_count", len(orderDTO.Shipments)).
		Msg("Order processed and saved successfully from queue")

	return nil
}
