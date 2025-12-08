package usecaseordermapping

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

// MapAndSaveOrder recibe una orden en formato canónico y la guarda en todas las tablas relacionadas
// Este es el punto de entrada principal para todas las integraciones después de mapear sus datos
func (uc *UseCaseOrderMapping) MapAndSaveOrder(ctx context.Context, dto *domain.CanonicalOrderDTO) (*domain.OrderResponse, error) {
	// 0. Validar datos obligatorios de integración
	if dto.IntegrationID == 0 {
		return nil, errors.New("integration_id is required")
	}
	if dto.BusinessID == nil || *dto.BusinessID == 0 {
		return nil, errors.New("business_id is required")
	}

	// 1. Validar que no exista una orden con el mismo external_id para la misma integración
	exists, err := uc.repo.OrderExists(ctx, dto.ExternalID, dto.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("error checking if order exists: %w", err)
	}
	if exists {
		return nil, domain.ErrOrderAlreadyExists
	}

	// 1.5. Validar/Crear Cliente
	client, err := uc.GetOrCreateCustomer(ctx, *dto.BusinessID, dto)
	if err != nil {
		return nil, fmt.Errorf("error processing customer: %w", err)
	}
	var clientID *uint
	if client != nil {
		clientID = &client.ID
	}

	// 2. Crear la entidad de dominio Order
	order := &domain.Order{
		// Identificadores de integración
		BusinessID:      dto.BusinessID,
		IntegrationID:   dto.IntegrationID,
		IntegrationType: dto.IntegrationType,

		// Identificadores de la orden
		Platform:       dto.Platform,
		ExternalID:     dto.ExternalID,
		OrderNumber:    dto.OrderNumber,
		InternalNumber: dto.InternalNumber,

		// Información financiera
		Subtotal:     dto.Subtotal,
		Tax:          dto.Tax,
		Discount:     dto.Discount,
		ShippingCost: dto.ShippingCost,
		TotalAmount:  dto.TotalAmount,
		Currency:     dto.Currency,
		CodTotal:     dto.CodTotal,

		// Información del cliente
		CustomerID:    clientID, // Usar el ID del cliente validado/creado
		CustomerName:  dto.CustomerName,
		CustomerEmail: dto.CustomerEmail,
		CustomerPhone: dto.CustomerPhone,
		CustomerDNI:   dto.CustomerDNI,

		// Tipo y estado
		OrderTypeID:    dto.OrderTypeID,
		OrderTypeName:  dto.OrderTypeName,
		Status:         dto.Status,
		OriginalStatus: dto.OriginalStatus,

		// Información adicional
		Notes:    dto.Notes,
		Coupon:   dto.Coupon,
		Approved: dto.Approved,
		UserID:   dto.UserID,
		UserName: dto.UserName,

		// Facturación
		Invoiceable:     dto.Invoiceable,
		InvoiceURL:      dto.InvoiceURL,
		InvoiceID:       dto.InvoiceID,
		InvoiceProvider: dto.InvoiceProvider,

		// Datos estructurados (JSONB)
		Items:              dto.Items,
		Metadata:           dto.Metadata,
		FinancialDetails:   dto.FinancialDetails,
		ShippingDetails:    dto.ShippingDetails,
		PaymentDetails:     dto.PaymentDetails,
		FulfillmentDetails: dto.FulfillmentDetails,

		// Timestamps
		OccurredAt: dto.OccurredAt,
		ImportedAt: dto.ImportedAt,
	}

	// 2.1. Asignar PaymentMethodID desde el primer pago
	order.PaymentMethodID = 1 // Valor por defecto
	if len(dto.Payments) > 0 && dto.Payments[0].PaymentMethodID > 0 {
		order.PaymentMethodID = dto.Payments[0].PaymentMethodID
		if dto.Payments[0].Status == "completed" && dto.Payments[0].PaidAt != nil {
			order.IsPaid = true
			order.PaidAt = dto.Payments[0].PaidAt
		}
	}

	// 3. Guardar la orden principal
	if err := uc.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	// 4. Guardar OrderItems
	if len(dto.OrderItems) > 0 {
		orderItems := make([]*domain.OrderItem, len(dto.OrderItems))
		for i, itemDTO := range dto.OrderItems {
			// Validar/Crear Producto
			_, err := uc.GetOrCreateProduct(ctx, *dto.BusinessID, itemDTO)
			if err != nil {
				return nil, fmt.Errorf("error processing product for item %s: %w", itemDTO.ProductSKU, err)
			}

			orderItems[i] = &domain.OrderItem{
				OrderID:          order.ID,
				ProductID:        itemDTO.ProductID,
				ProductSKU:       itemDTO.ProductSKU,
				ProductName:      itemDTO.ProductName,
				ProductTitle:     itemDTO.ProductTitle,
				VariantID:        itemDTO.VariantID,
				Quantity:         itemDTO.Quantity,
				UnitPrice:        itemDTO.UnitPrice,
				TotalPrice:       itemDTO.TotalPrice,
				Currency:         itemDTO.Currency,
				Discount:         itemDTO.Discount,
				Tax:              itemDTO.Tax,
				TaxRate:          itemDTO.TaxRate,
				ImageURL:         itemDTO.ImageURL,
				ProductURL:       itemDTO.ProductURL,
				Weight:           itemDTO.Weight,
				RequiresShipping: true,
				IsGiftCard:       false,
				Metadata:         itemDTO.Metadata,
			}
		}
		if err := uc.repo.CreateOrderItems(ctx, orderItems); err != nil {
			return nil, fmt.Errorf("error creating order items: %w", err)
		}
	}

	// 5. Guardar Addresses
	if len(dto.Addresses) > 0 {
		addresses := make([]*domain.Address, len(dto.Addresses))
		for i, addrDTO := range dto.Addresses {
			addresses[i] = &domain.Address{
				Type:         addrDTO.Type,
				OrderID:      order.ID,
				FirstName:    addrDTO.FirstName,
				LastName:     addrDTO.LastName,
				Company:      addrDTO.Company,
				Phone:        addrDTO.Phone,
				Street:       addrDTO.Street,
				Street2:      addrDTO.Street2,
				City:         addrDTO.City,
				State:        addrDTO.State,
				Country:      addrDTO.Country,
				PostalCode:   addrDTO.PostalCode,
				Latitude:     addrDTO.Latitude,
				Longitude:    addrDTO.Longitude,
				Instructions: addrDTO.Instructions,
				IsDefault:    false,
				Metadata:     addrDTO.Metadata,
			}
		}
		if err := uc.repo.CreateAddresses(ctx, addresses); err != nil {
			return nil, fmt.Errorf("error creating addresses: %w", err)
		}
	}

	// 6. Guardar Payments
	if len(dto.Payments) > 0 {
		payments := make([]*domain.Payment, len(dto.Payments))
		for i, payDTO := range dto.Payments {
			payments[i] = &domain.Payment{
				OrderID:          order.ID,
				PaymentMethodID:  payDTO.PaymentMethodID,
				Amount:           payDTO.Amount,
				Currency:         payDTO.Currency,
				ExchangeRate:     payDTO.ExchangeRate,
				Status:           payDTO.Status,
				PaidAt:           payDTO.PaidAt,
				ProcessedAt:      payDTO.ProcessedAt,
				TransactionID:    payDTO.TransactionID,
				PaymentReference: payDTO.PaymentReference,
				Gateway:          payDTO.Gateway,
				RefundAmount:     payDTO.RefundAmount,
				RefundedAt:       payDTO.RefundedAt,
				FailureReason:    payDTO.FailureReason,
				Metadata:         payDTO.Metadata,
			}
		}
		if err := uc.repo.CreatePayments(ctx, payments); err != nil {
			return nil, fmt.Errorf("error creating payments: %w", err)
		}
	}

	// 7. Guardar Shipments
	if len(dto.Shipments) > 0 {
		shipments := make([]*domain.Shipment, len(dto.Shipments))
		for i, shipDTO := range dto.Shipments {
			shipments[i] = &domain.Shipment{
				OrderID:           order.ID,
				TrackingNumber:    shipDTO.TrackingNumber,
				TrackingURL:       shipDTO.TrackingURL,
				Carrier:           shipDTO.Carrier,
				CarrierCode:       shipDTO.CarrierCode,
				GuideID:           shipDTO.GuideID,
				GuideURL:          shipDTO.GuideURL,
				Status:            shipDTO.Status,
				ShippedAt:         shipDTO.ShippedAt,
				DeliveredAt:       shipDTO.DeliveredAt,
				ShippingAddressID: shipDTO.ShippingAddressID,
				ShippingCost:      shipDTO.ShippingCost,
				InsuranceCost:     shipDTO.InsuranceCost,
				TotalCost:         shipDTO.TotalCost,
				Weight:            shipDTO.Weight,
				Height:            shipDTO.Height,
				Width:             shipDTO.Width,
				Length:            shipDTO.Length,
				WarehouseID:       shipDTO.WarehouseID,
				WarehouseName:     shipDTO.WarehouseName,
				DriverID:          shipDTO.DriverID,
				DriverName:        shipDTO.DriverName,
				IsLastMile:        shipDTO.IsLastMile,
				EstimatedDelivery: shipDTO.EstimatedDelivery,
				DeliveryNotes:     shipDTO.DeliveryNotes,
				Metadata:          shipDTO.Metadata,
			}
		}
		if err := uc.repo.CreateShipments(ctx, shipments); err != nil {
			return nil, fmt.Errorf("error creating shipments: %w", err)
		}
	}

	// 8. Guardar ChannelMetadata (datos crudos)
	if dto.ChannelMetadata != nil {
		metadata := &domain.OrderChannelMetadata{
			OrderID:       order.ID,
			ChannelSource: dto.ChannelMetadata.ChannelSource,
			IntegrationID: dto.IntegrationID,
			RawData:       dto.ChannelMetadata.RawData,
			Version:       dto.ChannelMetadata.Version,
			ReceivedAt:    dto.ChannelMetadata.ReceivedAt,
			ProcessedAt:   dto.ChannelMetadata.ProcessedAt,
			IsLatest:      dto.ChannelMetadata.IsLatest,
			LastSyncedAt:  dto.ChannelMetadata.LastSyncedAt,
			SyncStatus:    dto.ChannelMetadata.SyncStatus,
		}
		if metadata.ReceivedAt.IsZero() {
			metadata.ReceivedAt = time.Now()
		}
		if metadata.SyncStatus == "" {
			metadata.SyncStatus = "pending"
		}
		if err := uc.repo.CreateChannelMetadata(ctx, metadata); err != nil {
			return nil, fmt.Errorf("error creating channel metadata: %w", err)
		}
	}

	// 9. Publicar evento de orden creada
	if uc.eventPublisher != nil {
		eventData := domain.OrderEventData{
			OrderNumber:    order.OrderNumber,
			InternalNumber: order.InternalNumber,
			ExternalID:     order.ExternalID,
			CurrentStatus:  order.Status,
			CustomerEmail:  order.CustomerEmail,
			TotalAmount:    &order.TotalAmount,
			Currency:       order.Currency,
			Platform:       order.Platform,
		}
		event := domain.NewOrderEvent(domain.OrderEventTypeCreated, order.ID, eventData)
		event.BusinessID = order.BusinessID
		if order.IntegrationID > 0 {
			integrationID := order.IntegrationID
			event.IntegrationID = &integrationID
		}
		// Publicar de forma asíncrona (no bloquear si falla)
		go func() {
			if err := uc.eventPublisher.PublishOrderEvent(ctx, event); err != nil {
				uc.logger.Error(ctx).
					Err(err).
					Str("order_id", order.ID).
					Msg("Error al publicar evento de orden creada")
			}
		}()
	}

	// 10. Retornar la respuesta mapeada
	return mapOrderToResponse(order), nil
}

// mapOrderToResponse convierte un modelo Order a OrderResponse
func mapOrderToResponse(order *domain.Order) *domain.OrderResponse {
	return &domain.OrderResponse{
		ID:        order.ID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,

		// Identificadores de integración
		BusinessID:      order.BusinessID,
		IntegrationID:   order.IntegrationID,
		IntegrationType: order.IntegrationType,

		// Identificadores de la orden
		Platform:       order.Platform,
		ExternalID:     order.ExternalID,
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,

		// Información financiera
		Subtotal:     order.Subtotal,
		Tax:          order.Tax,
		Discount:     order.Discount,
		ShippingCost: order.ShippingCost,
		TotalAmount:  order.TotalAmount,
		Currency:     order.Currency,
		CodTotal:     order.CodTotal,

		// Información del cliente
		CustomerID:    order.CustomerID,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		CustomerPhone: order.CustomerPhone,
		CustomerDNI:   order.CustomerDNI,

		// Dirección de envío (desnormalizado)
		ShippingStreet:     order.ShippingStreet,
		ShippingCity:       order.ShippingCity,
		ShippingState:      order.ShippingState,
		ShippingCountry:    order.ShippingCountry,
		ShippingPostalCode: order.ShippingPostalCode,
		ShippingLat:        order.ShippingLat,
		ShippingLng:        order.ShippingLng,

		// Información de pago
		PaymentMethodID: order.PaymentMethodID,
		IsPaid:          order.IsPaid,
		PaidAt:          order.PaidAt,

		// Información de envío/logística
		TrackingNumber: order.TrackingNumber,
		TrackingLink:   order.TrackingLink,
		GuideID:        order.GuideID,
		GuideLink:      order.GuideLink,
		DeliveryDate:   order.DeliveryDate,
		DeliveredAt:    order.DeliveredAt,

		// Información de fulfillment
		WarehouseID:   order.WarehouseID,
		WarehouseName: order.WarehouseName,
		DriverID:      order.DriverID,
		DriverName:    order.DriverName,
		IsLastMile:    order.IsLastMile,

		// Dimensiones y peso
		Weight: order.Weight,
		Height: order.Height,
		Width:  order.Width,
		Length: order.Length,
		Boxes:  order.Boxes,

		// Tipo y estado
		OrderTypeID:    order.OrderTypeID,
		OrderTypeName:  order.OrderTypeName,
		Status:         order.Status,
		OriginalStatus: order.OriginalStatus,

		// Información adicional
		Notes:    order.Notes,
		Coupon:   order.Coupon,
		Approved: order.Approved,
		UserID:   order.UserID,
		UserName: order.UserName,

		// Facturación
		Invoiceable:     order.Invoiceable,
		InvoiceURL:      order.InvoiceURL,
		InvoiceID:       order.InvoiceID,
		InvoiceProvider: order.InvoiceProvider,

		// Datos estructurados
		Items:              order.Items,
		Metadata:           order.Metadata,
		FinancialDetails:   order.FinancialDetails,
		ShippingDetails:    order.ShippingDetails,
		PaymentDetails:     order.PaymentDetails,
		FulfillmentDetails: order.FulfillmentDetails,

		// Timestamps
		OccurredAt: order.OccurredAt,
		ImportedAt: order.ImportedAt,
	}
}
