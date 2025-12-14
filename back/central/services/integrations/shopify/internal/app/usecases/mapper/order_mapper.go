package mapper

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

func MapShopifyOrderToProbability(s *domain.ShopifyOrder) *domain.ProbabilityOrderDTO {
	orderItems := make([]domain.ProbabilityOrderItemDTO, len(s.Items))
	for i, item := range s.Items {
		orderItems[i] = domain.ProbabilityOrderItemDTO{
			ProductSKU:  item.SKU,
			ProductName: item.Name,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.UnitPrice * float64(item.Quantity),
			Currency:    s.Currency,
		}
	}

	addresses := []domain.ProbabilityAddressDTO{}
	if s.ShippingAddress.Street != "" || s.ShippingAddress.City != "" {
		address := domain.ProbabilityAddressDTO{
			Type:       "shipping",
			Street:     s.ShippingAddress.Street,
			Street2:    s.ShippingAddress.Address2,
			City:       s.ShippingAddress.City,
			State:      s.ShippingAddress.State,
			Country:    s.ShippingAddress.Country,
			PostalCode: s.ShippingAddress.PostalCode,
		}
		if s.ShippingAddress.Coordinates != nil {
			address.Latitude = &s.ShippingAddress.Coordinates.Lat
			address.Longitude = &s.ShippingAddress.Coordinates.Lng
		}
		addresses = append(addresses, address)
	}

	itemsJSON, _ := json.Marshal(orderItems)

	var metadataJSON []byte
	if s.Metadata != nil {
		metadataJSON, _ = json.Marshal(s.Metadata)
	}

	subtotal := s.TotalAmount

	probabilityOrder := &domain.ProbabilityOrderDTO{
		BusinessID:      s.BusinessID,
		IntegrationID:   s.IntegrationID,
		IntegrationType: s.IntegrationType,
		Platform:        s.Platform,
		ExternalID:      s.ExternalID,
		OrderNumber:     s.OrderNumber,
		Subtotal:        subtotal,
		TotalAmount:     s.TotalAmount,
		Currency:        s.Currency,
		CustomerName:    s.Customer.Name,
		CustomerEmail:   s.Customer.Email,
		CustomerPhone:   s.Customer.Phone,
		Status:          s.Status,
		OriginalStatus:  s.OriginalStatus,
		OccurredAt:      s.OccurredAt,
		ImportedAt:      s.ImportedAt,
		Items:           itemsJSON,
		Metadata:        metadataJSON,
		OrderItems:      orderItems,
		Addresses:       addresses,
		OrderStatusURL:  s.OrderStatusURL,
	}

	return probabilityOrder
}
