package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// ICredentialResolver resolves encrypted credentials and config for an integration.
// Satisfied by core.IIntegrationService (replicated interface — module isolation).
type ICredentialResolver interface {
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	GetIntegrationConfig(ctx context.Context, integrationID string) (map[string]interface{}, error)
}

// TransportRequestMessage is the message received from the transport router
// (replicated locally — module isolation rule)
type TransportRequestMessage struct {
	ShipmentID    *uint                  `json:"shipment_id,omitempty"`
	Provider      string                 `json:"provider"`
	Operation     string                 `json:"operation"`
	CorrelationID string                 `json:"correlation_id"`
	BusinessID    uint                   `json:"business_id"`
	IntegrationID uint                   `json:"integration_id"`
	BaseURL       string                 `json:"base_url,omitempty"`
	IsTest        bool                   `json:"is_test,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Payload       map[string]interface{} `json:"payload"`
}

const (
	QueueEnvioclickRequests = "transport.envioclick.requests"
)

// TransportRequestConsumer consumes transport requests for EnvioClick
type TransportRequestConsumer struct {
	rabbit             rabbitmq.IQueue
	useCase            *app.UseCase
	responsePublisher  *queue.ResponsePublisher
	credentialResolver ICredentialResolver
	log                log.ILogger
}

// NewTransportRequestConsumer creates a new consumer
func NewTransportRequestConsumer(
	rabbit rabbitmq.IQueue,
	useCase *app.UseCase,
	responsePublisher *queue.ResponsePublisher,
	credentialResolver ICredentialResolver,
	logger log.ILogger,
) *TransportRequestConsumer {
	return &TransportRequestConsumer{
		rabbit:             rabbit,
		useCase:            useCase,
		responsePublisher:  responsePublisher,
		credentialResolver: credentialResolver,
		log:                logger.WithModule("transport.envioclick.consumer"),
	}
}

// Start starts consuming from the EnvioClick requests queue
func (c *TransportRequestConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		c.log.Warn(ctx).Msg("RabbitMQ client is nil, consumer cannot start")
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", QueueEnvioclickRequests).
		Msg("Starting EnvioClick transport request consumer")

	if err := c.rabbit.DeclareQueue(QueueEnvioclickRequests, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare queue")
		return err
	}

	if err := c.rabbit.Consume(ctx, QueueEnvioclickRequests, c.handleRequest); err != nil {
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

	// Resolve integration config (for use_platform_token and base_url_test)
	var integrationConfig map[string]interface{}
	if request.IntegrationID != 0 {
		cfg, cfgErr := c.credentialResolver.GetIntegrationConfig(ctx, fmt.Sprintf("%d", request.IntegrationID))
		if cfgErr != nil {
			c.log.Warn(ctx).Err(cfgErr).Uint("integration_id", request.IntegrationID).Msg("Failed to get integration config, using defaults")
		} else {
			integrationConfig = cfg
		}
	}

	// Resolve API key from integration credentials (or platform token)
	apiKey, err := c.resolveAPIKey(ctx, &request, integrationConfig)
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

	// Resolve effective base URL (base_url_test from config overrides msg.BaseURL)
	baseURL := c.resolveBaseURL(&request, integrationConfig)

	var response *queue.TransportResponseMessage

	switch request.Operation {
	case "quote":
		response = c.processQuote(ctx, &request, baseURL, apiKey)
	case "generate":
		response = c.processGenerate(ctx, &request, baseURL, apiKey)
	case "track":
		response = c.processTrack(ctx, &request, baseURL, apiKey)
	case "cancel":
		response = c.processCancel(ctx, &request, baseURL, apiKey)
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

// resolveAPIKey resolves the API key for EnvioClick.
// If config has use_platform_token=true, uses ENVIOCLICK_API_KEY env var (shared platform account).
// Otherwise decrypts the api_key stored in integration credentials.
func (c *TransportRequestConsumer) resolveAPIKey(ctx context.Context, request *TransportRequestMessage, config map[string]interface{}) (string, error) {
	// Check if the integration uses the platform shared token
	if config != nil {
		if usePlatform, ok := config["use_platform_token"].(bool); ok && usePlatform {
			apiKey := os.Getenv("ENVIOCLICK_API_KEY")
			if apiKey == "" {
				return "", fmt.Errorf("use_platform_token está activo pero ENVIOCLICK_API_KEY no está configurado en la plataforma")
			}
			c.log.Info(ctx).Uint("integration_id", request.IntegrationID).Msg("Using platform EnvioClick token")
			return apiKey, nil
		}
	}

	if request.IntegrationID == 0 {
		// No integration ID — try env var fallback
		if envKey := os.Getenv("ENVIOCLICK_API_KEY"); envKey != "" {
			c.log.Warn(ctx).Msg("integration_id is 0, using ENVIOCLICK_API_KEY env var fallback")
			return envKey, nil
		}
		return "", fmt.Errorf("integration_id is 0, cannot resolve credentials")
	}

	apiKey, err := c.credentialResolver.DecryptCredential(ctx, fmt.Sprintf("%d", request.IntegrationID), "api_key")
	if err != nil {
		// Fallback to env var when credentials are not stored in the integration
		if envKey := os.Getenv("ENVIOCLICK_API_KEY"); envKey != "" {
			c.log.Warn(ctx).
				Uint("integration_id", request.IntegrationID).
				Msg("Failed to decrypt credentials, using ENVIOCLICK_API_KEY env var fallback")
			return envKey, nil
		}
		return "", fmt.Errorf("failed to decrypt api_key for integration %d: %w", request.IntegrationID, err)
	}

	return apiKey, nil
}

// resolveBaseURL determines the effective base URL for EnvioClick.
// Priority: config.base_url_test > msg.BaseURL > DefaultBaseURL
func (c *TransportRequestConsumer) resolveBaseURL(request *TransportRequestMessage, config map[string]interface{}) string {
	if config != nil {
		if testURL, ok := config["base_url_test"].(string); ok && testURL != "" {
			return testURL
		}
	}
	if request.BaseURL != "" {
		return request.BaseURL
	}
	return "https://api.envioclickpro.com.co/api/v2"
}

// processQuote handles shipping rate quotes
func (c *TransportRequestConsumer) processQuote(ctx context.Context, request *TransportRequestMessage, baseURL, apiKey string) *queue.TransportResponseMessage {
	payloadBytes, err := json.Marshal(request.Payload)
	if err != nil {
		return c.errorResponse(request, "Failed to marshal payload: "+err.Error())
	}

	var req domain.QuoteRequest
	if err := json.Unmarshal(payloadBytes, &req); err != nil {
		return c.errorResponse(request, "Failed to unmarshal payload as QuoteRequest: "+err.Error())
	}

	resp, err := c.useCase.Quote(ctx, baseURL, apiKey, req)
	if err != nil {
		return c.errorResponse(request, err.Error())
	}

	return &queue.TransportResponseMessage{
		ShipmentID:    request.ShipmentID,
		BusinessID:    request.BusinessID,
		Provider:      "envioclick",
		Operation:     "quote",
		Status:        "success",
		CorrelationID: request.CorrelationID,
		IsTest:        request.IsTest,
		Timestamp:     time.Now(),
		Data:          toMap(resp),
	}
}

// processGenerate handles guide generation
func (c *TransportRequestConsumer) processGenerate(ctx context.Context, request *TransportRequestMessage, baseURL, apiKey string) *queue.TransportResponseMessage {
	payloadBytes, err := json.Marshal(request.Payload)
	if err != nil {
		return c.errorResponse(request, "Failed to marshal payload: "+err.Error())
	}

	var req domain.QuoteRequest
	if err := json.Unmarshal(payloadBytes, &req); err != nil {
		return c.errorResponse(request, "Failed to unmarshal payload as QuoteRequest: "+err.Error())
	}

	resp, err := c.useCase.Generate(ctx, baseURL, apiKey, req)
	if err != nil {
		return c.errorResponse(request, err.Error())
	}

	return &queue.TransportResponseMessage{
		ShipmentID:    request.ShipmentID,
		BusinessID:    request.BusinessID,
		Provider:      "envioclick",
		Operation:     "generate",
		Status:        "success",
		CorrelationID: request.CorrelationID,
		IsTest:        request.IsTest,
		Timestamp:     time.Now(),
		Data:          toMap(resp),
	}
}

// processTrack handles shipment tracking
func (c *TransportRequestConsumer) processTrack(ctx context.Context, request *TransportRequestMessage, baseURL, apiKey string) *queue.TransportResponseMessage {
	trackingNumber, _ := request.Payload["tracking_number"].(string)
	if trackingNumber == "" {
		return c.errorResponse(request, "tracking_number is required in payload")
	}

	resp, err := c.useCase.Track(ctx, baseURL, apiKey, trackingNumber)
	if err != nil {
		return c.errorResponse(request, err.Error())
	}

	return &queue.TransportResponseMessage{
		ShipmentID:    request.ShipmentID,
		BusinessID:    request.BusinessID,
		Provider:      "envioclick",
		Operation:     "track",
		Status:        "success",
		CorrelationID: request.CorrelationID,
		IsTest:        request.IsTest,
		Timestamp:     time.Now(),
		Data:          toMap(resp),
	}
}

// processCancel handles shipment cancellation
func (c *TransportRequestConsumer) processCancel(ctx context.Context, request *TransportRequestMessage, baseURL, apiKey string) *queue.TransportResponseMessage {
	idShipment, _ := request.Payload["id_shipment"].(string)
	if idShipment == "" {
		return c.errorResponse(request, "id_shipment is required in payload")
	}

	resp, err := c.useCase.Cancel(ctx, baseURL, apiKey, idShipment)
	if err != nil {
		return c.errorResponse(request, err.Error())
	}

	return &queue.TransportResponseMessage{
		ShipmentID:    request.ShipmentID,
		BusinessID:    request.BusinessID,
		Provider:      "envioclick",
		Operation:     "cancel",
		Status:        "success",
		CorrelationID: request.CorrelationID,
		IsTest:        request.IsTest,
		Timestamp:     time.Now(),
		Data:          toMap(resp),
	}
}

// errorResponse creates an error response
func (c *TransportRequestConsumer) errorResponse(request *TransportRequestMessage, errMsg string) *queue.TransportResponseMessage {
	return &queue.TransportResponseMessage{
		ShipmentID:    request.ShipmentID,
		BusinessID:    request.BusinessID,
		Provider:      "envioclick",
		Operation:     request.Operation,
		Status:        "error",
		CorrelationID: request.CorrelationID,
		Timestamp:     time.Now(),
		Error:         errMsg,
	}
}

// toMap converts a struct to map[string]interface{} via JSON
func toMap(v interface{}) map[string]interface{} {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}
