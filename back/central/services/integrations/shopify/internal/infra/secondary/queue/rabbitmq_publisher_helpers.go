package queue

import (
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

func mapUnifiedToCanonical(u *domain.UnifiedOrder) map[string]interface{} {
	// We use a map to simulate the CanonicalOrderDTO structure since we can't import it directly
	// (cyclic dependency if we import orders module, or code duplication).
	// Ideally, CanonicalOrderDTO should be in a shared lib.
	// For now, we manually flatten the structure to match what Consumer expects.

	return map[string]interface{}{
		"business_id":      u.BusinessID,
		"integration_id":   u.IntegrationID,
		"integration_type": u.IntegrationType,
		"platform":         u.Platform,
		"external_id":      u.ExternalID,
		"order_number":     u.OrderNumber,
		"total_amount":     u.TotalAmount,
		"currency":         u.Currency,
		"status":           u.Status,
		"original_status":  u.OriginalStatus,
		"customer_name":    u.Customer.Name,
		"customer_email":   u.Customer.Email,
		"customer_phone":   u.Customer.Phone,

		// Flattened Address Fields (Consumer expects these at root)
		"shipping_street":      u.ShippingAddress.Street,
		"shipping_city":        u.ShippingAddress.City,
		"shipping_state":       u.ShippingAddress.State,
		"shipping_country":     u.ShippingAddress.Country,
		"shipping_postal_code": u.ShippingAddress.PostalCode,

		// Items (Consumer expects 'items' jsonb or 'order_items' list)
		// We map to 'order_items' list structure
		"order_items": mapItems(u.Items),

		// Map to 'items' JSONB for frontend compatibility
		"items": mapItems(u.Items),

		"occurred_at": u.OccurredAt,
		"imported_at": u.ImportedAt,
		"metadata":    u.Metadata,
		"subtotal":    u.TotalAmount, // Fallback if subtotal is missing in Unified
	}
}

func mapItems(items []domain.UnifiedOrderItem) []map[string]interface{} {
	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"product_sku":  item.SKU,
			"product_name": item.Name,
			"quantity":     item.Quantity,
			"unit_price":   item.UnitPrice,
			"total_price":  item.UnitPrice * float64(item.Quantity),
			"currency":     "COP", // Default or pass from parent
		}
	}
	return result
}
