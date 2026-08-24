package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

type fulfillmentOrderResponse struct {
	ID     string          `json:"id"`
	Number json.Number     `json:"number"`
	Status json.RawMessage `json:"status"`
}

type fulfillmentStatusObject struct {
	Status string `json:"status"`
}

func parseFulfillmentStatus(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var texto string
	if err := json.Unmarshal(raw, &texto); err == nil {
		return strings.ToUpper(strings.TrimSpace(texto))
	}
	var objeto fulfillmentStatusObject
	if err := json.Unmarshal(raw, &objeto); err == nil {
		return strings.ToUpper(strings.TrimSpace(objeto.Status))
	}
	return ""
}

func (c *TiendanubeClient) ListFulfillmentOrders(ctx context.Context, cred domain.Credential, orderID string) ([]domain.FulfillmentOrder, error) {
	raw, _, err := c.do(ctx, cred, http.MethodGet, "/orders/"+orderID+"/fulfillment-orders", nil, nil)
	if err != nil {
		return nil, err
	}

	var respuesta []fulfillmentOrderResponse
	if err := json.Unmarshal(raw, &respuesta); err != nil {
		return nil, fmt.Errorf("tiendanube client: parsing fulfillment orders: %w", err)
	}

	ordenes := make([]domain.FulfillmentOrder, 0, len(respuesta))
	for _, item := range respuesta {
		numero, _ := item.Number.Int64()
		ordenes = append(ordenes, domain.FulfillmentOrder{
			ID:     item.ID,
			Number: int(numero),
			Status: parseFulfillmentStatus(item.Status),
		})
	}
	return ordenes, nil
}

func (c *TiendanubeClient) UpdateFulfillmentOrder(ctx context.Context, cred domain.Credential, orderID, fulfillmentOrderID, status string, tracking *domain.TrackingInfo) error {
	body := map[string]interface{}{"status": status}
	if tracking != nil && strings.TrimSpace(tracking.Code) != "" {
		info := map[string]interface{}{
			"code":            tracking.Code,
			"notify_customer": tracking.NotifyCustomer,
		}
		if strings.TrimSpace(tracking.URL) != "" {
			info["url"] = tracking.URL
		}
		body["tracking_info"] = info
	}

	_, _, err := c.do(ctx, cred, http.MethodPatch, "/orders/"+orderID+"/fulfillment-orders/"+fulfillmentOrderID, nil, body)
	return err
}

func (c *TiendanubeClient) CreateTrackingEvent(ctx context.Context, cred domain.Credential, orderID, fulfillmentOrderID string, event domain.TrackingEvent) error {
	body := map[string]interface{}{"status": event.Status}
	if strings.TrimSpace(event.Description) != "" {
		body["description"] = event.Description
	}
	if strings.TrimSpace(event.Address) != "" {
		body["address"] = event.Address
	}
	if !event.HappenedAt.IsZero() {
		body["happened_at"] = event.HappenedAt.UTC().Format(time.RFC3339)
	}
	if event.EstimatedDeliveryAt != nil && !event.EstimatedDeliveryAt.IsZero() {
		body["estimated_delivery_at"] = event.EstimatedDeliveryAt.UTC().Format(time.RFC3339)
	}

	_, _, err := c.do(ctx, cred, http.MethodPost, "/orders/"+orderID+"/fulfillment-orders/"+fulfillmentOrderID+"/tracking-events", nil, body)
	return err
}

func (c *TiendanubeClient) CancelOrder(ctx context.Context, cred domain.Credential, orderID string) error {
	_, _, err := c.do(ctx, cred, http.MethodPost, "/orders/"+orderID+"/cancel", nil, nil)
	return err
}
