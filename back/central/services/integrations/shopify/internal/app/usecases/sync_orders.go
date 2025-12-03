package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

type SyncOrdersUseCase struct {
	coreIntegration core.IIntegrationCore
	shopifyClient   domain.ShopifyClient
	publisher       domain.OrderPublisher
}

func New(
	coreIntegration core.IIntegrationCore,
	shopifyClient domain.ShopifyClient,
	publisher domain.OrderPublisher,
) *SyncOrdersUseCase {
	return &SyncOrdersUseCase{
		coreIntegration: coreIntegration,
		shopifyClient:   shopifyClient,
		publisher:       publisher,
	}
}

func (uc *SyncOrdersUseCase) Execute(ctx context.Context, businessID *uint, createdMin *time.Time) error {
	// 1. Get Integration credentials
	integration, err := uc.coreIntegration.GetIntegrationByType(ctx, "shopify", businessID)
	if err != nil {
		return fmt.Errorf("failed to get shopify integration: %w", err)
	}

	if !integration.IsActive {
		return fmt.Errorf("integration is not active")
	}

	// Extract credentials
	var credsMap map[string]interface{}
	if len(integration.Credentials) > 0 {
		if err := json.Unmarshal(integration.Credentials, &credsMap); err != nil {
			return fmt.Errorf("failed to unmarshal credentials: %w", err)
		}
	}

	accessToken, ok := credsMap["access_token"].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("access token not found in credentials")
	}

	// Extract config
	var config map[string]interface{}
	if len(integration.Config) > 0 {
		_ = json.Unmarshal(integration.Config, &config)
	}
	storeName, _ := config["store_name"].(string)
	if storeName == "" {
		return fmt.Errorf("store name not found in config")
	}

	// 2. Prepare params
	params := map[string]string{
		"status": "any",
		"limit":  "250",
	}
	if createdMin != nil {
		params["created_at_min"] = createdMin.Format(time.RFC3339)
	}

	// 3. Fetch orders with pagination
	for {
		ordersData, nextPageURL, err := uc.shopifyClient.FetchOrders(ctx, storeName, accessToken, params)
		if err != nil {
			return fmt.Errorf("failed to fetch orders: %w", err)
		}

		for _, data := range ordersData {
			// Map to Unified Order
			unifiedOrder, err := uc.mapToUnified(integration, data)
			if err != nil {
				// Log error but continue?
				continue
			}

			// Publish to queue
			if err := uc.publisher.Publish(ctx, unifiedOrder); err != nil {
				return fmt.Errorf("failed to publish order: %w", err)
			}
		}

		if nextPageURL == "" {
			break
		}
		// TODO: Handle pagination properly with next page URL
		break
	}

	return nil
}

func (uc *SyncOrdersUseCase) mapToUnified(integration *core.IntegrationWithCredentials, data map[string]interface{}) (*domain.UnifiedOrder, error) {
	// Helper to safely get string
	getString := func(m map[string]interface{}, key string) string {
		if v, ok := m[key].(string); ok {
			return v
		}
		return ""
	}

	idVal := data["id"]
	externalID := fmt.Sprintf("%v", idVal)
	if f, ok := idVal.(float64); ok {
		externalID = fmt.Sprintf("%.0f", f)
	}

	orderNumber := fmt.Sprintf("%v", data["order_number"])

	var totalAmount float64
	if tp, ok := data["total_price"].(string); ok {
		totalAmount, _ = strconv.ParseFloat(tp, 64)
	}

	createdAtStr, _ := data["created_at"].(string)
	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)

	shipping, _ := data["shipping_address"].(map[string]interface{})

	// Map Items
	var items []domain.UnifiedOrderItem
	if lineItems, ok := data["line_items"].([]interface{}); ok {
		for _, item := range lineItems {
			if itemMap, ok := item.(map[string]interface{}); ok {
				price, _ := strconv.ParseFloat(getString(itemMap, "price"), 64)
				qty := 0
				if q, ok := itemMap["quantity"].(float64); ok {
					qty = int(q)
				}

				items = append(items, domain.UnifiedOrderItem{
					ExternalID: fmt.Sprintf("%v", itemMap["id"]),
					Name:       getString(itemMap, "name"),
					SKU:        getString(itemMap, "sku"),
					Quantity:   qty,
					UnitPrice:  price,
				})
			}
		}
	}

	// Extraer el tipo de integración del IntegrationType relacionado
	integrationType := "shopify" // Default
	if integration.IntegrationType != nil {
		integrationType = integration.IntegrationType.Code
	}

	return &domain.UnifiedOrder{
		// Identificadores de integración
		BusinessID:      integration.BusinessID,
		IntegrationID:   integration.ID,
		IntegrationType: integrationType,

		// Identificadores de orden
		Platform:    "shopify",
		ExternalID:  externalID,
		OrderNumber: orderNumber,

		// Información financiera
		TotalAmount: totalAmount,
		Currency:    getString(data, "currency"),

		// Cliente y envío
		Customer: domain.UnifiedCustomer{
			Name:  getString(shipping, "name"),
			Email: getString(data, "email"),
			Phone: getString(shipping, "phone"),
		},
		ShippingAddress: domain.UnifiedAddress{
			Street:     getString(shipping, "address1"),
			City:       getString(shipping, "city"),
			State:      getString(shipping, "province"),
			Country:    getString(shipping, "country"),
			PostalCode: getString(shipping, "zip"),
		},

		// Estado
		Status:         getString(data, "financial_status"),
		OriginalStatus: getString(data, "financial_status"),

		// Items
		Items: items,

		// Metadata
		Metadata: data,

		// Timestamps
		OccurredAt: createdAt,
		ImportedAt: time.Now(),
	}, nil
}
