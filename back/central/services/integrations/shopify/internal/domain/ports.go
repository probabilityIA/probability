package domain

import (
	"context"
)

// OrderPublisher defines the interface for publishing orders to the system (e.g., via RabbitMQ)
type OrderPublisher interface {
	Publish(ctx context.Context, order *UnifiedOrder) error
}

// ShopifyClient defines the interface for interacting with the Shopify API
type ShopifyClient interface {
	// ValidateToken checks if the access token is valid for the store
	ValidateToken(ctx context.Context, storeName, accessToken string) (bool, map[string]interface{}, error)

	// FetchOrders retrieves orders from Shopify.
	// Returns a list of orders (as maps/structs) and the next page cursor/link if any.
	FetchOrders(ctx context.Context, storeName, accessToken string, params map[string]string) ([]map[string]interface{}, string, error)
}
