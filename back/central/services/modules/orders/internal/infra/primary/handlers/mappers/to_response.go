package mappers

import (
	"gorm.io/datatypes"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/response"
)

// OrderToResponse convierte DTO de dominio a HTTP response
// ✅ Conversión: []byte → datatypes.JSON
func OrderToResponse(dto *dtos.OrderResponse) *response.Order {
	// Convertir []byte a datatypes.JSON
	var itemsJSON datatypes.JSON
	if len(dto.Items) > 0 {
		itemsJSON = datatypes.JSON(dto.Items)
	}

	var metadataJSON datatypes.JSON
	if len(dto.Metadata) > 0 {
		metadataJSON = datatypes.JSON(dto.Metadata)
	}

	var financialDetailsJSON datatypes.JSON
	if len(dto.FinancialDetails) > 0 {
		financialDetailsJSON = datatypes.JSON(dto.FinancialDetails)
	}

	var shippingDetailsJSON datatypes.JSON
	if len(dto.ShippingDetails) > 0 {
		shippingDetailsJSON = datatypes.JSON(dto.ShippingDetails)
	}

	var paymentDetailsJSON datatypes.JSON
	if len(dto.PaymentDetails) > 0 {
		paymentDetailsJSON = datatypes.JSON(dto.PaymentDetails)
	}

	var fulfillmentDetailsJSON datatypes.JSON
	if len(dto.FulfillmentDetails) > 0 {
		fulfillmentDetailsJSON = datatypes.JSON(dto.FulfillmentDetails)
	}

	// Mapear información de estados
	var orderStatus *response.OrderStatusInfo
	if dto.OrderStatus != nil {
		orderStatus = mapOrderStatusToResponse(dto.OrderStatus)
	}

	var paymentStatus *response.PaymentStatusInfo
	if dto.PaymentStatus != nil {
		paymentStatus = mapPaymentStatusToResponse(dto.PaymentStatus)
	}

	var fulfillmentStatus *response.FulfillmentStatusInfo
	if dto.FulfillmentStatus != nil {
		fulfillmentStatus = mapFulfillmentStatusToResponse(dto.FulfillmentStatus)
	}

	return &response.Order{
		ID:                      dto.ID,
		CreatedAt:               dto.CreatedAt,
		UpdatedAt:               dto.UpdatedAt,
		DeletedAt:               dto.DeletedAt,
		BusinessID:              dto.BusinessID,
		IntegrationID:           dto.IntegrationID,
		IntegrationType:         dto.IntegrationType,
		IntegrationLogoURL:      dto.IntegrationLogoURL,
		Platform:                dto.Platform,
		ExternalID:              dto.ExternalID,
		OrderNumber:             dto.OrderNumber,
		InternalNumber:          dto.InternalNumber,
		Subtotal:                dto.Subtotal,
		Tax:                     dto.Tax,
		Discount:                dto.Discount,
		ShippingCost:            dto.ShippingCost,
		TotalAmount:             dto.TotalAmount,
		Currency:                dto.Currency,
		CodTotal:                dto.CodTotal,
		SubtotalPresentment:     dto.SubtotalPresentment,
		TaxPresentment:          dto.TaxPresentment,
		DiscountPresentment:     dto.DiscountPresentment,
		ShippingCostPresentment: dto.ShippingCostPresentment,
		TotalAmountPresentment:  dto.TotalAmountPresentment,
		CurrencyPresentment:     dto.CurrencyPresentment,
		CustomerID:              dto.CustomerID,
		CustomerName:            dto.CustomerName,
		CustomerEmail:           dto.CustomerEmail,
		CustomerPhone:           dto.CustomerPhone,
		CustomerDNI:             dto.CustomerDNI,
		ShippingStreet:          dto.ShippingStreet,
		ShippingCity:            dto.ShippingCity,
		ShippingState:           dto.ShippingState,
		ShippingCountry:         dto.ShippingCountry,
		ShippingPostalCode:      dto.ShippingPostalCode,
		ShippingLat:             dto.ShippingLat,
		ShippingLng:             dto.ShippingLng,
		PaymentMethodID:         dto.PaymentMethodID,
		IsPaid:                  dto.IsPaid,
		PaidAt:                  dto.PaidAt,
		TrackingNumber:          dto.TrackingNumber,
		TrackingLink:            dto.TrackingLink,
		GuideID:                 dto.GuideID,
		GuideLink:               dto.GuideLink,
		DeliveryDate:            dto.DeliveryDate,
		DeliveredAt:             dto.DeliveredAt,
		DeliveryProbability:     dto.DeliveryProbability,
		WarehouseID:             dto.WarehouseID,
		WarehouseName:           dto.WarehouseName,
		DriverID:                dto.DriverID,
		DriverName:              dto.DriverName,
		IsLastMile:              dto.IsLastMile,
		Weight:                  dto.Weight,
		Height:                  dto.Height,
		Width:                   dto.Width,
		Length:                  dto.Length,
		Boxes:                   dto.Boxes,
		OrderTypeID:             dto.OrderTypeID,
		OrderTypeName:           dto.OrderTypeName,
		Status:                  dto.Status,
		OriginalStatus:          dto.OriginalStatus,
		StatusID:                dto.StatusID,
		OrderStatus:             orderStatus,
		PaymentStatusID:         dto.PaymentStatusID,
		FulfillmentStatusID:     dto.FulfillmentStatusID,
		PaymentStatus:           paymentStatus,
		FulfillmentStatus:       fulfillmentStatus,
		Notes:                   dto.Notes,
		Coupon:                  dto.Coupon,
		Approved:                dto.Approved,
		UserID:                  dto.UserID,
		UserName:                dto.UserName,
		IsConfirmed:             dto.IsConfirmed,
		Novelty:                 dto.Novelty,
		Invoiceable:             dto.Invoiceable,
		InvoiceURL:              dto.InvoiceURL,
		InvoiceID:               dto.InvoiceID,
		InvoiceProvider:         dto.InvoiceProvider,
		OrderStatusURL:          dto.OrderStatusURL,
		Items:                   itemsJSON,
		Metadata:                metadataJSON,
		FinancialDetails:        financialDetailsJSON,
		ShippingDetails:         shippingDetailsJSON,
		PaymentDetails:          paymentDetailsJSON,
		FulfillmentDetails:      fulfillmentDetailsJSON,
		NegativeFactors:         dto.NegativeFactors,
		OccurredAt:              dto.OccurredAt,
		ImportedAt:              dto.ImportedAt,
	}
}

// OrderSummaryToResponse convierte resumen de orden de dominio a HTTP response
func OrderSummaryToResponse(dto *dtos.OrderSummary) *response.OrderSummary {
	var orderStatus *response.OrderStatusInfo
	if dto.OrderStatus != nil {
		orderStatus = mapOrderStatusToResponse(dto.OrderStatus)
	}

	var paymentStatus *response.PaymentStatusInfo
	if dto.PaymentStatus != nil {
		paymentStatus = mapPaymentStatusToResponse(dto.PaymentStatus)
	}

	var fulfillmentStatus *response.FulfillmentStatusInfo
	if dto.FulfillmentStatus != nil {
		fulfillmentStatus = mapFulfillmentStatusToResponse(dto.FulfillmentStatus)
	}

	return &response.OrderSummary{
		ID:                     dto.ID,
		CreatedAt:              dto.CreatedAt,
		BusinessID:             dto.BusinessID,
		IntegrationID:          dto.IntegrationID,
		IntegrationType:        dto.IntegrationType,
		IntegrationLogoURL:     dto.IntegrationLogoURL,
		Platform:               dto.Platform,
		ExternalID:             dto.ExternalID,
		OrderNumber:            dto.OrderNumber,
		TotalAmount:            dto.TotalAmount,
		Currency:               dto.Currency,
		TotalAmountPresentment: dto.TotalAmountPresentment,
		CurrencyPresentment:    dto.CurrencyPresentment,
		CustomerName:           dto.CustomerName,
		CustomerEmail:          dto.CustomerEmail,
		CustomerPhone:          dto.CustomerPhone,
		ShippingStreet:         dto.ShippingStreet,
		ShippingCity:           dto.ShippingCity,
		ShippingState:          dto.ShippingState,
		Weight:                 dto.Weight,
		Height:                 dto.Height,
		Width:                  dto.Width,
		Length:                 dto.Length,
		Status:                 dto.Status,
		ItemsCount:             dto.ItemsCount,
		DeliveryProbability:    dto.DeliveryProbability,
		NegativeFactors:        dto.NegativeFactors,
		OrderStatus:            orderStatus,
		PaymentStatus:          paymentStatus,
		FulfillmentStatus:      fulfillmentStatus,
		IsConfirmed:            dto.IsConfirmed,
		Novelty:                dto.Novelty,
	}
}

// OrderRawToResponse convierte respuesta raw de dominio a HTTP response
func OrderRawToResponse(dto *dtos.OrderRawResponse) *response.OrderRaw {
	var rawDataJSON datatypes.JSON
	if len(dto.RawData) > 0 {
		rawDataJSON = datatypes.JSON(dto.RawData)
	}

	return &response.OrderRaw{
		OrderID:       dto.OrderID,
		ChannelSource: dto.ChannelSource,
		RawData:       rawDataJSON,
	}
}

// OrdersListToResponse convierte lista paginada de dominio a HTTP response
func OrdersListToResponse(dto *dtos.OrdersListResponse) *response.OrdersList {
	summaries := make([]response.OrderSummary, len(dto.Data))
	for i, summary := range dto.Data {
		summaries[i] = *OrderSummaryToResponse(&summary)
	}

	return &response.OrdersList{
		Data:       summaries,
		Total:      dto.Total,
		Page:       dto.Page,
		PageSize:   dto.PageSize,
		TotalPages: dto.TotalPages,
	}
}

// mapOrderStatusToResponse convierte OrderStatusInfo de entities a response
func mapOrderStatusToResponse(status *entities.OrderStatusInfo) *response.OrderStatusInfo {
	return &response.OrderStatusInfo{
		ID:          status.ID,
		Code:        status.Code,
		Name:        status.Name,
		Description: status.Description,
		Category:    status.Category,
		Color:       status.Color,
	}
}

// mapPaymentStatusToResponse convierte PaymentStatusInfo de entities a response
func mapPaymentStatusToResponse(status *entities.PaymentStatusInfo) *response.PaymentStatusInfo {
	return &response.PaymentStatusInfo{
		ID:          status.ID,
		Code:        status.Code,
		Name:        status.Name,
		Description: status.Description,
		Category:    status.Category,
		Color:       status.Color,
	}
}

// mapFulfillmentStatusToResponse convierte FulfillmentStatusInfo de entities a response
func mapFulfillmentStatusToResponse(status *entities.FulfillmentStatusInfo) *response.FulfillmentStatusInfo {
	return &response.FulfillmentStatusInfo{
		ID:          status.ID,
		Code:        status.Code,
		Name:        status.Name,
		Description: status.Description,
		Category:    status.Category,
		Color:       status.Color,
	}
}
