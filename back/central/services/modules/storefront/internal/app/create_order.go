package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
)

func (uc *UseCase) CreateOrder(ctx context.Context, businessID, userID uint, dto *dtos.StorefrontCreateOrderDTO) error {
	if len(dto.Items) == 0 {
		return domainerrors.ErrNoItems
	}

	// Validate quantities
	for _, item := range dto.Items {
		if item.Quantity < 1 {
			return domainerrors.ErrInvalidQuantity
		}
	}

	// Get client data for the authenticated user
	client, err := uc.repo.GetClientByUserID(ctx, businessID, userID)
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	// Get platform integration_id for this business
	integrationID, err := uc.repo.GetPlatformIntegrationID(ctx, businessID)
	if err != nil {
		return fmt.Errorf("failed to get platform integration: %w", err)
	}

	// Look up product details and build order items
	var totalAmount float64
	orderItems := make([]map[string]interface{}, 0, len(dto.Items))

	for _, item := range dto.Items {
		product, err := uc.repo.GetProductByID(ctx, businessID, item.ProductID)
		if err != nil {
			return fmt.Errorf("product %s not found: %w", item.ProductID, err)
		}

		itemTotal := product.Price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, map[string]interface{}{
			"product_id":   product.ID,
			"product_sku":  product.SKU,
			"product_name": product.Name,
			"quantity":     item.Quantity,
			"unit_price":   product.Price,
			"total_price":  itemTotal,
			"currency":     product.Currency,
			"image_url":    product.ImageURL,
		})
	}

	// Build addresses
	var addresses []map[string]interface{}
	if dto.Address != nil {
		addresses = append(addresses, map[string]interface{}{
			"type":         "shipping",
			"first_name":   dto.Address.FirstName,
			"last_name":    dto.Address.LastName,
			"phone":        dto.Address.Phone,
			"street":       dto.Address.Street,
			"street2":      dto.Address.Street2,
			"city":         dto.Address.City,
			"state":        dto.Address.State,
			"country":      dto.Address.Country,
			"postal_code":  dto.Address.PostalCode,
			"instructions": dto.Address.Instructions,
		})
	}

	// Generate unique external_id and order_number
	now := time.Now()
	externalID := fmt.Sprintf("sf-%d", now.UnixNano())
	orderNumber := fmt.Sprintf("SF-%d", now.UnixNano()%1000000)

	customerEmail := ""
	if client.Email != nil {
		customerEmail = *client.Email
	}

	// Build the canonical order DTO (same format all ecommerce integrations use)
	canonicalOrder := map[string]interface{}{
		"business_id":      businessID,
		"integration_id":   integrationID,
		"integration_type": "platform",
		"platform":         "storefront",
		"external_id":      externalID,
		"order_number":     orderNumber,
		"subtotal":         totalAmount,
		"tax":              0,
		"discount":         0,
		"shipping_cost":    0,
		"total_amount":     totalAmount,
		"currency":         "COP",
		"customer_name":    client.Name,
		"customer_email":   customerEmail,
		"customer_phone":   client.Phone,
		"status":           "pending",
		"original_status":  "pending",
		"notes":            dto.Notes,
		"user_id":          userID,
		"invoiceable":      false,
		"occurred_at":      now.Format(time.RFC3339),
		"imported_at":      now.Format(time.RFC3339),
		"order_items":      orderItems,
		"addresses":        addresses,
		"payments":         []interface{}{},
		"shipments":        []interface{}{},
	}

	if client.Dni != nil {
		canonicalOrder["customer_dni"] = *client.Dni
	}

	orderJSON, err := json.Marshal(canonicalOrder)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	if err := uc.publisher.PublishOrder(ctx, orderJSON); err != nil {
		return fmt.Errorf("failed to publish order: %w", err)
	}

	uc.logger.Info(ctx).
		Str("external_id", externalID).
		Str("order_number", orderNumber).
		Uint("business_id", businessID).
		Uint("user_id", userID).
		Msg("Storefront order published to queue")

	return nil
}
