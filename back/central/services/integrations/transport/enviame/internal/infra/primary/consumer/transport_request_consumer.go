package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// ICredentialResolver resolves encrypted credentials for an integration.
// Satisfied by core.IIntegrationService (replicated interface — module isolation).
type ICredentialResolver interface {
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
}

// TransportRequestMessage is the message received from the transport router
type TransportRequestMessage struct {
	ShipmentID    *uint                  `json:"shipment_id,omitempty"`
	Provider      string                 `json:"provider"`
	Operation     string                 `json:"operation"`
	CorrelationID string                 `json:"correlation_id"`
	BusinessID    uint                   `json:"business_id"`
	IntegrationID uint                   `json:"integration_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Payload       map[string]interface{} `json:"payload"`
}

const (
	QueueEnviameRequests = rabbitmq.QueueTransportEnviameRequests
)

// TransportRequestConsumer consumes transport requests for Enviame
type TransportRequestConsumer struct {
	rabbit             rabbitmq.IQueue
	useCase            app.IUseCase
	responsePublisher  *queue.ResponsePublisher
	credentialResolver ICredentialResolver
	log                log.ILogger
}

// NewTransportRequestConsumer creates a new consumer
func NewTransportRequestConsumer(
	rabbit rabbitmq.IQueue,
	useCase app.IUseCase,
	responsePublisher *queue.ResponsePublisher,
	credentialResolver ICredentialResolver,
	logger log.ILogger,
) *TransportRequestConsumer {
	return &TransportRequestConsumer{
		rabbit:             rabbit,
		useCase:            useCase,
		responsePublisher:  responsePublisher,
		credentialResolver: credentialResolver,
		log:                logger.WithModule("transport.enviame.consumer"),
	}
}

// Start starts consuming from the Enviame requests queue
func (c *TransportRequestConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		c.log.Warn(ctx).Msg("RabbitMQ client is nil, consumer cannot start")
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", QueueEnviameRequests).
		Msg("Starting Enviame transport request consumer")

	if err := c.rabbit.DeclareQueue(QueueEnviameRequests, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare queue")
		return err
	}

	if err := c.rabbit.Consume(ctx, QueueEnviameRequests, c.handleRequest); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to start consuming")
		return err
	}

	return nil
}

// handleRequest processes a transport request
func (c *TransportRequestConsumer) handleRequest(message []byte) error {
	ctx := context.Background()

	var request TransportRequestMessage
	if err := json.Unmarshal(message, &request); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("body", string(message)).
			Msg("Failed to unmarshal request")
		return err
	}

	c.log.Info(ctx).
		Str("operation", request.Operation).
		Str("correlation_id", request.CorrelationID).
		Uint("integration_id", request.IntegrationID).
		Msg("📨 Received transport request")

	apiKey, err := c.resolveAPIKey(ctx, &request)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("correlation_id", request.CorrelationID).
			Msg("Failed to resolve API key")
		response := c.errorResponse(&request, "Error al resolver credenciales: "+err.Error())
		if pubErr := c.responsePublisher.PublishResponse(ctx, response); pubErr != nil {
			c.log.Error(ctx).Err(pubErr).Msg("Failed to publish error response")
		}
		return err
	}

	var response *queue.TransportResponseMessage

	switch request.Operation {
	case "quote":
		response = c.processQuote(ctx, &request, apiKey)
	case "generate":
		response = c.processGenerate(ctx, &request, apiKey)
	case "track":
		response = c.processTrack(ctx, &request, apiKey)
	case "cancel":
		response = c.processCancel(ctx, &request, apiKey)
	default:
		c.log.Warn(ctx).
			Str("operation", request.Operation).
			Msg("Unknown operation")
		response = c.errorResponse(&request, "Unknown operation: "+request.Operation)
	}

	if err := c.responsePublisher.PublishResponse(ctx, response); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("correlation_id", request.CorrelationID).
			Msg("Failed to publish response")
		return err
	}

	return nil
}

// resolveAPIKey decrypts the api_key credential for the integration.
func (c *TransportRequestConsumer) resolveAPIKey(ctx context.Context, request *TransportRequestMessage) (string, error) {
	if request.IntegrationID == 0 {
		if envKey := os.Getenv("ENVIAME_API_KEY"); envKey != "" {
			c.log.Warn(ctx).Msg("integration_id is 0, using ENVIAME_API_KEY env var fallback")
			return envKey, nil
		}
		return "", fmt.Errorf("integration_id is 0, cannot resolve credentials")
	}

	apiKey, err := c.credentialResolver.DecryptCredential(ctx, fmt.Sprintf("%d", request.IntegrationID), "api_key")
	if err != nil {
		if envKey := os.Getenv("ENVIAME_API_KEY"); envKey != "" {
			c.log.Warn(ctx).
				Uint("integration_id", request.IntegrationID).
				Msg("Failed to decrypt credentials, using ENVIAME_API_KEY env var fallback")
			return envKey, nil
		}
		return "", fmt.Errorf("failed to decrypt api_key for integration %d: %w", request.IntegrationID, err)
	}

	return apiKey, nil
}

// processQuote handles shipping rate quotes
// TODO: Implement when Enviame API is ready
func (c *TransportRequestConsumer) processQuote(ctx context.Context, request *TransportRequestMessage, apiKey string) *queue.TransportResponseMessage {
	c.log.Warn(ctx).Msg("⚠️ Enviame processQuote not yet implemented")
	return c.errorResponse(request, "enviame: quote not yet implemented")
}

// processGenerate handles guide generation
// TODO: Implement when Enviame API is ready
func (c *TransportRequestConsumer) processGenerate(ctx context.Context, request *TransportRequestMessage, apiKey string) *queue.TransportResponseMessage {
	c.log.Warn(ctx).Msg("⚠️ Enviame processGenerate not yet implemented")
	return c.errorResponse(request, "enviame: generate not yet implemented")
}

// processTrack handles shipment tracking
// TODO: Implement when Enviame API is ready
func (c *TransportRequestConsumer) processTrack(ctx context.Context, request *TransportRequestMessage, apiKey string) *queue.TransportResponseMessage {
	c.log.Warn(ctx).Msg("⚠️ Enviame processTrack not yet implemented")
	return c.errorResponse(request, "enviame: track not yet implemented")
}

// processCancel handles shipment cancellation
// TODO: Implement when Enviame API is ready
func (c *TransportRequestConsumer) processCancel(ctx context.Context, request *TransportRequestMessage, apiKey string) *queue.TransportResponseMessage {
	c.log.Warn(ctx).Msg("⚠️ Enviame processCancel not yet implemented")
	return c.errorResponse(request, "enviame: cancel not yet implemented")
}

// errorResponse creates an error response
func (c *TransportRequestConsumer) errorResponse(request *TransportRequestMessage, errMsg string) *queue.TransportResponseMessage {
	return &queue.TransportResponseMessage{
		ShipmentID:    request.ShipmentID,
		BusinessID:    request.BusinessID,
		Provider:      "enviame",
		Operation:     request.Operation,
		Status:        "error",
		CorrelationID: request.CorrelationID,
		Timestamp:     time.Now(),
		Error:         errMsg,
	}
}
