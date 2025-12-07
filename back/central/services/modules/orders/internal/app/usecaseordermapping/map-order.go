package usecaseordermapping

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// MapAndSaveOrder recibe una orden en formato canónico y la guarda en todas las tablas relacionadas
// Este es el punto de entrada principal para todas las integraciones después de mapear sus datos
func (uc *UseCaseOrderMapping) MapAndSaveOrder(ctx context.Context, dto *domain.CanonicalOrderDTO) (*domain.OrderResponse, error) {
	// 1. Validar que no exista una orden con el mismo external_id para la misma integración
	exists, err := uc.repo.OrderExists(ctx, dto.ExternalID, dto.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("error checking if order exists: %w", err)
	}
	if exists {
		return nil, errors.New("order with this external_id already exists for this integration")
	}

	// 2. Crear el modelo de orden principal
	order := &models.Order{
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

		// Información del cliente (desnormalizado para compatibilidad)
		CustomerID:    dto.CustomerID,
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

		// Datos estructurados (JSONB) - Para compatibilidad
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

	// 2.1. Asignar PaymentMethodID desde el primer pago (si existe y es válido)
	// El modelo Order requiere PaymentMethodID (not null), así que lo tomamos del primer pago
	// Si no hay pagos o el PaymentMethodID es 0, usamos un valor por defecto
	order.PaymentMethodID = 1 // Valor por defecto (debe existir en payment_methods)
	if len(dto.Payments) > 0 && dto.Payments[0].PaymentMethodID > 0 {
		order.PaymentMethodID = dto.Payments[0].PaymentMethodID
		// También establecer IsPaid y PaidAt si el pago está completado
		if dto.Payments[0].Status == "completed" && dto.Payments[0].PaidAt != nil {
			order.IsPaid = true
			order.PaidAt = dto.Payments[0].PaidAt
		}
	}

	// 3. Guardar la orden principal primero (para obtener el ID)
	if err := uc.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	// 4. Guardar OrderItems
	if len(dto.OrderItems) > 0 {
		orderItems := make([]*models.OrderItem, len(dto.OrderItems))
		for i, itemDTO := range dto.OrderItems {
			orderItems[i] = &models.OrderItem{
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
				RequiresShipping: true, // Por defecto requiere envío
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
		addresses := make([]*models.Address, len(dto.Addresses))
		for i, addrDTO := range dto.Addresses {
			addresses[i] = &models.Address{
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
		payments := make([]*models.Payment, len(dto.Payments))
		for i, payDTO := range dto.Payments {
			payments[i] = &models.Payment{
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
		shipments := make([]*models.Shipment, len(dto.Shipments))
		for i, shipDTO := range dto.Shipments {
			shipments[i] = &models.Shipment{
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
		metadata := &models.OrderChannelMetadata{
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

	// 9. Retornar la respuesta mapeada
	return mapOrderToResponse(order), nil
}

// mapOrderToResponse convierte un modelo Order a OrderResponse
func mapOrderToResponse(order *models.Order) *domain.OrderResponse {
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
