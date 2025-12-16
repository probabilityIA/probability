package usecaseordermapping

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"gorm.io/datatypes"
)

// MapAndSaveOrder recibe una orden en formato canónico y la guarda en todas las tablas relacionadas
// Este es el punto de entrada principal para todas las integraciones después de mapear sus datos
func (uc *UseCaseOrderMapping) MapAndSaveOrder(ctx context.Context, dto *domain.ProbabilityOrderDTO) (*domain.OrderResponse, error) {
	// 0. Validar datos obligatorios de integración
	if dto.IntegrationID == 0 {
		return nil, errors.New("integration_id is required")
	}
	if dto.BusinessID == nil || *dto.BusinessID == 0 {
		return nil, errors.New("business_id is required")
	}

	// 1. Verificar si existe una orden con el mismo external_id para la misma integración
	exists, err := uc.repo.OrderExists(ctx, dto.ExternalID, dto.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("error checking if order exists: %w", err)
	}
	if exists {
		existingOrder, err := uc.repo.GetOrderByExternalID(ctx, dto.ExternalID, dto.IntegrationID)
		if err != nil {
			return nil, fmt.Errorf("error getting existing order: %w", err)
		}
		return uc.UpdateOrder(ctx, existingOrder, dto)
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

	// 2. Crear la entidad de dominio ProbabilityOrder
	order := &domain.ProbabilityOrder{
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
		OrderStatusURL:  dto.OrderStatusURL,

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

	// 2.3 POPULAR CAMPOS JSONB ITEMS SI ESTÁN VACÍOS (BACKWARD COMPATIBILITY)
	// Si items (JSONB) está vacío pero tenemos OrderItems (Relación), serializamos para llenar el campo JSONB
	if len(dto.OrderItems) > 0 && (dto.Items == nil || len(dto.Items) <= 4) {
		itemsJSON, err := json.Marshal(dto.OrderItems)
		if err == nil {
			order.Items = datatypes.JSON(itemsJSON)
		}
	}

	// 2.2 POPULAR CAMPOS PLANOS DE DIRECCIÓN (FLAT FIELDS)
	// Iterar sobre addresses para encontrar la de shipping y llenar los campos del modelo Order
	// Esto es crucial para compatibilidad con frontend que usa shipping_street, shipping_city, etc.
	for _, addr := range dto.Addresses {
		if addr.Type == "shipping" {
			order.ShippingStreet = addr.Street
			order.ShippingCity = addr.City
			order.ShippingState = addr.State
			order.ShippingCountry = addr.Country
			order.ShippingPostalCode = addr.PostalCode
			order.ShippingLat = addr.Latitude
			order.ShippingLng = addr.Longitude

			// Si hay complemento (Street2), lo concatenamos o guardamos en shipping_details si el modelo no tiene campo
			// El modelo actual no parece tener ShippingStreet2 plano, así que concatenamos o dependemos de la tabla addresses.
			// Por ahora, solo mapeamos lo que cabe.
			if addr.Street2 != "" {
				order.ShippingStreet = fmt.Sprintf("%s %s", order.ShippingStreet, addr.Street2)
				order.ShippingStreet2 = addr.Street2 // Populate for scoring
			}
			break
		}
	}

	// 3. Guardar la orden principal (sin score por ahora, se calculará mediante evento)
	if err := uc.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	fmt.Printf("[MapAndSaveOrder] Orden %s guardada exitosamente. Publicando evento para calcular score...\n", order.ID)

	// 4. Guardar OrderItems
	if len(dto.OrderItems) > 0 {
		orderItems := make([]*domain.ProbabilityOrderItem, len(dto.OrderItems))
		for i, itemDTO := range dto.OrderItems {
			// Validar/Crear Producto
			_, err := uc.GetOrCreateProduct(ctx, *dto.BusinessID, itemDTO)
			if err != nil {
				return nil, fmt.Errorf("error processing product for item %s: %w", itemDTO.ProductSKU, err)
			}

			orderItems[i] = &domain.ProbabilityOrderItem{
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
		addresses := make([]*domain.ProbabilityAddress, len(dto.Addresses))
		for i, addrDTO := range dto.Addresses {
			addresses[i] = &domain.ProbabilityAddress{
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
		payments := make([]*domain.ProbabilityPayment, len(dto.Payments))
		for i, payDTO := range dto.Payments {
			payments[i] = &domain.ProbabilityPayment{
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
		shipments := make([]*domain.ProbabilityShipment, len(dto.Shipments))
		for i, shipDTO := range dto.Shipments {
			shipments[i] = &domain.ProbabilityShipment{
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
		metadata := &domain.ProbabilityOrderChannelMetadata{
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

	// 9. Publicar eventos
	if uc.eventPublisher != nil {
		// 9.1. Publicar evento de orden creada
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
		go func() {
			// Usar context.Background() para evitar que el contexto cancelado afecte la publicación
			bgCtx := context.Background()
			uc.logger.Info(bgCtx).
				Str("order_id", order.ID).
				Str("event_type", string(event.Type)).
				Interface("business_id", event.BusinessID).
				Interface("integration_id", event.IntegrationID).
				Str("order_number", order.OrderNumber).
				Msg("📤 Publicando evento order.created a Redis...")
			if err := uc.eventPublisher.PublishOrderEvent(bgCtx, event); err != nil {
				uc.logger.Error(bgCtx).
					Err(err).
					Str("order_id", order.ID).
					Str("event_type", string(event.Type)).
					Msg("❌ Error al publicar evento de orden creada")
			} else {
				uc.logger.Info(bgCtx).
					Str("order_id", order.ID).
					Str("event_type", string(event.Type)).
					Msg("✅ Evento order.created publicado exitosamente a Redis")
			}
		}()

		// 9.2. Calcular score directamente cuando llega por RabbitMQ
		go func() {
			fmt.Printf("[MapAndSaveOrder] Calculando score directamente para orden %s\n", order.ID)
			if err := uc.scoreUseCase.CalculateAndUpdateOrderScore(ctx, order.ID); err != nil {
				uc.logger.Error(ctx).
					Err(err).
					Str("order_id", order.ID).
					Msg("Error al calcular score de la orden")
			} else {
				uc.logger.Info(ctx).
					Str("order_id", order.ID).
					Str("order_number", order.OrderNumber).
					Msg("✅ Score calculado exitosamente para la orden")
			}
		}()

		// 9.3. Publicar evento para calcular score (mantener para otros consumidores)
		scoreEventData := domain.OrderEventData{
			OrderNumber:    order.OrderNumber,
			InternalNumber: order.InternalNumber,
			ExternalID:     order.ExternalID,
		}
		scoreEvent := domain.NewOrderEvent(domain.OrderEventTypeScoreCalculationRequested, order.ID, scoreEventData)
		scoreEvent.BusinessID = order.BusinessID
		if order.IntegrationID > 0 {
			integrationID := order.IntegrationID
			scoreEvent.IntegrationID = &integrationID
		}
		go func() {
			fmt.Printf("[MapAndSaveOrder] Publicando evento order.score_calculation_requested para orden %s\n", order.ID)
			if err := uc.eventPublisher.PublishOrderEvent(ctx, scoreEvent); err != nil {
				uc.logger.Error(ctx).
					Err(err).
					Str("order_id", order.ID).
					Msg("Error al publicar evento de cálculo de score")
			} else {
				fmt.Printf("[MapAndSaveOrder] Evento order.score_calculation_requested publicado exitosamente para orden %s\n", order.ID)
			}
		}()
	}

	// 10. Retornar la respuesta mapeada
	return mapOrderToResponse(order), nil
}

// mapOrderToResponse convierte un modelo Order a OrderResponse
func mapOrderToResponse(order *domain.ProbabilityOrder) *domain.OrderResponse {
	return &domain.OrderResponse{
		ID:        order.ID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,

		// Identificadores de integración
		BusinessID:         order.BusinessID,
		IntegrationID:      order.IntegrationID,
		IntegrationType:    order.IntegrationType,
		IntegrationLogoURL: order.IntegrationLogoURL,

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
		TrackingNumber:      order.TrackingNumber,
		TrackingLink:        order.TrackingLink,
		GuideID:             order.GuideID,
		GuideLink:           order.GuideLink,
		DeliveryDate:        order.DeliveryDate,
		DeliveredAt:         order.DeliveredAt,
		DeliveryProbability: order.DeliveryProbability,

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
		OrderStatusURL:  order.OrderStatusURL,

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
