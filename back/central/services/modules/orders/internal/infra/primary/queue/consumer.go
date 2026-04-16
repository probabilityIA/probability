package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const OrdersCanonicalQueueName = rabbitmq.QueueOrdersCanonical

type OrderConsumer struct {
	queue          rabbitmq.IQueue
	logger         log.ILogger
	createUC       ports.IOrderCreateUseCase
	repo           ports.IRepository
	eventPublisher ports.IIntegrationEventPublisher
}

func New(
	queue rabbitmq.IQueue,
	logger log.ILogger,
	createUC ports.IOrderCreateUseCase,
	repo ports.IRepository,
	eventPublisher ports.IIntegrationEventPublisher,
) ports.IOrderConsumer {
	return &OrderConsumer{
		queue:          queue,
		logger:         logger,
		createUC:       createUC,
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}

func (c *OrderConsumer) Start(ctx context.Context) error {
	if err := c.queue.DeclareQueue(OrdersCanonicalQueueName, true); err != nil {
		c.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Msg("Failed to declare queue")
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := c.queue.Consume(ctx, OrdersCanonicalQueueName, c.handleMessage); err != nil {
		c.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Msg("Failed to start consumer")
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	return nil
}

func (c *OrderConsumer) handleMessage(messageBody []byte) error {
	ctx := context.Background()

	c.logger.Debug().
		Str("queue", OrdersCanonicalQueueName).
		Int("message_size", len(messageBody)).
		Msg("Processing order message from queue")

	var orderDTO dtos.ProbabilityOrderDTO
	if err := json.Unmarshal(messageBody, &orderDTO); err != nil {
		c.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Str("message_body", string(messageBody)).
			Msg("Failed to unmarshal order message")

		c.saveOrderError(ctx, nil, err, "unmarshal_error", messageBody)
		c.publishRejected(ctx, &orderDTO, "Error de formato: "+err.Error())
		return nil
	}

	if orderDTO.ExternalID == "" {
		err := fmt.Errorf("order message missing external_id")
		c.logger.Error().
			Str("queue", OrdersCanonicalQueueName).
			Msg("Order message missing external_id")
		c.saveOrderError(ctx, &orderDTO, err, "validation_error", messageBody)
		c.publishRejected(ctx, &orderDTO, "Falta external_id")
		return nil
	}

	if orderDTO.IntegrationID == 0 {
		err := fmt.Errorf("order message missing integration_id")
		c.logger.Error().
			Str("queue", OrdersCanonicalQueueName).
			Str("external_id", orderDTO.ExternalID).
			Msg("Order message missing integration_id")
		c.saveOrderError(ctx, &orderDTO, err, "validation_error", messageBody)
		c.publishRejected(ctx, &orderDTO, "Falta integration_id")
		return nil
	}

	orderResponse, err := c.createUC.MapAndSaveOrder(ctx, &orderDTO)
	if err != nil {
		errStr := err.Error()
		if errors.Is(err, domainerrors.ErrOrderAlreadyExists) {
			c.logger.Info().
				Str("queue", OrdersCanonicalQueueName).
				Str("external_id", orderDTO.ExternalID).
				Msg("Order already exists, skipping")
			c.publishRejected(ctx, &orderDTO, "Orden duplicada")
			return nil
		}

		if contains(errStr, "business_id is required") || contains(errStr, "integration_id is required") {
			c.logger.Warn().
				Str("queue", OrdersCanonicalQueueName).
				Str("external_id", orderDTO.ExternalID).
				Msg("Discarding invalid message: missing required fields (drain queue)")
			c.publishRejected(ctx, &orderDTO, "Campos requeridos faltantes")
			return nil
		}

		if contains(errStr, "violates foreign key constraint") {
			c.logger.Warn().
				Err(err).
				Str("queue", OrdersCanonicalQueueName).
				Str("external_id", orderDTO.ExternalID).
				Msg("Order failed with data integrity error (FK violation), discarding message")
			c.publishRejected(ctx, &orderDTO, "Error de integridad de datos")
			return nil
		}

		if contains(errStr, "duplicate key value violates unique constraint") &&
			(contains(errStr, "idx_integration_external_id") || contains(errStr, "SQLSTATE 23505")) {
			c.logger.Info().
				Str("queue", OrdersCanonicalQueueName).
				Str("external_id", orderDTO.ExternalID).
				Uint("integration_id", orderDTO.IntegrationID).
				Msg("Order already exists (race condition detected), skipping duplicate message")
			c.publishRejected(ctx, &orderDTO, "Orden duplicada")
			return nil
		}

		if contains(errStr, "value too long for type") || contains(errStr, "SQLSTATE 22001") {
			c.logger.Warn().
				Err(err).
				Str("queue", OrdersCanonicalQueueName).
				Str("external_id", orderDTO.ExternalID).
				Uint("integration_id", orderDTO.IntegrationID).
				Msg("Order failed with data length error (varchar overflow), discarding message")
			c.saveOrderError(ctx, &orderDTO, err, "data_length_error", messageBody)
			c.publishRejected(ctx, &orderDTO, "Error de longitud de datos")
			return nil
		}

		c.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Str("external_id", orderDTO.ExternalID).
			Uint("integration_id", orderDTO.IntegrationID).
			Str("platform", orderDTO.Platform).
			Msg("Failed to map and save order")

		c.saveOrderError(ctx, &orderDTO, err, "processing_error", messageBody)
		c.publishRejected(ctx, &orderDTO, "Error procesando: "+err.Error())
		c.publishAIOrderResult(ctx, &orderDTO, false, "", err.Error())
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

	c.publishAIOrderResult(ctx, &orderDTO, true, orderResponse.ID, "")

	return nil
}

func (c *OrderConsumer) saveOrderError(ctx context.Context, orderDTO *dtos.ProbabilityOrderDTO, err error, errorType string, messageBody []byte) {
	if c.repo == nil {
		c.logger.Warn().Msg("Repository not available, cannot save order error")
		return
	}

	if errorType == "" {
		errMsg := err.Error()
		if strings.Contains(errMsg, "validation") || strings.Contains(errMsg, "required") {
			errorType = "validation_error"
		} else if strings.Contains(errMsg, "database") || strings.Contains(errMsg, "constraint") {
			errorType = "database_error"
		} else {
			errorType = "processing_error"
		}
	}

	var externalID string
	var integrationID uint
	var businessID *uint
	var integrationType string
	var platform string

	if orderDTO != nil {
		externalID = orderDTO.ExternalID
		integrationID = orderDTO.IntegrationID
		businessID = orderDTO.BusinessID
		integrationType = orderDTO.IntegrationType
		platform = orderDTO.Platform
	} else {
		var rawMap map[string]interface{}
		if json.Unmarshal(messageBody, &rawMap) == nil {
			if extID, ok := rawMap["external_id"].(string); ok {
				externalID = extID
			}
			if intID, ok := rawMap["integration_id"].(float64); ok {
				integrationID = uint(intID)
			}
			if busID, ok := rawMap["business_id"].(float64); ok {
				bid := uint(busID)
				businessID = &bid
			}
			if intType, ok := rawMap["integration_type"].(string); ok {
				integrationType = intType
			}
			if plat, ok := rawMap["platform"].(string); ok {
				platform = plat
			}
		}
	}

	orderError := &entities.OrderError{
		ExternalID:      externalID,
		IntegrationID:   integrationID,
		BusinessID:      businessID,
		IntegrationType: integrationType,
		Platform:        platform,
		ErrorType:       errorType,
		ErrorMessage:    err.Error(),
		RawData:         messageBody, // JSON original
		Status:          "new",
	}

	if saveErr := c.repo.CreateOrderError(ctx, orderError); saveErr != nil {
		c.logger.Error().
			Err(saveErr).
			Msg("Failed to save order error to database")
	}
}

func (c *OrderConsumer) publishRejected(ctx context.Context, dto *dtos.ProbabilityOrderDTO, reason string) {
	if c.eventPublisher == nil || dto == nil {
		return
	}

	c.eventPublisher.PublishSyncOrderRejected(ctx, dto.IntegrationID, dto.BusinessID, map[string]interface{}{
		"order_number": dto.OrderNumber,
		"external_id":  dto.ExternalID,
		"reason":       reason,
		"rejected_at":  time.Now().Format(time.RFC3339),
	})
}

func (c *OrderConsumer) publishAIOrderResult(ctx context.Context, dto *dtos.ProbabilityOrderDTO, success bool, orderID string, errMsg string) {
	if dto == nil || dto.Platform != "whatsapp_ai" {
		return
	}

	var businessID uint
	if dto.BusinessID != nil {
		businessID = *dto.BusinessID
	}

	result := map[string]interface{}{
		"external_id":   dto.ExternalID,
		"phone_number":  dto.CustomerPhone,
		"business_id":   businessID,
		"success":       success,
		"order_id":      orderID,
		"order_number":  dto.OrderNumber,
		"error_message": errMsg,
	}

	payload, err := json.Marshal(result)
	if err != nil {
		c.logger.Error().Err(err).Msg("Error marshaling AI order result")
		return
	}

	if err := c.queue.Publish(ctx, rabbitmq.QueueAIOrderResult, payload); err != nil {
		c.logger.Error().
			Err(err).
			Str("external_id", dto.ExternalID).
			Msg("Error publishing AI order result")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
