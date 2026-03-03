package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const (
	QueueTransportResponses = rabbitmq.QueueTransportResponses
)

// TransportResponseMessage is the response message from a transport provider
// (replicated locally from integrations/transport)
type TransportResponseMessage struct {
	ShipmentID    *uint                  `json:"shipment_id,omitempty"`
	BusinessID    uint                   `json:"business_id"`
	Provider      string                 `json:"provider"`
	Operation     string                 `json:"operation"`
	Status        string                 `json:"status"` // "success", "error"
	CorrelationID string                 `json:"correlation_id"`
	IsTest        bool                   `json:"is_test,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data,omitempty"`
	Error         string                 `json:"error,omitempty"`
}

// ResponseConsumer consumes transport responses
type ResponseConsumer struct {
	queue        rabbitmq.IQueue
	repo         domain.IRepository
	log          log.ILogger
	ssePublisher domain.IShipmentSSEPublisher
	redisClient  redis.IRedis
}

// NewResponseConsumer creates a new transport response consumer
func NewResponseConsumer(
	queue rabbitmq.IQueue,
	repo domain.IRepository,
	logger log.ILogger,
	ssePublisher domain.IShipmentSSEPublisher,
	redisClient redis.IRedis,
) *ResponseConsumer {
	return &ResponseConsumer{
		queue:        queue,
		repo:         repo,
		log:          logger.WithModule("shipments.transport_response_consumer"),
		ssePublisher: ssePublisher,
		redisClient:  redisClient,
	}
}

// Start begins consuming transport responses
func (c *ResponseConsumer) Start(ctx context.Context) error {
	if err := c.queue.DeclareQueue(QueueTransportResponses, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al declarar cola de transport responses")
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	c.log.Info(ctx).
		Str("queue", QueueTransportResponses).
		Msg("📥 Starting transport response consumer")

	if err := c.queue.Consume(ctx, QueueTransportResponses, c.handleResponse); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al iniciar consumer de transport responses")
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	return nil
}

// handleResponse processes a transport response
func (c *ResponseConsumer) handleResponse(message []byte) error {
	ctx := context.Background()

	var response TransportResponseMessage
	if err := json.Unmarshal(message, &response); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al deserializar transport response")
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.log.Info(ctx).
		Str("provider", response.Provider).
		Str("operation", response.Operation).
		Str("status", response.Status).
		Str("correlation_id", response.CorrelationID).
		Msg("📨 Processing transport response")

	switch response.Operation {
	case "quote":
		c.handleQuoteResponse(ctx, &response)
	case "generate":
		c.handleGenerateResponse(ctx, &response)
	case "track":
		c.handleTrackResponse(ctx, &response)
	case "cancel":
		c.handleCancelResponse(ctx, &response)
	default:
		c.log.Warn(ctx).
			Str("operation", response.Operation).
			Msg("Unknown transport operation in response")
	}

	return nil
}

// handleGenerateResponse processes a guide generation response
func (c *ResponseConsumer) handleGenerateResponse(ctx context.Context, response *TransportResponseMessage) {
	businessID := c.resolveBusinessID(ctx, response)

	if response.Status == "error" {
		c.log.Error(ctx).
			Str("error", response.Error).
			Str("correlation_id", response.CorrelationID).
			Msg("❌ Guide generation failed")

		// If we have a shipment ID, update status to failed
		if response.ShipmentID != nil {
			shipment, err := c.repo.GetShipmentByID(ctx, *response.ShipmentID)
			if err == nil && shipment != nil {
				shipment.Status = "failed"
				if err := c.repo.UpdateShipment(ctx, shipment); err != nil {
					c.log.Error(ctx).Err(err).Msg("Failed to update shipment status to failed")
				}
			}
			c.ssePublisher.PublishGuideFailed(ctx, businessID, *response.ShipmentID, response.CorrelationID, response.Error)
		}
		return
	}

	// Success: extract tracking data from response
	data := response.Data
	if data == nil {
		c.log.Warn(ctx).Msg("Generate response has no data")
		return
	}

	// Extract nested data field
	dataField, _ := data["data"].(map[string]interface{})
	if dataField == nil {
		dataField = data
	}

	trackingNumber, _ := dataField["tracker"].(string)
	labelURL, _ := dataField["url"].(string)

	c.log.Info(ctx).
		Str("tracking_number", trackingNumber).
		Str("label_url", labelURL).
		Str("correlation_id", response.CorrelationID).
		Msg("✅ Guide generated successfully")

	// If we have a shipment ID, update the shipment
	if response.ShipmentID != nil {
		shipment, err := c.repo.GetShipmentByID(ctx, *response.ShipmentID)
		if err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to get shipment")
			return
		}
		if shipment != nil {
			if trackingNumber != "" {
				shipment.TrackingNumber = &trackingNumber
			}
			if labelURL != "" {
				shipment.GuideURL = &labelURL
			}
			shipment.Status = "pending"
			shipment.IsTest = response.IsTest

			if err := c.repo.UpdateShipment(ctx, shipment); err != nil {
				c.log.Error(ctx).Err(err).Msg("Failed to update shipment with tracking data")
			}

			// Sync guide_link and tracking_number to the order immediately
			if shipment.OrderID != nil && *shipment.OrderID != "" {
				if err := c.repo.UpdateOrderGuideLink(ctx, *shipment.OrderID, labelURL, trackingNumber); err != nil {
					c.log.Error(ctx).Err(err).
						Str("order_id", *shipment.OrderID).
						Msg("Failed to sync guide_link to order")
				} else {
					c.log.Info(ctx).
						Str("order_id", *shipment.OrderID).
						Str("guide_link", labelURL).
						Msg("✅ guide_link synced to order")
				}
			}
		}

		c.ssePublisher.PublishGuideGenerated(ctx, businessID, *response.ShipmentID, response.CorrelationID, trackingNumber, labelURL)
	}
}

// handleQuoteResponse processes a quote response
func (c *ResponseConsumer) handleQuoteResponse(ctx context.Context, response *TransportResponseMessage) {
	businessID := response.BusinessID

	if response.Status == "error" {
		c.log.Error(ctx).
			Str("error", response.Error).
			Str("correlation_id", response.CorrelationID).
			Msg("❌ Quote request failed")

		c.storeQuoteResult(ctx, response.CorrelationID, nil, response.Error)
		c.ssePublisher.PublishQuoteFailed(ctx, businessID, response.CorrelationID, response.Error)
		return
	}

	c.log.Info(ctx).
		Str("correlation_id", response.CorrelationID).
		Msg("✅ Quote response received")

	c.storeQuoteResult(ctx, response.CorrelationID, response.Data, "")
	c.ssePublisher.PublishQuoteReceived(ctx, businessID, response.CorrelationID, response.Data)
}

// storeQuoteResult stores the quote result in Redis so the HTTP handler can poll for it synchronously.
func (c *ResponseConsumer) storeQuoteResult(ctx context.Context, correlationID string, data map[string]interface{}, errMsg string) {
	if c.redisClient == nil {
		return
	}

	status := "success"
	if errMsg != "" {
		status = "error"
	}

	result := map[string]interface{}{
		"status": status,
		"data":   data,
		"error":  errMsg,
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		c.log.Warn(ctx).Err(err).Str("correlation_id", correlationID).Msg("Failed to marshal quote result for Redis")
		return
	}

	key := fmt.Sprintf("shipment:quote:result:%s", correlationID)
	if err := c.redisClient.Set(ctx, key, string(bytes), 60*time.Second); err != nil {
		c.log.Warn(ctx).Err(err).Str("correlation_id", correlationID).Msg("Failed to store quote result in Redis")
	}
}

// handleTrackResponse processes a tracking response
func (c *ResponseConsumer) handleTrackResponse(ctx context.Context, response *TransportResponseMessage) {
	businessID := response.BusinessID

	if response.Status == "error" {
		c.log.Error(ctx).
			Str("error", response.Error).
			Str("correlation_id", response.CorrelationID).
			Msg("❌ Tracking request failed")

		c.ssePublisher.PublishTrackingFailed(ctx, businessID, response.CorrelationID, response.Error)
		return
	}

	c.log.Info(ctx).
		Str("correlation_id", response.CorrelationID).
		Msg("✅ Tracking response received")

	c.ssePublisher.PublishTrackingUpdated(ctx, businessID, response.CorrelationID, response.Data)
}

// handleCancelResponse processes a cancellation response
func (c *ResponseConsumer) handleCancelResponse(ctx context.Context, response *TransportResponseMessage) {
	businessID := c.resolveBusinessID(ctx, response)

	if response.Status == "error" {
		c.log.Error(ctx).
			Str("error", response.Error).
			Str("correlation_id", response.CorrelationID).
			Msg("❌ Shipment cancellation failed")

		shipmentID := uint(0)
		if response.ShipmentID != nil {
			shipmentID = *response.ShipmentID
		}
		c.ssePublisher.PublishCancelFailed(ctx, businessID, shipmentID, response.CorrelationID, response.Error)
		return
	}

	c.log.Info(ctx).
		Str("correlation_id", response.CorrelationID).
		Msg("✅ Shipment cancelled successfully")

	// If we have a shipment ID, update status
	if response.ShipmentID != nil {
		shipment, err := c.repo.GetShipmentByID(ctx, *response.ShipmentID)
		if err == nil && shipment != nil {
			shipment.Status = "cancelled"
			if err := c.repo.UpdateShipment(ctx, shipment); err != nil {
				c.log.Error(ctx).Err(err).Msg("Failed to update shipment status to cancelled")
			}
		}

		c.ssePublisher.PublishShipmentCancelled(ctx, businessID, *response.ShipmentID)
	}
}

// resolveBusinessID resolves the business ID from the response message.
// The transport router should always echo back business_id from the original request.
func (c *ResponseConsumer) resolveBusinessID(ctx context.Context, response *TransportResponseMessage) uint {
	if response.BusinessID != 0 {
		return response.BusinessID
	}

	c.log.Warn(ctx).
		Str("correlation_id", response.CorrelationID).
		Str("operation", response.Operation).
		Msg("Transport response missing business_id")

	return 0
}
