package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

type ResponseConsumer struct {
	queue        rabbitmq.IQueue
	repo         domain.IRepository
	log          log.ILogger
	ssePublisher domain.IShipmentSSEPublisher
	redisClient  redis.IRedis
	marginReader domain.IShippingMarginReader
}

func NewResponseConsumer(
	queue rabbitmq.IQueue,
	repo domain.IRepository,
	logger log.ILogger,
	ssePublisher domain.IShipmentSSEPublisher,
	redisClient redis.IRedis,
	marginReader domain.IShippingMarginReader,
) *ResponseConsumer {
	return &ResponseConsumer{
		queue:        queue,
		repo:         repo,
		log:          logger.WithModule("shipments.transport_response_consumer"),
		ssePublisher: ssePublisher,
		redisClient:  redisClient,
		marginReader: marginReader,
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
	case "webhook_update":
		c.handleWebhookUpdate(ctx, &response)
	case "sync_batch":
		c.log.Info(ctx).
			Str("correlation_id", response.CorrelationID).
			Interface("summary", response.Data).
			Msg("Sync batch summary received")
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

	c.log.Info(ctx).Interface("response_data_keys", getKeys(data)).Msg("DEBUG: Response.Data keys")

	// Extract nested data field
	dataField, _ := data["data"].(map[string]interface{})
	if dataField == nil {
		dataField = data
		c.log.Info(ctx).Msg("DEBUG: Using data directly as dataField (no nested 'data' field)")
	} else {
		c.log.Info(ctx).Interface("dataField_keys", getKeys(dataField)).Msg("DEBUG: Using nested dataField")
	}

	trackingNumber, _ := dataField["tracker"].(string)
	labelURL, _ := dataField["url"].(string)
	carrier, _ := dataField["carrier"].(string)
	idOrder, _ := dataField["idOrder"].(float64)

	// Si el carrier viene vacío de la respuesta, inferirlo del tracking_number
	if carrier == "" && trackingNumber != "" {
		carrier = inferCarrierFromTrackingNumber(trackingNumber)
	}

	c.log.Info(ctx).
		Str("tracking_number", trackingNumber).
		Str("label_url", labelURL).
		Str("carrier", carrier).
		Float64("id_order", idOrder).
		Str("correlation_id", response.CorrelationID).
		Interface("all_datafield_values", dataField).
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
			if carrier != "" {
				shipment.Carrier = &carrier
			}

			// Store idOrder in Metadata
			if idOrder != 0 {
				var meta map[string]interface{}
				if len(shipment.Metadata) > 0 {
					json.Unmarshal(shipment.Metadata, &meta)
				}
				if meta == nil {
					meta = make(map[string]interface{})
				}
				meta["envioclick_id_order"] = int64(idOrder)
				if updatedBytes, err := json.Marshal(meta); err == nil {
					shipment.Metadata = updatedBytes
				}
			}

			shipment.Status = "pending"
			shipment.IsTest = response.IsTest

			appendGuideGeneratedEvent(shipment, response.Provider, trackingNumber, carrier)

			if err := c.repo.UpdateShipment(ctx, shipment); err != nil {
				c.log.Error(ctx).Err(err).Msg("Failed to update shipment with tracking data")
			}

			if shipment.TotalCost != nil && *shipment.TotalCost > 0 && businessID != 0 {
				if err := c.repo.DebitWalletForGuide(ctx, businessID, *shipment.TotalCost, trackingNumber); err != nil {
					c.log.Error(ctx).Err(err).
						Uint("business_id", businessID).
						Float64("amount", *shipment.TotalCost).
						Str("tracking_number", trackingNumber).
						Msg("Failed to debit wallet for guide")
				} else {
					c.log.Info(ctx).
						Uint("business_id", businessID).
						Float64("amount", *shipment.TotalCost).
						Str("tracking_number", trackingNumber).
						Msg("Wallet debited for guide")
				}
			}

			// Sync guide_link, tracking_number, and carrier to the order immediately
			if shipment.OrderID != nil && *shipment.OrderID != "" {
				if err := c.repo.UpdateOrderGuideLink(ctx, *shipment.OrderID, labelURL, trackingNumber, carrier); err != nil {
					c.log.Error(ctx).Err(err).
						Str("order_id", *shipment.OrderID).
						Msg("Failed to sync guide_link to order")
				} else {
					c.log.Info(ctx).
						Str("order_id", *shipment.OrderID).
						Str("guide_link", labelURL).
						Str("carrier", carrier).
						Msg("✅ guide_link and carrier synced to order")
				}
			}
		}

		// Use shipment's existing carrier as fallback when provider response has empty carrier
		effectiveCarrier := carrier
		if effectiveCarrier == "" && shipment != nil && shipment.Carrier != nil && *shipment.Carrier != "" {
			effectiveCarrier = *shipment.Carrier
		}

		// Enrich with customer/business/integration data for WhatsApp notifications
		var notification *domain.GuideNotificationData
		if shipment != nil {
			notification = &domain.GuideNotificationData{
				CustomerName:  shipment.CustomerName,
				CustomerPhone: shipment.CustomerPhone,
				OrderNumber:   shipment.OrderNumber,
				CodTotal:      shipment.CodTotal,
			}
			if trackingNumber != "" {
				url := "https://www.probabilityia.com.co/rastreo?tracking=" + trackingNumber
				if businessID > 0 {
					url += "&b=" + strconv.FormatUint(uint64(businessID), 10)
				}
				notification.TrackingURL = url
			}

			// Fetch business name
			if businessID != 0 {
				if bName, err := c.repo.GetBusinessName(ctx, businessID); err == nil {
					notification.BusinessName = bName
				}
			}

			// Fetch integration_id from the order
			if shipment.OrderID != nil && *shipment.OrderID != "" {
				if intID, err := c.repo.GetOrderIntegrationID(ctx, *shipment.OrderID); err == nil {
					notification.IntegrationID = intID
				}
			}
		}

		c.ssePublisher.PublishGuideGenerated(ctx, businessID, *response.ShipmentID, response.CorrelationID, trackingNumber, labelURL, effectiveCarrier, notification)
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

	c.applyServiceFeeToQuoteData(ctx, response.Data, response.Provider, businessID)

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

	// Update shipment status based on tracking data from carrier
	if response.ShipmentID != nil && response.Data != nil {
		shipment, err := c.repo.GetShipmentByID(ctx, *response.ShipmentID)
		if err == nil && shipment != nil {
			status, ok := response.Data["status"].(string)
			history, hasHistory := response.Data["history"]

			if hasHistory {
				var meta map[string]interface{}
				if len(shipment.Metadata) > 0 {
					json.Unmarshal(shipment.Metadata, &meta)
				}
				if meta == nil {
					meta = make(map[string]interface{})
				}
				meta["tracking_events"] = history
				if updatedBytes, err := json.Marshal(meta); err == nil {
					shipment.Metadata = updatedBytes
				}
			}

			if (ok && status != "") || hasHistory {
				if ok && status != "" {
					shipment.Status = status // Update to: in_transit, delivered, failed, etc.
				}

				if err := c.repo.UpdateShipment(ctx, shipment); err != nil {
					c.log.Error(ctx).
						Err(err).
						Str("shipment_id", fmt.Sprintf("%d", *response.ShipmentID)).
						Msg("Failed to update shipment status/history from tracking")
				} else {
					c.log.Info(ctx).
						Str("shipment_id", fmt.Sprintf("%d", *response.ShipmentID)).
						Msg("✅ Shipment status/history updated from tracking response")
				}
			}
		}
	}

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

	if response.ShipmentID != nil {
		shipment, err := c.repo.GetShipmentByID(ctx, *response.ShipmentID)
		if err == nil && shipment != nil {
			shipment.Status = "cancelled"
			appendCancelEvent(shipment, response.Provider)
			if err := c.repo.UpdateShipment(ctx, shipment); err != nil {
				c.log.Error(ctx).Err(err).Msg("Failed to update shipment status to cancelled")
			}

			if shipment.OrderID != nil && *shipment.OrderID != "" {
				if err := c.repo.ClearOrderGuideData(ctx, *shipment.OrderID); err != nil {
					c.log.Error(ctx).Err(err).
						Str("order_id", *shipment.OrderID).
						Msg("Failed to clear guide data from order")
				}
			}
		}

		c.ssePublisher.PublishShipmentCancelled(ctx, businessID, *response.ShipmentID)
	}
}

func (c *ResponseConsumer) handleWebhookUpdate(ctx context.Context, response *TransportResponseMessage) {
	if response.Data == nil {
		c.log.Warn(ctx).Str("correlation_id", response.CorrelationID).Msg("webhook_update response has no data")
		return
	}

	trackingNumber, _ := response.Data["tracking_number"].(string)
	probabilityStatus, _ := response.Data["probability_status"].(string)

	if trackingNumber == "" || probabilityStatus == "" {
		c.log.Warn(ctx).
			Str("correlation_id", response.CorrelationID).
			Str("tracking_number", trackingNumber).
			Str("probability_status", probabilityStatus).
			Msg("webhook_update response missing required fields")
		return
	}

	shipment, err := c.repo.GetShipmentByTrackingNumber(ctx, trackingNumber)
	if err != nil || shipment == nil {
		c.log.Warn(ctx).
			Err(err).
			Str("tracking_number", trackingNumber).
			Str("correlation_id", response.CorrelationID).
			Msg("Shipment not found for webhook update")
		return
	}

	previousStatus := shipment.Status
	shipment.Status = probabilityStatus

	if shippedAtStr, ok := response.Data["shipped_at"].(string); ok && shippedAtStr != "" {
		if t := parseFlexibleTime(shippedAtStr); t != nil {
			shipment.ShippedAt = t
		}
	}

	if deliveredAtStr, ok := response.Data["delivered_at"].(string); ok && deliveredAtStr != "" {
		if t := parseFlexibleTime(deliveredAtStr); t != nil {
			shipment.DeliveredAt = t
		}
	}

	appendTrackingEvent(shipment, response, probabilityStatus)

	if err := c.repo.UpdateShipment(ctx, shipment); err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("shipment_id", shipment.ID).
			Str("correlation_id", response.CorrelationID).
			Msg("Failed to update shipment from webhook")
		return
	}

	if previousStatus != probabilityStatus && shipment.OrderID != nil && *shipment.OrderID != "" {
		if err := c.repo.UpdateOrderStatusByOrderID(ctx, *shipment.OrderID, probabilityStatus); err != nil {
			c.log.Warn(ctx).
				Err(err).
				Str("order_id", *shipment.OrderID).
				Str("new_status", probabilityStatus).
				Msg("Failed to sync order status from webhook")
		}
	}

	c.log.Info(ctx).
		Uint("shipment_id", shipment.ID).
		Str("tracking_number", trackingNumber).
		Str("previous_status", previousStatus).
		Str("new_status", probabilityStatus).
		Str("provider", response.Provider).
		Str("correlation_id", response.CorrelationID).
		Msg("✅ Shipment updated from provider webhook")

	businessID, _ := c.repo.GetShipmentBusinessIDByID(ctx, shipment.ID)

	response.Data["shipment_id"] = shipment.ID
	response.Data["previous_status"] = previousStatus
	response.Data["new_status"] = probabilityStatus
	if shipment.OrderID != nil {
		response.Data["order_id"] = *shipment.OrderID
	}
	if shipment.OrderNumber != "" {
		response.Data["order_number"] = shipment.OrderNumber
	}
	if shipment.CustomerName != "" {
		response.Data["customer_name"] = shipment.CustomerName
	}

	c.ssePublisher.PublishTrackingUpdated(ctx, businessID, response.CorrelationID, response.Data)
}

func appendCancelEvent(shipment *domain.Shipment, provider string) {
	event := map[string]any{
		"date":        time.Now().Format(time.RFC3339),
		"status":      "cancelled",
		"raw_status":  "Cancelado",
		"description": "Envío cancelado",
		"source":      provider,
	}
	mergeTrackingEvent(shipment, event)
}

func appendGuideGeneratedEvent(shipment *domain.Shipment, provider, trackingNumber, carrier string) {
	event := map[string]any{
		"date":        time.Now().Format(time.RFC3339),
		"status":      "pending",
		"raw_status":  "Guía generada",
		"description": fmt.Sprintf("Guía creada con %s (tracking: %s)", carrierOrDefault(carrier, provider), trackingNumber),
		"source":      provider,
	}
	mergeTrackingEvent(shipment, event)
}

func carrierOrDefault(carrier, provider string) string {
	if carrier != "" {
		return carrier
	}
	return provider
}

func mergeTrackingEvent(shipment *domain.Shipment, newEvent map[string]any) {
	var meta map[string]any
	if len(shipment.Metadata) > 0 {
		_ = json.Unmarshal(shipment.Metadata, &meta)
	}
	if meta == nil {
		meta = make(map[string]any)
	}

	var events []any
	if raw, ok := meta["tracking_events"].([]any); ok {
		events = raw
	}

	for _, existing := range events {
		em, ok := existing.(map[string]any)
		if !ok {
			continue
		}
		if em["raw_status"] == newEvent["raw_status"] && em["description"] == newEvent["description"] && em["date"] == newEvent["date"] {
			return
		}
	}

	events = append(events, newEvent)
	meta["tracking_events"] = events

	if updated, err := json.Marshal(meta); err == nil {
		shipment.Metadata = updated
	}
}

func appendTrackingEvent(shipment *domain.Shipment, response *TransportResponseMessage, probabilityStatus string) {
	rawStatus, _ := response.Data["raw_status"].(string)
	description, _ := response.Data["event_description"].(string)
	eventTimestamp, _ := response.Data["event_timestamp"].(string)
	hasIncidence, _ := response.Data["has_incidence"].(bool)

	if eventTimestamp == "" {
		eventTimestamp = time.Now().Format(time.RFC3339)
	}

	newEvent := map[string]any{
		"date":          eventTimestamp,
		"status":        probabilityStatus,
		"raw_status":    rawStatus,
		"description":   description,
		"has_incidence": hasIncidence,
		"source":        response.Provider,
	}
	mergeTrackingEvent(shipment, newEvent)
}

func parseFlexibleTime(s string) *time.Time {
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", "2006-01-02T15:04:05.000Z", time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
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

func getKeys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func normalizeCarrierCode(carrier string) string {
	c := strings.ToLower(strings.TrimSpace(carrier))
	c = strings.ReplaceAll(c, " ", "")
	return c
}

func toFloat(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case float32:
		return float64(x), true
	default:
		return 0, false
	}
}

func (c *ResponseConsumer) applyServiceFeeToQuoteData(ctx context.Context, data map[string]interface{}, provider string, businessID uint) {
	if data == nil || c.marginReader == nil || businessID == 0 {
		return
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		return
	}
	rawRates, ok := innerData["rates"]
	if !ok {
		return
	}
	rates, ok := rawRates.([]interface{})
	if !ok {
		return
	}

	for _, r := range rates {
		rate, ok := r.(map[string]interface{})
		if !ok {
			continue
		}

		carrierName, _ := rate["carrier"].(string)
		carrierCode := normalizeCarrierCode(carrierName)
		if carrierCode == "" {
			continue
		}

		margin, err := c.marginReader.Get(ctx, businessID, carrierCode)
		if err != nil {
			c.log.Warn(ctx).Err(err).Uint("business_id", businessID).Str("carrier", carrierCode).Msg("shipping margin lookup failed")
			continue
		}
		if margin.MarginAmount == 0 && margin.InsuranceMargin == 0 {
			continue
		}

		if val, exists := rate["flete"]; exists && margin.MarginAmount > 0 {
			if oldVal, ok := toFloat(val); ok {
				newVal := oldVal + margin.MarginAmount
				rate["flete"] = newVal
				c.log.Info(ctx).
					Str("provider", provider).
					Str("carrier", carrierName).
					Uint("business_id", businessID).
					Float64("original_flete", oldVal).
					Float64("margin_amount", margin.MarginAmount).
					Float64("final_flete", newVal).
					Msg("Margin applied to quote flete")
			}
		}

		if margin.InsuranceMargin > 0 {
			if val, exists := rate["minimumInsurance"]; exists {
				if oldVal, ok := toFloat(val); ok {
					newVal := oldVal + margin.InsuranceMargin
					rate["minimumInsurance"] = newVal
					c.log.Info(ctx).
						Str("carrier", carrierName).
						Uint("business_id", businessID).
						Float64("original_insurance", oldVal).
						Float64("insurance_margin", margin.InsuranceMargin).
						Float64("final_insurance", newVal).
						Msg("Insurance margin applied to quote")
				}
			}
		}
	}
}

// inferCarrierFromTrackingNumber deduce la transportadora basándose en el formato del tracking_number
// Usa prefijos conocidos de Envioclik para identificar la transportadora
func inferCarrierFromTrackingNumber(trackingNumber string) string {
	if trackingNumber == "" {
		return ""
	}

	upper := strings.ToUpper(trackingNumber)

	// Prefijos con guion (formato ENV-xxx, CRD-xxx, etc.)
	if strings.HasPrefix(upper, "ENV-") {
		return "ENVIA"
	}
	if strings.HasPrefix(upper, "IRP-") {
		return "INTERRAPIDISIMO"
	}
	if strings.HasPrefix(upper, "CRD-") {
		return "COORDINADORA"
	}
	if strings.HasPrefix(upper, "SRV-") {
		return "SERVIENTREGA"
	}
	if strings.HasPrefix(upper, "TCC-") {
		return "TODOCARGO"
	}

	// Prefijos numéricos (formato Envioclik)
	// ENVIA: 034056
	if strings.HasPrefix(trackingNumber, "034056") {
		return "ENVIA"
	}

	// INTERRAPIDISIMO: 2400 (rango amplio: 240047, 240048, 240050, etc.)
	if strings.HasPrefix(trackingNumber, "2400") {
		return "INTERRAPIDISIMO"
	}

	// COORDINADORA: 4005
	if strings.HasPrefix(trackingNumber, "4005") {
		return "COORDINADORA"
	}

	// SERVIENTREGA: 072
	if strings.HasPrefix(trackingNumber, "072") {
		return "SERVIENTREGA"
	}

	// Si no coincide con ningún prefijo conocido, retornar vacío
	return ""
}
