package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

type SyncOrdersUseCase struct {
	coreIntegration core.IIntegrationCore
	shopifyClient   domain.ShopifyClient
	orderPublisher  domain.OrderPublisher
}

func New(
	coreIntegration core.IIntegrationCore,
	shopifyClient domain.ShopifyClient,
	orderPublisher domain.OrderPublisher,
) *SyncOrdersUseCase {
	return &SyncOrdersUseCase{
		coreIntegration: coreIntegration,
		shopifyClient:   shopifyClient,
		orderPublisher:  orderPublisher,
	}
}

// SyncOrders fetches orders from Shopify and publishes them to the queue
// This runs in a goroutine for background processing
func (uc *SyncOrdersUseCase) SyncOrders(ctx context.Context, integrationID string) error {
	// Get integration from core
	integration, err := uc.coreIntegration.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	// Extract Shopify credentials from config
	config, ok := integration.Config.(map[string]interface{})
	if !ok || config == nil {
		fmt.Printf("[SyncOrders] Config is nil or invalid for integration %s\n", integrationID)
		// Fallback: try to use Name as store_name if config is empty
		if integration.Name != "" {
			fmt.Printf("[SyncOrders] Using integration Name as store_name: %s\n", integration.Name)
			config = map[string]interface{}{
				"store_name": integration.Name,
			}
		} else {
			return fmt.Errorf("invalid integration config format and Name is empty")
		}
	}

	var storeDomain string
	if url, ok := config["store_url"].(string); ok && url != "" {
		// Clean up URL to get domain
		url = strings.TrimPrefix(url, "https://")
		url = strings.TrimPrefix(url, "http://")
		url = strings.TrimSuffix(url, "/")
		storeDomain = url
		fmt.Printf("[SyncOrders] Extracted storeDomain from store_url: %s\n", storeDomain)
	} else {
		// Try store_name
		if storeName, ok := config["store_name"].(string); ok && storeName != "" {
			storeDomain = storeName
			fmt.Printf("[SyncOrders] Using store_name as storeDomain: %s\n", storeDomain)
		} else {
			// Try Name again if not in config
			if integration.Name != "" {
				storeDomain = integration.Name
				fmt.Printf("[SyncOrders] Using integration Name as storeDomain: %s\n", storeDomain)
			} else {
				fmt.Printf("[SyncOrders] Config content: %+v\n", config)
				return fmt.Errorf("store_url and store_name not found in integration config")
			}
		}
	}

	// Decrypt access token
	accessToken, err := uc.coreIntegration.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		// Try fallback key if needed, or log specific error
		fmt.Printf("[SyncOrders] Failed to decrypt access_token: %v. Checking if it's plaintext in config?\n", err)
		// Often migration puts token in config? No, usually credentials.
		return fmt.Errorf("failed to decrypt access_token: %w", err)
	}

	// Calculate date filter: last 15 days
	fifteenDaysAgo := time.Now().AddDate(0, 0, -15)
	createdAtMin := fifteenDaysAgo.Format(time.RFC3339)

	// Prepare params for Shopify API
	params := map[string]string{
		"status":         "any",
		"limit":          "250",
		"created_at_min": createdAtMin,
	}
	fmt.Printf("[SyncOrders] Starting sync for integration %s. Params: %+v\n", integrationID, params)

	// Start background goroutine for syncing
	go func() {
		ctx := context.Background() // New context for background work
		totalOrders := 0

		// Fetch first page
		fmt.Println("[SyncOrders] Fetching first page...")
		orders, nextURL, err := uc.shopifyClient.FetchOrders(ctx, storeDomain, accessToken, params)
		if err != nil {
			fmt.Printf("[SyncOrders] Error fetching orders: %v\n", err)
			return
		}
		fmt.Printf("[SyncOrders] Fetched %d orders from first page. NextURL: %s\n", len(orders), nextURL)

		// Process first page
		for _, orderData := range orders {
			// Log order ID for debug
			if id, ok := orderData["id"]; ok {
				fmt.Printf("[SyncOrders] Processing order ID: %v\n", id)
			}
			if err := uc.processOrder(ctx, integration, orderData); err != nil {
				fmt.Printf("[SyncOrders] Error processing order: %v\n", err)
				continue
			}
			totalOrders++
		}

		// Handle pagination
		for nextURL != "" {
			// Shopify rate limit: 2 requests/second
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("[SyncOrders] Fetching next page: %s\n", nextURL)

			orders, nextURL, err = uc.shopifyClient.FetchOrders(ctx, storeDomain, accessToken, params)
			if err != nil {
				fmt.Printf("[SyncOrders] Error fetching paginated orders: %v\n", err)
				break
			}
			fmt.Printf("[SyncOrders] Fetched %d orders from page.\n", len(orders))

			for _, orderData := range orders {
				if err := uc.processOrder(ctx, integration, orderData); err != nil {
					fmt.Printf("[SyncOrders] Error processing order: %v\n", err)
					continue
				}
				totalOrders++
			}
		}

		// Update last_sync timestamp
		if err := uc.coreIntegration.UpdateLastSync(ctx, integrationID); err != nil {
			fmt.Printf("[SyncOrders] Error updating last_sync: %v\n", err)
		}

		fmt.Printf("[SyncOrders] Sync completed: %d orders processed\n", totalOrders)
	}()

	return nil
}

func (uc *SyncOrdersUseCase) processOrder(ctx context.Context, integration *core.Integration, orderData map[string]interface{}) error {
	// Map Shopify order to UnifiedOrder
	unifiedOrder, err := mapShopifyOrderToUnified(integration, orderData)
	if err != nil {
		return fmt.Errorf("failed to map order: %w", err)
	}

	// Publish to queue for Orders module to consume
	if err := uc.orderPublisher.Publish(ctx, unifiedOrder); err != nil {
		return fmt.Errorf("failed to publish order: %w", err)
	}

	return nil
}

func mapShopifyOrderToUnified(integration *core.Integration, orderData map[string]interface{}) (*domain.UnifiedOrder, error) {
	// Extract customer info
	customer := domain.UnifiedCustomer{}
	if customerData, ok := orderData["customer"].(map[string]interface{}); ok {
		if firstName, ok := customerData["first_name"].(string); ok {
			customer.Name = firstName
			if lastName, ok := customerData["last_name"].(string); ok {
				customer.Name += " " + lastName
			}
		}
		if email, ok := customerData["email"].(string); ok {
			customer.Email = email
		}
		if phone, ok := customerData["phone"].(string); ok {
			customer.Phone = phone
		}
	}

	// Extract shipping address
	shippingAddr := domain.UnifiedAddress{}
	if addrData, ok := orderData["shipping_address"].(map[string]interface{}); ok {
		if addr1, ok := addrData["address1"].(string); ok {
			shippingAddr.Street = addr1
		}
		if addr2, ok := addrData["address2"].(string); ok {
			shippingAddr.Address2 = addr2
		}
		if city, ok := addrData["city"].(string); ok {
			shippingAddr.City = city
		}
		if province, ok := addrData["province"].(string); ok {
			shippingAddr.State = province
		}
		if country, ok := addrData["country"].(string); ok {
			shippingAddr.Country = country
		}
		if zip, ok := addrData["zip"].(string); ok {
			shippingAddr.PostalCode = zip
		}
		// Extract Coordinates
		var lat, lng float64
		hasCoords := false
		if val, ok := addrData["latitude"].(float64); ok {
			lat = val
			hasCoords = true
		}
		if val, ok := addrData["longitude"].(float64); ok {
			lng = val
			hasCoords = true // Should check if both exist or at least one
		}
		if hasCoords {
			shippingAddr.Coordinates = &struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			}{
				Lat: lat,
				Lng: lng,
			}
		}
	}

	// Extract line items
	items := []domain.UnifiedOrderItem{}
	if lineItems, ok := orderData["line_items"].([]interface{}); ok {
		for _, item := range lineItems {
			if itemData, ok := item.(map[string]interface{}); ok {
				unifiedItem := domain.UnifiedOrderItem{}
				if id, ok := itemData["id"].(float64); ok {
					unifiedItem.ExternalID = fmt.Sprintf("%.0f", id)
				}
				if name, ok := itemData["name"].(string); ok {
					unifiedItem.Name = name
				}
				if sku, ok := itemData["sku"].(string); ok {
					unifiedItem.SKU = sku
				}
				if qty, ok := itemData["quantity"].(float64); ok {
					unifiedItem.Quantity = int(qty)
				}
				if price, ok := itemData["price"].(string); ok {
					fmt.Sscanf(price, "%f", &unifiedItem.UnitPrice)
				}
				items = append(items, unifiedItem)
			}
		}
	}

	// Parse created_at
	var occurredAt time.Time
	if createdAtStr, ok := orderData["created_at"].(string); ok {
		occurredAt, _ = time.Parse(time.RFC3339, createdAtStr)
	}

	// Extract total amount
	var totalAmount float64
	if totalPriceStr, ok := orderData["total_price"].(string); ok {
		fmt.Sscanf(totalPriceStr, "%f", &totalAmount)
	}

	// Determine BusinessID (fallback to 1 if nil/0, typical for SuperAdmin creation in dev)
	var finalBusinessID *uint
	if integration.BusinessID != nil && *integration.BusinessID > 0 {
		finalBusinessID = integration.BusinessID
	} else {
		fmt.Printf("[SyncOrders] Warning: Integration %d has nil/0 BusinessID. Using fallback BusinessID=1.\n", integration.ID)
		one := uint(1)
		finalBusinessID = &one
	}

	// Build unified order
	unifiedOrder := &domain.UnifiedOrder{
		BusinessID:      finalBusinessID,
		IntegrationID:   integration.ID,
		IntegrationType: "shopify",
		Platform:        "shopify",
		ExternalID:      fmt.Sprintf("%.0f", orderData["id"].(float64)),
		OrderNumber:     orderData["name"].(string),
		TotalAmount:     totalAmount,
		Currency:        orderData["currency"].(string),
		Customer:        customer,
		ShippingAddress: shippingAddr,
		Status:          "pending", // Default status
		OriginalStatus:  orderData["financial_status"].(string),
		Items:           items,
		Metadata:        orderData,
		OccurredAt:      occurredAt,
		ImportedAt:      time.Now(),
		OrderStatusURL:  orderData["order_status_url"].(string),
	}

	return unifiedOrder, nil
}
