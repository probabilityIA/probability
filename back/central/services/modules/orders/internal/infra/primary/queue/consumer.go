package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	integrationevents "github.com/secamc93/probability/back/central/services/integrations/events"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"gorm.io/datatypes"
)

const (
	// OrdersCanonicalQueueName es el nombre de la cola donde se reciben órdenes canónicas
	OrdersCanonicalQueueName = "probability.orders.canonical"
)

// OrderConsumer consume órdenes canónicas de RabbitMQ y las procesa
// Implementa ports.IOrderConsumer
type OrderConsumer struct {
	queue          rabbitmq.IQueue
	logger         log.ILogger
	orderMappingUC ports.IOrderMappingUseCase
	repo           ports.IRepository
}

// New crea una nueva instancia del consumidor de órdenes
func New(
	queue rabbitmq.IQueue,
	logger log.ILogger,
	orderMappingUC ports.IOrderMappingUseCase,
	repo ports.IRepository,
) ports.IOrderConsumer {
	return &OrderConsumer{
		queue:          queue,
		logger:         logger,
		orderMappingUC: orderMappingUC,
		repo:           repo,
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

	return nil
}

// handleMessage procesa cada mensaje recibido de la cola
func (c *OrderConsumer) handleMessage(messageBody []byte) error {
	ctx := context.Background()

	c.logger.Debug().
		Str("queue", OrdersCanonicalQueueName).
		Int("message_size", len(messageBody)).
		Msg("Processing order message from queue")

	// Deserializar el mensaje a ProbabilityOrderDTO
	var orderDTO dtos.ProbabilityOrderDTO
	if err := json.Unmarshal(messageBody, &orderDTO); err != nil {
		c.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Str("message_body", string(messageBody)).
			Msg("Failed to unmarshal order message")

		// Guardar error con JSON original
		c.saveOrderError(ctx, nil, err, "unmarshal_error", messageBody)
		return fmt.Errorf("failed to unmarshal order message: %w", err)
	}

	// Validar que la orden tenga los campos mínimos requeridos
	if orderDTO.ExternalID == "" {
		err := fmt.Errorf("order message missing external_id")
		c.logger.Error().
			Str("queue", OrdersCanonicalQueueName).
			Msg("Order message missing external_id")
		c.saveOrderError(ctx, &orderDTO, err, "validation_error", messageBody)
		return err
	}

	if orderDTO.IntegrationID == 0 {
		err := fmt.Errorf("order message missing integration_id")
		c.logger.Error().
			Str("queue", OrdersCanonicalQueueName).
			Str("external_id", orderDTO.ExternalID).
			Msg("Order message missing integration_id")
		c.saveOrderError(ctx, &orderDTO, err, "validation_error", messageBody)
		return err
	}

	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "C",
			"location":     "consumer.go:116",
			"message":      "Consumer - Processing order from queue",
			"data": map[string]interface{}{
				"external_id":    orderDTO.ExternalID,
				"order_number":   orderDTO.OrderNumber,
				"integration_id": orderDTO.IntegrationID,
			},
			"timestamp": time.Now().UnixMilli(),
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	// Llamar al caso de uso para mapear y guardar la orden
	orderResponse, err := c.orderMappingUC.MapAndSaveOrder(ctx, &orderDTO)
	if err != nil {
		errStr := err.Error()
		// Check for specific errors to discard message
		if errors.Is(err, domainerrors.ErrOrderAlreadyExists) {
			c.logger.Info().
				Str("queue", OrdersCanonicalQueueName).
				Str("external_id", orderDTO.ExternalID).
				Msg("Order already exists, skipping")
			// No publicar evento de rechazo para órdenes duplicadas (es comportamiento esperado)
			return nil
		}

		// Discard messages with missing required business/integration checks from domain logic
		if contains(errStr, "business_id is required") || contains(errStr, "integration_id is required") {
			c.logger.Warn().
				Str("queue", OrdersCanonicalQueueName).
				Str("external_id", orderDTO.ExternalID).
				Msg("Discarding invalid message: missing required fields (drain queue)")
			// Publicar evento de orden rechazada
			integrationevents.PublishSyncOrderRejected(
				ctx,
				orderDTO.IntegrationID,
				orderDTO.BusinessID,
				integrationevents.SyncOrderRejectedEvent{
					OrderNumber: orderDTO.OrderNumber,
					ExternalID:  orderDTO.ExternalID,
					Platform:    orderDTO.Platform,
					Reason:      "Campos requeridos faltantes",
					Error:       errStr,
					RejectedAt:  time.Now(),
				},
			)
			return nil
		}

		// If error is a FK violation (data integrity), discard message to avoid loop
		if contains(errStr, "violates foreign key constraint") {
			c.logger.Warn().
				Err(err).
				Str("queue", OrdersCanonicalQueueName).
				Str("external_id", orderDTO.ExternalID).
				Msg("Order failed with data integrity error (FK violation), discarding message")
			// Publicar evento de orden rechazada
			integrationevents.PublishSyncOrderRejected(
				ctx,
				orderDTO.IntegrationID,
				orderDTO.BusinessID,
				integrationevents.SyncOrderRejectedEvent{
					OrderNumber: orderDTO.OrderNumber,
					ExternalID:  orderDTO.ExternalID,
					Platform:    orderDTO.Platform,
					Reason:      "Error de integridad de datos (FK violation)",
					Error:       errStr,
					RejectedAt:  time.Now(),
				},
			)
			return nil
		}

		// If error is a duplicate key violation for external_id (race condition), discard message
		if contains(errStr, "duplicate key value violates unique constraint") &&
			(contains(errStr, "idx_integration_external_id") || contains(errStr, "SQLSTATE 23505")) {
			c.logger.Info().
				Str("queue", OrdersCanonicalQueueName).
				Str("external_id", orderDTO.ExternalID).
				Uint("integration_id", orderDTO.IntegrationID).
				Msg("Order already exists (race condition detected), skipping duplicate message")
			return nil
		}

		c.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Str("external_id", orderDTO.ExternalID).
			Uint("integration_id", orderDTO.IntegrationID).
			Str("platform", orderDTO.Platform).
			Msg("Failed to map and save order")

		// Publicar evento de orden rechazada
		integrationevents.PublishSyncOrderRejected(
			ctx,
			orderDTO.IntegrationID,
			orderDTO.BusinessID,
			integrationevents.SyncOrderRejectedEvent{
				OrderNumber: orderDTO.OrderNumber,
				ExternalID:  orderDTO.ExternalID,
				Platform:    orderDTO.Platform,
				Reason:      "Error al procesar orden",
				Error:       errStr,
				RejectedAt:  time.Now(),
			},
		)

		// Guardar error con JSON original
		c.saveOrderError(ctx, &orderDTO, err, "processing_error", messageBody)
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

// saveOrderError guarda un error en la tabla order_errors con el JSON original
func (c *OrderConsumer) saveOrderError(ctx context.Context, orderDTO *dtos.ProbabilityOrderDTO, err error, errorType string, messageBody []byte) {
	if c.repo == nil {
		c.logger.Warn().Msg("Repository not available, cannot save order error")
		return
	}

	// Determinar el tipo de error basado en el mensaje
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

	// Extraer información del DTO si está disponible
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
		// Intentar extraer del JSON si el DTO no está disponible
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
		RawData:         datatypes.JSON(messageBody), // JSON original
		Status:          "new",
	}

	// Intentar guardar el error (no bloqueamos si falla)
	if saveErr := c.repo.CreateOrderError(ctx, orderError); saveErr != nil {
		c.logger.Error().
			Err(saveErr).
			Msg("Failed to save order error to database")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
