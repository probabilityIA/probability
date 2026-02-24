package mapper

import (
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/magento/internal/infra/secondary/queue/request"
)

// MapDomainToSerializable mapea desde el dominio (sin etiquetas JSON) a la estructura
// serializable (con etiquetas JSON) para publicar en RabbitMQ.
func MapDomainToSerializable(order *canonical.ProbabilityOrderDTO) *request.SerializableProbabilityOrderDTO {
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

	payments := make([]request.SerializablePaymentDTO, len(order.Payments))
	for i, p := range order.Payments {
		payments[i] = request.SerializablePaymentDTO{
			PaymentMethodID:  p.PaymentMethodID,
			Amount:           p.Amount,
			Currency:         p.Currency,
			ExchangeRate:     p.ExchangeRate,
			Status:           p.Status,
			PaidAt:           formatTimePtr(p.PaidAt),
			ProcessedAt:      formatTimePtr(p.ProcessedAt),
			TransactionID:    p.TransactionID,
			PaymentReference: p.PaymentReference,
			Gateway:          p.Gateway,
			RefundAmount:     p.RefundAmount,
			RefundedAt:       formatTimePtr(p.RefundedAt),
			FailureReason:    p.FailureReason,
			Metadata:         json.RawMessage(p.Metadata),
		}
	}

	shipments := make([]request.SerializableShipmentDTO, len(order.Shipments))
	for i, s := range order.Shipments {
		shipments[i] = request.SerializableShipmentDTO{
			TrackingNumber:    s.TrackingNumber,
			TrackingURL:       s.TrackingURL,
			Carrier:           s.Carrier,
			CarrierCode:       s.CarrierCode,
			GuideID:           s.GuideID,
			GuideURL:          s.GuideURL,
			Status:            s.Status,
			ShippedAt:         formatTimePtr(s.ShippedAt),
			DeliveredAt:       formatTimePtr(s.DeliveredAt),
			ShippingAddressID: s.ShippingAddressID,
			ShippingCost:      s.ShippingCost,
			InsuranceCost:     s.InsuranceCost,
			TotalCost:         s.TotalCost,
			Weight:            s.Weight,
			Height:            s.Height,
			Width:             s.Width,
			Length:            s.Length,
			WarehouseID:       s.WarehouseID,
			WarehouseName:     s.WarehouseName,
			DriverID:          s.DriverID,
			DriverName:        s.DriverName,
			IsLastMile:        s.IsLastMile,
			EstimatedDelivery: formatTimePtr(s.EstimatedDelivery),
			DeliveryNotes:     s.DeliveryNotes,
			Metadata:          json.RawMessage(s.Metadata),
		}
	}

	var channelMeta *request.SerializableChannelMetadataDTO
	if order.ChannelMetadata != nil {
		cm := order.ChannelMetadata
		channelMeta = &request.SerializableChannelMetadataDTO{
			ChannelSource: cm.ChannelSource,
			RawData:       json.RawMessage(cm.RawData),
			Version:       cm.Version,
			ReceivedAt:    cm.ReceivedAt.Format(time.RFC3339),
			ProcessedAt:   formatTimePtr(cm.ProcessedAt),
			IsLatest:      cm.IsLatest,
			LastSyncedAt:  formatTimePtr(cm.LastSyncedAt),
			SyncStatus:    cm.SyncStatus,
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
		Items:              json.RawMessage(order.Items),
		Metadata:           json.RawMessage(order.Metadata),
		FinancialDetails:   json.RawMessage(order.FinancialDetails),
		ShippingDetails:    json.RawMessage(order.ShippingDetails),
		PaymentDetails:     json.RawMessage(order.PaymentDetails),
		FulfillmentDetails: json.RawMessage(order.FulfillmentDetails),
		OrderItems:         orderItems,
		Addresses:          addresses,
		Payments:           payments,
		Shipments:          shipments,
		ChannelMetadata:    channelMeta,
	}
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
