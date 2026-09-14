package mapper

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"gorm.io/datatypes"
)

func ToOrderResponse(order *entities.ProbabilityOrder) *dtos.OrderResponse {
	if order == nil {
		return nil
	}

	return &dtos.OrderResponse{
		ID:        order.ID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,

		BusinessID:         order.BusinessID,
		IntegrationID:      order.IntegrationID,
		IntegrationType:    order.IntegrationType,
		IntegrationLogoURL: order.IntegrationLogoURL,
		IntegrationName:    order.IntegrationName,

		Platform:       order.Platform,
		ExternalID:     order.ExternalID,
		ChannelPackID:  order.ChannelPackID,
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,

		Subtotal:                    order.Subtotal,
		Tax:                         order.Tax,
		Discount:                    order.Discount,
		ShippingCost:                order.ShippingCost,
		ShippingDiscount:            order.ShippingDiscount,
		ShippingDiscountPresentment: order.ShippingDiscountPresentment,
		TotalAmount:                 order.TotalAmount,
		Currency:                    order.Currency,
		IsCod:                       order.IsCod,
		CodTotal:                    order.CodTotal,
		CodIncludesShipping:         order.CodIncludesShipping,
		CodCheckoutCarrierFee:       order.CodCheckoutCarrierFee,
		CodCutConfirmed:             order.CodCutConfirmed,
		SubtotalPresentment:         order.SubtotalPresentment,
		TaxPresentment:              order.TaxPresentment,
		DiscountPresentment:         order.DiscountPresentment,
		ShippingCostPresentment:     order.ShippingCostPresentment,
		TotalAmountPresentment:      order.TotalAmountPresentment,
		CurrencyPresentment:         order.CurrencyPresentment,

		CustomerID:        order.CustomerID,
		CustomerName:      order.CustomerName,
		CustomerFirstName: order.CustomerFirstName,
		CustomerLastName:  order.CustomerLastName,
		CustomerEmail:     order.CustomerEmail,
		CustomerPhone:     order.CustomerPhone,
		CustomerDNI:       order.CustomerDNI,

		ShippingStreet:           order.ShippingStreet,
		ShippingCity:             order.ShippingCity,
		ShippingState:            order.ShippingState,
		ShippingCountry:          order.ShippingCountry,
		ShippingPostalCode:       order.ShippingPostalCode,
		ShippingLat:              order.ShippingLat,
		ShippingLng:              order.ShippingLng,
		ShippingGeoConfidence:    order.ShippingGeoConfidence,
		ShippingAddressSource:    order.ShippingAddressSource,
		ShippingDeliveryType:     order.ShippingDeliveryType,
		ShippingOfficeCarrier:    order.ShippingOfficeCarrier,
		ShippingOfficeName:       order.ShippingOfficeName,
		ShippingOfficeDetails:    order.ShippingOfficeDetails,
		ShippingNeighborhood:     order.ShippingNeighborhood,
		ShippingComplementType:   order.ShippingComplementType,
		ShippingComplementNumber: order.ShippingComplementNumber,
		ShippingTower:            order.ShippingTower,
		ShippingBuilding:         order.ShippingBuilding,
		DestinationDaneCode:      order.DestinationDaneCode,

		PaymentMethodID: order.PaymentMethodID,
		IsPaid:          order.IsPaid,
		PaidAt:          order.PaidAt,

		TrackingNumber:      order.TrackingNumber,
		TrackingLink:        order.TrackingLink,
		GuideID:             order.GuideID,
		GuideLink:           order.GuideLink,
		DeliveryDate:        order.DeliveryDate,
		DeliveredAt:         order.DeliveredAt,
		DeliveryProbability: order.DeliveryProbability,

		WarehouseID:   order.WarehouseID,
		WarehouseName: order.WarehouseName,
		DriverID:      order.DriverID,
		DriverName:    order.DriverName,
		IsLastMile:    order.IsLastMile,

		Weight: order.Weight,
		Height: order.Height,
		Width:  order.Width,
		Length: order.Length,
		Boxes:  order.Boxes,

		OrderTypeID:         order.OrderTypeID,
		OrderTypeName:       order.OrderTypeName,
		Status:              order.Status,
		OriginalStatus:      order.OriginalStatus,
		StatusID:            order.StatusID,
		OrderStatus:         order.OrderStatus,
		PaymentStatusID:     order.PaymentStatusID,
		FulfillmentStatusID: order.FulfillmentStatusID,
		PaymentStatus:       order.PaymentStatus,
		FulfillmentStatus:   order.FulfillmentStatus,

		Notes:    order.Notes,
		Coupon:   order.Coupon,
		Approved: order.Approved,
		UserID:   order.UserID,
		UserName: order.UserName,

		IsConfirmed: order.IsConfirmed,
		Novelty:     order.Novelty,

		IsTest: order.IsTest,

		Invoiceable:     order.Invoiceable,
		InvoiceURL:      order.InvoiceURL,
		InvoiceID:       order.InvoiceID,
		InvoiceProvider: order.InvoiceProvider,
		OrderStatusURL:  order.OrderStatusURL,

		OrderItems: order.OrderItems,

		Shipment: mapShipmentToResponse(order.Shipments),

		Metadata:           order.Metadata,
		FinancialDetails:   order.FinancialDetails,
		ShippingDetails:    order.ShippingDetails,
		FreeShipping:       order.FreeShipping,
		PaymentDetails:     order.PaymentDetails,
		FulfillmentDetails: order.FulfillmentDetails,

		OccurredAt: order.OccurredAt,
		ImportedAt: order.ImportedAt,

		NegativeFactors: UnmarshalNegativeFactors(order.NegativeFactors),
		ScoreBreakdown:  json.RawMessage(order.ScoreBreakdown),
	}
}

func UnmarshalNegativeFactors(jsonData datatypes.JSON) []string {
	if len(jsonData) == 0 || string(jsonData) == "null" {
		return []string{}
	}
	var factors []string
	_ = json.Unmarshal(jsonData, &factors)
	return factors
}

func mapShipmentToResponse(shipments []entities.ProbabilityShipment) *dtos.ShipmentData {
	if len(shipments) == 0 {
		return nil
	}

	s := shipments[0]
	codCarrierFee := s.CodCarrierFee
	if codCarrierFee == nil || *codCarrierFee == 0 {
		for i := range shipments {
			if shipments[i].CodCarrierFee != nil && *shipments[i].CodCarrierFee > 0 {
				codCarrierFee = shipments[i].CodCarrierFee
				break
			}
		}
	}
	return &dtos.ShipmentData{
		ID:                  s.ID,
		Carrier:             s.Carrier,
		TrackingNumber:      s.TrackingNumber,
		GuideURL:            s.GuideURL,
		Status:              s.Status,
		CarrierStatus:       s.CarrierStatus,
		CarrierStatusDetail: s.CarrierStatusDetail,
		TotalCost:           s.TotalCost,
		CodCarrierFee:       codCarrierFee,
	}
}

func ToOrderSummary(order *entities.ProbabilityOrder) dtos.OrderSummary {
	var businessID uint
	if order.BusinessID != nil {
		businessID = *order.BusinessID
	}

	var shipment *dtos.ShipmentSummary
	if len(order.Shipments) > 0 {
		s := order.Shipments[0]
		shipment = &dtos.ShipmentSummary{
			ID:                  s.ID,
			Carrier:             s.Carrier,
			TrackingNumber:      s.TrackingNumber,
			GuideURL:            s.GuideURL,
			Status:              s.Status,
			CarrierStatus:       s.CarrierStatus,
			CarrierStatusDetail: s.CarrierStatusDetail,
			TotalCost:           s.TotalCost,
			CodCarrierFee:       s.CodCarrierFee,
		}
	}

	return dtos.OrderSummary{
		ID:                     order.ID,
		ShippingDetails:        order.ShippingDetails,
		FreeShipping:           order.FreeShipping,
		StatusSource:           order.StatusSource,
		StatusChangedBy:        order.StatusChangedBy,
		StatusChangedAt:        order.StatusChangedAt,
		CreatedAt:              order.CreatedAt,
		BusinessID:             businessID,
		IntegrationID:          order.IntegrationID,
		IntegrationType:        order.IntegrationType,
		IntegrationLogoURL:     order.IntegrationLogoURL,
		Platform:               order.Platform,
		ExternalID:             order.ExternalID,
		ChannelPackID:          order.ChannelPackID,
		OrderNumber:            order.OrderNumber,
		TotalAmount:            order.TotalAmount,
		Currency:               order.Currency,
		TotalAmountPresentment: order.TotalAmountPresentment,
		CurrencyPresentment:    order.CurrencyPresentment,
		CustomerName:           order.CustomerName,
		CustomerFirstName:      order.CustomerFirstName,
		CustomerLastName:       order.CustomerLastName,
		CustomerEmail:          order.CustomerEmail,
		CustomerPhone:          order.CustomerPhone,
		ShippingStreet:         order.ShippingStreet,
		ShippingCity:           order.ShippingCity,
		ShippingState:          order.ShippingState,
		ShippingGeoConfidence:  order.ShippingGeoConfidence,
		ShippingAddressSource:  order.ShippingAddressSource,
		ShippingDeliveryType:   order.ShippingDeliveryType,
		ShippingOfficeCarrier:  order.ShippingOfficeCarrier,
		ShippingOfficeName:     order.ShippingOfficeName,
		Weight:                 order.Weight,
		Height:                 order.Height,
		Width:                  order.Width,
		Length:                 order.Length,
		Status:                 order.Status,
		ItemsCount:             len(order.OrderItems),
		DeliveryProbability:    order.DeliveryProbability,
		NegativeFactors:        UnmarshalNegativeFactors(order.NegativeFactors),
		ScoreBreakdown:         json.RawMessage(order.ScoreBreakdown),
		OrderStatus:            order.OrderStatus,
		PaymentStatus:          order.PaymentStatus,
		FulfillmentStatus:      order.FulfillmentStatus,
		OrderStatusURL:         order.OrderStatusURL,
		GuideLink:              order.GuideLink,
		IsPaid:                 order.IsPaid,
		IsCod:                  order.IsCod,
		CodTotal:               order.CodTotal,
		IsConfirmed:            order.IsConfirmed,
		Novelty:                order.Novelty,
		IsTest:                 order.IsTest,
		InvoiceStatus:          order.InvoiceStatus,
		CodCutConfirmed:        order.CodCutConfirmed,
		Shipment:               shipment,
	}
}

func ToDomainOrderItems(items []map[string]interface{}) []entities.ProbabilityOrderItem {
	if len(items) == 0 {
		return []entities.ProbabilityOrderItem{}
	}

	result := make([]entities.ProbabilityOrderItem, len(items))
	for i, itemMap := range items {

		qty := 1.0
		if q, ok := itemMap["quantity"].(float64); ok {
			qty = q
		}

		price := 0.0
		if p, ok := itemMap["price"].(float64); ok {
			price = p
		}

		sku := ""
		if s, ok := itemMap["sku"].(string); ok {
			sku = s
		}

		name := ""
		if n, ok := itemMap["name"].(string); ok {
			name = n
		}

		variantLabel := ""
		if vl, ok := itemMap["variant_label"].(string); ok {
			variantLabel = vl
		}

		var productID *string
		if productIDVal, ok := itemMap["product_id"].(string); ok {
			productID = &productIDVal
		}

		result[i] = entities.ProbabilityOrderItem{
			ProductID:    productID,
			ProductSKU:   sku,
			ProductName:  name,
			VariantLabel: variantLabel,
			VariantID:    nil,
			Quantity:     int(qty),
			UnitPrice:    price,
			TotalPrice:   price * qty,
			Currency:     "COP",
			Discount:     0,
			Tax:          0,
		}
	}
	return result
}
