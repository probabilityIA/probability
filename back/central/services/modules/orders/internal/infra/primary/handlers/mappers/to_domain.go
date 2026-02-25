package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/request"
)

// MapUpdateOrderRequestToDomain convierte HTTP update request a DTO de dominio
// ✅ Conversión: datatypes.JSON → []byte, maneja snake_case JSON tags
func MapUpdateOrderRequestToDomain(req *request.UpdateOrder) *dtos.UpdateOrderRequest {
	if req == nil {
		return nil
	}
	return &dtos.UpdateOrderRequest{
		Subtotal:            req.Subtotal,
		Tax:                 req.Tax,
		Discount:            req.Discount,
		ShippingCost:        req.ShippingCost,
		TotalAmount:         req.TotalAmount,
		Currency:            req.Currency,
		CodTotal:            req.CodTotal,
		CustomerName:        req.CustomerName,
		CustomerEmail:       req.CustomerEmail,
		CustomerPhone:       req.CustomerPhone,
		CustomerDNI:         req.CustomerDNI,
		CustomerOrderCount:  req.CustomerOrderCount,
		CustomerTotalSpent:  req.CustomerTotalSpent,
		ShippingStreet:      req.ShippingStreet,
		ShippingCity:        req.ShippingCity,
		ShippingState:       req.ShippingState,
		ShippingCountry:     req.ShippingCountry,
		ShippingPostalCode:  req.ShippingPostalCode,
		ShippingLat:         req.ShippingLat,
		ShippingLng:         req.ShippingLng,
		PaymentMethodID:     req.PaymentMethodID,
		IsPaid:              req.IsPaid,
		PaidAt:              req.PaidAt,
		TrackingNumber:      req.TrackingNumber,
		TrackingLink:        req.TrackingLink,
		GuideID:             req.GuideID,
		GuideLink:           req.GuideLink,
		DeliveryDate:        req.DeliveryDate,
		DeliveredAt:         req.DeliveredAt,
		WarehouseID:         req.WarehouseID,
		WarehouseName:       req.WarehouseName,
		DriverID:            req.DriverID,
		DriverName:          req.DriverName,
		IsLastMile:          req.IsLastMile,
		Weight:              req.Weight,
		Height:              req.Height,
		Width:               req.Width,
		Length:              req.Length,
		Boxes:               req.Boxes,
		OrderTypeID:         req.OrderTypeID,
		OrderTypeName:       req.OrderTypeName,
		Status:              req.Status,
		OriginalStatus:      req.OriginalStatus,
		StatusID:            req.StatusID,
		PaymentStatusID:     req.PaymentStatusID,
		FulfillmentStatusID: req.FulfillmentStatusID,
		Notes:               req.Notes,
		Coupon:              req.Coupon,
		Approved:            req.Approved,
		UserID:              req.UserID,
		UserName:            req.UserName,
		IsConfirmed:         req.IsConfirmed,
		ConfirmationStatus:  req.ConfirmationStatus,
		Novelty:             req.Novelty,
		Invoiceable:         req.Invoiceable,
		InvoiceURL:          req.InvoiceURL,
		InvoiceID:           req.InvoiceID,
		InvoiceProvider:     req.InvoiceProvider,
		Items:               []byte(req.Items),
		Metadata:            []byte(req.Metadata),
		FinancialDetails:    []byte(req.FinancialDetails),
		ShippingDetails:     []byte(req.ShippingDetails),
		PaymentDetails:      []byte(req.PaymentDetails),
		FulfillmentDetails:  []byte(req.FulfillmentDetails),
	}
}

// MapOrderRequestToDomain convierte HTTP request a DTO de dominio
// ✅ Conversión: datatypes.JSON → []byte
func MapOrderRequestToDomain(req *request.MapOrder) *dtos.ProbabilityOrderDTO {
	// Convertir datatypes.JSON ([]byte) a []byte para domain
	var itemsBytes []byte
	if req.Items != nil {
		itemsBytes = []byte(req.Items)
	}

	var metadataBytes []byte
	if req.Metadata != nil {
		metadataBytes = []byte(req.Metadata)
	}

	var financialDetailsBytes []byte
	if req.FinancialDetails != nil {
		financialDetailsBytes = []byte(req.FinancialDetails)
	}

	var shippingDetailsBytes []byte
	if req.ShippingDetails != nil {
		shippingDetailsBytes = []byte(req.ShippingDetails)
	}

	var paymentDetailsBytes []byte
	if req.PaymentDetails != nil {
		paymentDetailsBytes = []byte(req.PaymentDetails)
	}

	var fulfillmentDetailsBytes []byte
	if req.FulfillmentDetails != nil {
		fulfillmentDetailsBytes = []byte(req.FulfillmentDetails)
	}

	// Convertir relaciones anidadas
	var channelMetadata *dtos.ProbabilityChannelMetadataDTO
	if req.ChannelMetadata != nil {
		channelMetadata = mapChannelMetadataRequestToDomain(req.ChannelMetadata)
	}

	return &dtos.ProbabilityOrderDTO{
		BusinessID:              req.BusinessID,
		IntegrationID:           req.IntegrationID,
		IntegrationType:         req.IntegrationType,
		Platform:                req.Platform,
		ExternalID:              req.ExternalID,
		OrderNumber:             req.OrderNumber,
		InternalNumber:          req.InternalNumber,
		Subtotal:                req.Subtotal,
		Tax:                     req.Tax,
		Discount:                req.Discount,
		ShippingCost:            req.ShippingCost,
		TotalAmount:             req.TotalAmount,
		Currency:                req.Currency,
		CodTotal:                req.CodTotal,
		SubtotalPresentment:     req.SubtotalPresentment,
		TaxPresentment:          req.TaxPresentment,
		DiscountPresentment:     req.DiscountPresentment,
		ShippingCostPresentment: req.ShippingCostPresentment,
		TotalAmountPresentment:  req.TotalAmountPresentment,
		CurrencyPresentment:     req.CurrencyPresentment,
		CustomerID:              req.CustomerID,
		CustomerName:            req.CustomerName,
		CustomerEmail:           req.CustomerEmail,
		CustomerPhone:           req.CustomerPhone,
		CustomerDNI:             req.CustomerDNI,
		CustomerOrderCount:      req.CustomerOrderCount,
		CustomerTotalSpent:      req.CustomerTotalSpent,
		OrderTypeID:             req.OrderTypeID,
		OrderTypeName:           req.OrderTypeName,
		Status:                  req.Status,
		OriginalStatus:          req.OriginalStatus,
		StatusID:                req.StatusID,
		PaymentStatusID:         req.PaymentStatusID,
		FulfillmentStatusID:     req.FulfillmentStatusID,
		Notes:                   req.Notes,
		Coupon:                  req.Coupon,
		Approved:                req.Approved,
		UserID:                  req.UserID,
		UserName:                req.UserName,
		Invoiceable:             req.Invoiceable,
		InvoiceURL:              req.InvoiceURL,
		InvoiceID:               req.InvoiceID,
		InvoiceProvider:         req.InvoiceProvider,
		OrderStatusURL:          req.OrderStatusURL,
		OccurredAt:              req.OccurredAt,
		ImportedAt:              req.ImportedAt,
		Items:                   itemsBytes,
		Metadata:                metadataBytes,
		FinancialDetails:        financialDetailsBytes,
		ShippingDetails:         shippingDetailsBytes,
		PaymentDetails:          paymentDetailsBytes,
		FulfillmentDetails:      fulfillmentDetailsBytes,
		OrderItems:              mapOrderItemsRequestToDomain(req.OrderItems),
		Addresses:               mapAddressesRequestToDomain(req.Addresses),
		Payments:                mapPaymentsRequestToDomain(req.Payments),
		Shipments:               mapShipmentsRequestToDomain(req.Shipments),
		ChannelMetadata:         channelMetadata,
	}
}

// mapOrderItemsRequestToDomain convierte items HTTP a domain
func mapOrderItemsRequestToDomain(items []request.MapOrderItem) []dtos.ProbabilityOrderItemDTO {
	result := make([]dtos.ProbabilityOrderItemDTO, len(items))
	for i, item := range items {
		var metadataBytes []byte
		if item.Metadata != nil {
			metadataBytes = []byte(item.Metadata)
		}

		result[i] = dtos.ProbabilityOrderItemDTO{
			ProductID:             item.ProductID,
			ProductSKU:            item.ProductSKU,
			ProductName:           item.ProductName,
			ProductTitle:          item.ProductTitle,
			VariantID:             item.VariantID,
			Quantity:              item.Quantity,
			UnitPrice:             item.UnitPrice,
			TotalPrice:            item.TotalPrice,
			Currency:              item.Currency,
			Discount:              item.Discount,
			Tax:                   item.Tax,
			TaxRate:               item.TaxRate,
			UnitPricePresentment:  item.UnitPricePresentment,
			TotalPricePresentment: item.TotalPricePresentment,
			DiscountPresentment:   item.DiscountPresentment,
			TaxPresentment:        item.TaxPresentment,
			ImageURL:              item.ImageURL,
			ProductURL:            item.ProductURL,
			Weight:                item.Weight,
			Metadata:              metadataBytes,
		}
	}
	return result
}

// mapAddressesRequestToDomain convierte direcciones HTTP a domain
func mapAddressesRequestToDomain(addresses []request.MapAddress) []dtos.ProbabilityAddressDTO {
	result := make([]dtos.ProbabilityAddressDTO, len(addresses))
	for i, addr := range addresses {
		var metadataBytes []byte
		if addr.Metadata != nil {
			metadataBytes = []byte(addr.Metadata)
		}

		result[i] = dtos.ProbabilityAddressDTO{
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
			Metadata:     metadataBytes,
		}
	}
	return result
}

// mapPaymentsRequestToDomain convierte pagos HTTP a domain
func mapPaymentsRequestToDomain(payments []request.MapPayment) []dtos.ProbabilityPaymentDTO {
	result := make([]dtos.ProbabilityPaymentDTO, len(payments))
	for i, payment := range payments {
		var metadataBytes []byte
		if payment.Metadata != nil {
			metadataBytes = []byte(payment.Metadata)
		}

		result[i] = dtos.ProbabilityPaymentDTO{
			PaymentMethodID:  payment.PaymentMethodID,
			Amount:           payment.Amount,
			Currency:         payment.Currency,
			ExchangeRate:     payment.ExchangeRate,
			Status:           payment.Status,
			PaidAt:           payment.PaidAt,
			ProcessedAt:      payment.ProcessedAt,
			TransactionID:    payment.TransactionID,
			PaymentReference: payment.PaymentReference,
			Gateway:          payment.Gateway,
			RefundAmount:     payment.RefundAmount,
			RefundedAt:       payment.RefundedAt,
			FailureReason:    payment.FailureReason,
			Metadata:         metadataBytes,
		}
	}
	return result
}

// mapShipmentsRequestToDomain convierte envíos HTTP a domain
func mapShipmentsRequestToDomain(shipments []request.MapShipment) []dtos.ProbabilityShipmentDTO {
	result := make([]dtos.ProbabilityShipmentDTO, len(shipments))
	for i, shipment := range shipments {
		var metadataBytes []byte
		if shipment.Metadata != nil {
			metadataBytes = []byte(shipment.Metadata)
		}

		result[i] = dtos.ProbabilityShipmentDTO{
			TrackingNumber:    shipment.TrackingNumber,
			TrackingURL:       shipment.TrackingURL,
			Carrier:           shipment.Carrier,
			CarrierCode:       shipment.CarrierCode,
			GuideID:           shipment.GuideID,
			GuideURL:          shipment.GuideURL,
			Status:            shipment.Status,
			ShippedAt:         shipment.ShippedAt,
			DeliveredAt:       shipment.DeliveredAt,
			ShippingAddressID: shipment.ShippingAddressID,
			ShippingCost:      shipment.ShippingCost,
			InsuranceCost:     shipment.InsuranceCost,
			TotalCost:         shipment.TotalCost,
			Weight:            shipment.Weight,
			Height:            shipment.Height,
			Width:             shipment.Width,
			Length:            shipment.Length,
			WarehouseID:       shipment.WarehouseID,
			WarehouseName:     shipment.WarehouseName,
			DriverID:          shipment.DriverID,
			DriverName:        shipment.DriverName,
			IsLastMile:        shipment.IsLastMile,
			EstimatedDelivery: shipment.EstimatedDelivery,
			DeliveryNotes:     shipment.DeliveryNotes,
			Metadata:          metadataBytes,
		}
	}
	return result
}

// mapChannelMetadataRequestToDomain convierte metadata del canal HTTP a domain
func mapChannelMetadataRequestToDomain(meta *request.MapChannelMetadata) *dtos.ProbabilityChannelMetadataDTO {
	var rawDataBytes []byte
	if meta.RawData != nil {
		rawDataBytes = []byte(meta.RawData)
	}

	return &dtos.ProbabilityChannelMetadataDTO{
		ChannelSource: meta.ChannelSource,
		RawData:       rawDataBytes,
		Version:       meta.Version,
		ReceivedAt:    meta.ReceivedAt,
		ProcessedAt:   meta.ProcessedAt,
		IsLatest:      meta.IsLatest,
		LastSyncedAt:  meta.LastSyncedAt,
		SyncStatus:    meta.SyncStatus,
	}
}
