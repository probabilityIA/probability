package mapper

import (
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/queue/request"
)

// MapDomainToSerializable mapea desde el dominio (sin etiquetas) a estructura serializable (con etiquetas)
func MapDomainToSerializable(order *domain.ProbabilityOrderDTO) *request.SerializableProbabilityOrderDTO {
	// Mapear items
	orderItems := make([]request.SerializableOrderItemDTO, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = request.SerializableOrderItemDTO{
			ProductID:    item.ProductID,
			ProductSKU:   item.ProductSKU,
			ProductName:  item.ProductName,
			ProductTitle: item.ProductTitle,
			VariantID:    item.VariantID,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			TotalPrice:   item.TotalPrice,
			Currency:     item.Currency,
			Discount:     item.Discount,
			Tax:          item.Tax,
			TaxRate:      item.TaxRate,
			ImageURL:     item.ImageURL,
			ProductURL:   item.ProductURL,
			Weight:       item.Weight,
			Metadata:     item.Metadata,
		}
	}

	// Mapear addresses
	addresses := make([]request.SerializableAddressDTO, len(order.Addresses))
	for i, addr := range order.Addresses {
		addresses[i] = request.SerializableAddressDTO{
			Type:         addr.Type,
			FirstName:    addr.FirstName,
			LastName:     addr.LastName,
			Company:      addr.Company,
			Phone:        addr.Phone,
			Street:       addr.Street,
			Street2:      addr.Street2,
			City:         addr.City,
			State:        addr.State,
			Country:      addr.Country,
			PostalCode:   addr.PostalCode,
			Latitude:     addr.Latitude,
			Longitude:    addr.Longitude,
			Instructions: addr.Instructions,
			Metadata:     addr.Metadata,
		}
	}

	// Mapear payments
	payments := make([]request.SerializablePaymentDTO, len(order.Payments))
	for i, pay := range order.Payments {
		paidAt := formatTime(pay.PaidAt)
		processedAt := formatTime(pay.ProcessedAt)
		refundedAt := formatTime(pay.RefundedAt)
		payments[i] = request.SerializablePaymentDTO{
			PaymentMethodID:  pay.PaymentMethodID,
			Amount:           pay.Amount,
			Currency:         pay.Currency,
			ExchangeRate:     pay.ExchangeRate,
			Status:           pay.Status,
			PaidAt:           paidAt,
			ProcessedAt:      processedAt,
			TransactionID:    pay.TransactionID,
			PaymentReference: pay.PaymentReference,
			Gateway:          pay.Gateway,
			RefundAmount:     pay.RefundAmount,
			RefundedAt:       refundedAt,
			FailureReason:    pay.FailureReason,
			Metadata:         pay.Metadata,
		}
	}

	// Mapear shipments
	shipments := make([]request.SerializableShipmentDTO, len(order.Shipments))
	for i, ship := range order.Shipments {
		shippedAt := formatTime(ship.ShippedAt)
		deliveredAt := formatTime(ship.DeliveredAt)
		estimatedDelivery := formatTime(ship.EstimatedDelivery)
		shipments[i] = request.SerializableShipmentDTO{
			TrackingNumber:    ship.TrackingNumber,
			TrackingURL:       ship.TrackingURL,
			Carrier:           ship.Carrier,
			CarrierCode:       ship.CarrierCode,
			GuideID:           ship.GuideID,
			GuideURL:          ship.GuideURL,
			Status:            ship.Status,
			ShippedAt:         shippedAt,
			DeliveredAt:       deliveredAt,
			ShippingAddressID: ship.ShippingAddressID,
			ShippingCost:      ship.ShippingCost,
			InsuranceCost:     ship.InsuranceCost,
			TotalCost:         ship.TotalCost,
			Weight:            ship.Weight,
			Height:            ship.Height,
			Width:             ship.Width,
			Length:            ship.Length,
			WarehouseID:       ship.WarehouseID,
			WarehouseName:     ship.WarehouseName,
			DriverID:          ship.DriverID,
			DriverName:        ship.DriverName,
			IsLastMile:        ship.IsLastMile,
			EstimatedDelivery: estimatedDelivery,
			DeliveryNotes:     ship.DeliveryNotes,
			Metadata:          ship.Metadata,
		}
	}

	// Mapear channel metadata
	var channelMetadata *request.SerializableChannelMetadataDTO
	if order.ChannelMetadata != nil {
		processedAt := formatTime(order.ChannelMetadata.ProcessedAt)
		lastSyncedAt := formatTime(order.ChannelMetadata.LastSyncedAt)
		channelMetadata = &request.SerializableChannelMetadataDTO{
			ChannelSource: order.ChannelMetadata.ChannelSource,
			RawData:       json.RawMessage(order.ChannelMetadata.RawData), // Convertir []byte a json.RawMessage
			Version:       order.ChannelMetadata.Version,
			ReceivedAt:    order.ChannelMetadata.ReceivedAt.Format(time.RFC3339),
			ProcessedAt:   processedAt,
			IsLatest:      order.ChannelMetadata.IsLatest,
			LastSyncedAt:  lastSyncedAt,
			SyncStatus:    order.ChannelMetadata.SyncStatus,
		}
	}

	return &request.SerializableProbabilityOrderDTO{
		BusinessID:         order.BusinessID,
		IntegrationID:      order.IntegrationID,
		IntegrationType:    order.IntegrationType,
		Platform:           order.Platform,
		ExternalID:         order.ExternalID,
		OrderNumber:        order.OrderNumber,
		InternalNumber:     order.InternalNumber,
		Subtotal:           order.Subtotal,
		Tax:                order.Tax,
		Discount:           order.Discount,
		ShippingCost:       order.ShippingCost,
		TotalAmount:        order.TotalAmount,
		Currency:           order.Currency,
		CodTotal:           order.CodTotal,
		CustomerID:         order.CustomerID,
		CustomerName:       order.CustomerName,
		CustomerEmail:      order.CustomerEmail,
		CustomerPhone:      order.CustomerPhone,
		CustomerDNI:        order.CustomerDNI,
		OrderTypeID:        order.OrderTypeID,
		OrderTypeName:      order.OrderTypeName,
		Status:             order.Status,
		OriginalStatus:     order.OriginalStatus,
		Notes:              order.Notes,
		Coupon:             order.Coupon,
		Approved:           order.Approved,
		UserID:             order.UserID,
		UserName:           order.UserName,
		Invoiceable:        order.Invoiceable,
		InvoiceURL:         order.InvoiceURL,
		InvoiceID:          order.InvoiceID,
		InvoiceProvider:    order.InvoiceProvider,
		OrderStatusURL:     order.OrderStatusURL,
		OccurredAt:         order.OccurredAt.Format(time.RFC3339),
		ImportedAt:         order.ImportedAt.Format(time.RFC3339),
		Items:              order.Items,
		Metadata:           order.Metadata,
		FinancialDetails:   order.FinancialDetails,
		ShippingDetails:    order.ShippingDetails,
		PaymentDetails:     order.PaymentDetails,
		FulfillmentDetails: order.FulfillmentDetails,
		OrderItems:         orderItems,
		Addresses:          addresses,
		Payments:           payments,
		Shipments:          shipments,
		ChannelMetadata:    channelMetadata,
	}
}

func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format(time.RFC3339)
	return &formatted
}
