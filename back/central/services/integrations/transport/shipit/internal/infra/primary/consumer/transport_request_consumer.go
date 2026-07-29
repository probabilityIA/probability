package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type ICredentialResolver interface {
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	GetIntegrationConfig(ctx context.Context, integrationID string) (map[string]interface{}, error)
}

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
	QueueShipitRequests = rabbitmq.QueueTransportShipitRequests
	providerName        = "shipit"
)

type TransportRequestConsumer struct {
	rabbit             rabbitmq.IQueue
	useCase            app.IUseCase
	responsePublisher  *queue.ResponsePublisher
	webhookPublisher   domain.IWebhookResponsePublisher
	credentialResolver ICredentialResolver
	log                log.ILogger
}

func NewTransportRequestConsumer(
	rabbit rabbitmq.IQueue,
	useCase app.IUseCase,
	responsePublisher *queue.ResponsePublisher,
	webhookPublisher domain.IWebhookResponsePublisher,
	credentialResolver ICredentialResolver,
	logger log.ILogger,
) *TransportRequestConsumer {
	return &TransportRequestConsumer{
		rabbit:             rabbit,
		useCase:            useCase,
		responsePublisher:  responsePublisher,
		webhookPublisher:   webhookPublisher,
		credentialResolver: credentialResolver,
		log:                logger.WithModule("transport.shipit.consumer"),
	}
}

func (c *TransportRequestConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		c.log.Warn(ctx).Msg("RabbitMQ client is nil, consumer cannot start")
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", QueueShipitRequests).
		Msg("Starting Shipit transport request consumer")

	if err := c.rabbit.DeclareQueue(QueueShipitRequests, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare queue")
		return err
	}

	if err := c.rabbit.Consume(ctx, QueueShipitRequests, c.handleRequest); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to start consuming")
		return err
	}

	return nil
}

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
		Msg("Received transport request")

	var integrationConfig map[string]interface{}
	if request.IntegrationID != 0 {
		cfg, cfgErr := c.credentialResolver.GetIntegrationConfig(ctx, fmt.Sprintf("%d", request.IntegrationID))
		if cfgErr != nil {
			c.log.Warn(ctx).Err(cfgErr).Uint("integration_id", request.IntegrationID).Msg("Failed to get integration config, using defaults")
		} else {
			integrationConfig = cfg
		}
	}

	creds, err := c.resolveCredentials(ctx, &request)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("correlation_id", request.CorrelationID).
			Msg("Failed to resolve credentials")
		response := c.errorResponse(&request, "Error al resolver credenciales: "+err.Error())
		if pubErr := c.responsePublisher.PublishResponse(ctx, response); pubErr != nil {
			c.log.Error(ctx).Err(pubErr).Msg("Failed to publish error response")
		}
		return err
	}

	baseURL := c.resolveBaseURL(&request, integrationConfig)

	var response *queue.TransportResponseMessage

	switch request.Operation {
	case "quote":
		response = c.processQuote(ctx, &request, baseURL, creds)
	case "generate":
		response = c.processGenerate(ctx, &request, baseURL, creds)
	case "track":
		response = c.processTrack(ctx, &request, baseURL, creds)
	case "cancel", "cancel_batch":
		response = c.errorResponse(&request, domain.ErrCancelNotSupported.Error())
	case "sync_batch":
		response = c.processSyncBatch(ctx, &request, baseURL, creds)
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

func (c *TransportRequestConsumer) resolveCredentials(ctx context.Context, request *TransportRequestMessage) (domain.Credentials, error) {
	if request.IntegrationID != 0 {
		integrationIDStr := fmt.Sprintf("%d", request.IntegrationID)

		email, emailErr := c.credentialResolver.DecryptCredential(ctx, integrationIDStr, "email")
		token, tokenErr := c.credentialResolver.DecryptCredential(ctx, integrationIDStr, "access_token")

		if emailErr == nil && tokenErr == nil && email != "" && token != "" {
			return domain.Credentials{Email: email, AccessToken: token}, nil
		}

		c.log.Warn(ctx).
			Uint("integration_id", request.IntegrationID).
			Msg("Could not decrypt Shipit credentials, trying env fallback")
	}

	email := os.Getenv("SHIPIT_EMAIL")
	token := os.Getenv("SHIPIT_ACCESS_TOKEN")
	if email != "" && token != "" {
		return domain.Credentials{Email: email, AccessToken: token}, nil
	}

	return domain.Credentials{}, fmt.Errorf("no hay credenciales de Shipit para la integracion %d", request.IntegrationID)
}

func (c *TransportRequestConsumer) resolveBaseURL(request *TransportRequestMessage, config map[string]interface{}) string {
	if config != nil {
		if testURL, ok := config["base_url_test"].(string); ok && testURL != "" {
			return testURL
		}
		if customURL, ok := config["base_url"].(string); ok && customURL != "" {
			return customURL
		}
	}
	return request.BaseURL
}

func (c *TransportRequestConsumer) parseGuideRequest(request *TransportRequestMessage) (*domain.GuideRequest, error) {
	payloadBytes, err := json.Marshal(request.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var req domain.GuideRequest
	if err := json.Unmarshal(payloadBytes, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload as GuideRequest: %w", err)
	}
	return &req, nil
}

func (c *TransportRequestConsumer) processQuote(ctx context.Context, request *TransportRequestMessage, baseURL string, creds domain.Credentials) *queue.TransportResponseMessage {
	req, err := c.parseGuideRequest(request)
	if err != nil {
		return c.errorResponse(request, err.Error())
	}

	resp, err := c.useCase.Quote(ctx, baseURL, creds, *req)
	if err != nil {
		return c.errorResponse(request, err.Error())
	}

	return c.successResponse(request, "quote", toMap(resp))
}

func (c *TransportRequestConsumer) processGenerate(ctx context.Context, request *TransportRequestMessage, baseURL string, creds domain.Credentials) *queue.TransportResponseMessage {
	req, err := c.parseGuideRequest(request)
	if err != nil {
		return c.errorResponse(request, err.Error())
	}

	resp, err := c.useCase.Generate(ctx, baseURL, creds, *req, request.IsTest)
	if err != nil {
		return c.errorResponse(request, err.Error())
	}

	if resp.Data.Carrier == "" {
		if carrierFromPayload, ok := request.Payload["carrier"].(string); ok && carrierFromPayload != "" {
			resp.Data.Carrier = carrierFromPayload
		}
	}

	return c.successResponse(request, "generate", toMap(resp))
}

func (c *TransportRequestConsumer) processTrack(ctx context.Context, request *TransportRequestMessage, baseURL string, creds domain.Credentials) *queue.TransportResponseMessage {
	trackingNumber, _ := request.Payload["tracking_number"].(string)
	if trackingNumber == "" {
		return c.errorResponse(request, "tracking_number is required in payload")
	}

	resp, err := c.useCase.Track(ctx, baseURL, creds, trackingNumber)
	if err != nil {
		return c.errorResponse(request, err.Error())
	}

	return c.successResponse(request, "track", toMap(resp))
}

func (c *TransportRequestConsumer) processSyncBatch(ctx context.Context, request *TransportRequestMessage, baseURL string, creds domain.Credentials) *queue.TransportResponseMessage {
	itemsRaw, ok := request.Payload["items"]
	if !ok {
		return c.errorResponse(request, "payload.items missing for sync_batch")
	}

	itemsBytes, err := json.Marshal(itemsRaw)
	if err != nil {
		return c.errorResponse(request, "failed to marshal items: "+err.Error())
	}

	var items []struct {
		ShipmentID     uint   `json:"shipment_id"`
		TrackingNumber string `json:"tracking_number"`
		Carrier        string `json:"carrier"`
	}
	if err := json.Unmarshal(itemsBytes, &items); err != nil {
		return c.errorResponse(request, "failed to unmarshal items: "+err.Error())
	}

	total := len(items)
	processed := 0
	failed := 0
	notFound := 0

	for _, item := range items {
		if item.TrackingNumber == "" {
			notFound++
			continue
		}

		output, trackErr := c.useCase.Track(ctx, baseURL, creds, item.TrackingNumber)
		if trackErr != nil {
			failed++
			c.log.Warn(ctx).
				Err(trackErr).
				Str("tracking_number", item.TrackingNumber).
				Msg("Shipit sync_batch item failed")
			continue
		}

		shipmentID := item.ShipmentID
		carrier := output.Carrier
		if carrier == "" {
			carrier = item.Carrier
		}
		if carrier == "" {
			carrier = providerName
		}

		msg := &domain.WebhookUpdateMessage{
			ShipmentID:      &shipmentID,
			BusinessID:      request.BusinessID,
			CorrelationID:   request.CorrelationID,
			TrackingNumber:  item.TrackingNumber,
			Carrier:         carrier,
			Status:          domain.ProbabilityShipmentStatus(output.Status),
			RawStatus:       output.StatusDetail,
			RawStatusDetail: output.StatusDetail,
			Description:     output.StatusDetail,
		}

		if pubErr := c.webhookPublisher.PublishWebhookUpdate(ctx, msg); pubErr != nil {
			failed++
			continue
		}
		processed++
	}

	return &queue.TransportResponseMessage{
		BusinessID:    request.BusinessID,
		Provider:      providerName,
		Operation:     "sync_batch",
		Status:        "success",
		CorrelationID: request.CorrelationID,
		IsTest:        request.IsTest,
		Timestamp:     time.Now(),
		Data: map[string]any{
			"total":     total,
			"processed": processed,
			"failed":    failed,
			"not_found": notFound,
		},
	}
}

func (c *TransportRequestConsumer) successResponse(request *TransportRequestMessage, operation string, data map[string]interface{}) *queue.TransportResponseMessage {
	return &queue.TransportResponseMessage{
		ShipmentID:    request.ShipmentID,
		BusinessID:    request.BusinessID,
		Provider:      providerName,
		Operation:     operation,
		Status:        "success",
		CorrelationID: request.CorrelationID,
		IsTest:        request.IsTest,
		Timestamp:     time.Now(),
		Data:          data,
	}
}

func (c *TransportRequestConsumer) errorResponse(request *TransportRequestMessage, errMsg string) *queue.TransportResponseMessage {
	return &queue.TransportResponseMessage{
		ShipmentID:    request.ShipmentID,
		BusinessID:    request.BusinessID,
		Provider:      providerName,
		Operation:     request.Operation,
		Status:        "error",
		CorrelationID: request.CorrelationID,
		Timestamp:     time.Now(),
		Error:         errMsg,
	}
}

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
