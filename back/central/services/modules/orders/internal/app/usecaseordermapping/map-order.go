package usecaseordermapping

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/helpers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseordermapping/mappers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// MapAndSaveOrder recibe una orden en formato canónico y la guarda en todas las tablas relacionadas
// Este es el punto de entrada principal para todas las integraciones después de mapear sus datos
func (uc *UseCaseOrderMapping) MapAndSaveOrder(ctx context.Context, dto *dtos.ProbabilityOrderDTO) (*dtos.OrderResponse, error) {
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
		uc.publishSyncOrderRejected(ctx, dto.IntegrationID, dto.BusinessID, dto.OrderNumber, dto.ExternalID, dto.Platform, "Error al procesar cliente", err.Error())
		return nil, fmt.Errorf("error processing customer: %w", err)
	}
	var clientID *uint
	if client != nil {
		clientID = &client.ID
	}

	// 1.6. Mapear estados de la orden
	statusMapping := uc.mapOrderStatuses(ctx, dto)

	// 2. Crear la entidad de dominio ProbabilityOrder
	order := uc.buildOrderEntity(dto, clientID, statusMapping)

	// 2.1. Asignar PaymentMethodID desde el primer pago
	uc.assignPaymentMethodID(order, dto)

	// 2.1.1. Mantener IsPaid actualizado según PaymentStatusID
	uc.syncIsPaidFromPaymentStatus(ctx, order, statusMapping.PaymentStatusID)

	// 2.3. Popular campos JSONB y planos de dirección
	uc.populateOrderFields(order, dto)

	// 3. Guardar la orden principal (sin score por ahora, se calculará mediante evento)
	if err := uc.repo.CreateOrder(ctx, order); err != nil {
		uc.publishSyncOrderRejected(ctx, dto.IntegrationID, dto.BusinessID, dto.OrderNumber, dto.ExternalID, dto.Platform, "Error al guardar en base de datos", err.Error())

		return nil, fmt.Errorf("error creating order: %w", err)
	}

	// Publicar evento de orden creada exitosamente
	uc.publishSyncOrderCreated(ctx, dto.IntegrationID, dto.BusinessID, map[string]interface{}{
		"order_id":       order.ID,
		"order_number":   dto.OrderNumber,
		"external_id":    dto.ExternalID,
		"platform":       dto.Platform,
		"customer_email": dto.CustomerEmail,
		"currency":       dto.Currency,
		"status":         order.Status,
		"created_at":     order.CreatedAt,
		"total_amount":   order.TotalAmount,
		"synced_at":      time.Now(),
	})

	// 4-8. Guardar entidades relacionadas
	if err := uc.saveRelatedEntities(ctx, order, dto); err != nil {
		return nil, err
	}

	// 9. Publicar eventos
	uc.publishOrderEvents(ctx, order)

	// 10. Retornar la respuesta mapeada
	return uc.mapOrderToResponse(order), nil
}

// ───────────────────────────────────────────
//
//	FUNCIONES DE MAPEO DE ESTADOS
//
// ───────────────────────────────────────────

// orderStatusMapping contiene los IDs de los estados mapeados
type orderStatusMapping struct {
	OrderStatusID       *uint
	PaymentStatusID     *uint
	FulfillmentStatusID *uint
}

// mapOrderStatuses mapea todos los estados de la orden (OrderStatus, PaymentStatus, FulfillmentStatus)
func (uc *UseCaseOrderMapping) mapOrderStatuses(ctx context.Context, dto *dtos.ProbabilityOrderDTO) orderStatusMapping {
	mapping := orderStatusMapping{}

	// Mapear OrderStatusID
	mapping.OrderStatusID = uc.mapOrderStatusID(ctx, dto)

	// Mapear PaymentStatusID
	mapping.PaymentStatusID = uc.mapPaymentStatusID(ctx, dto)

	// Mapear FulfillmentStatusID
	mapping.FulfillmentStatusID = uc.mapFulfillmentStatusID(ctx, dto)

	return mapping
}

// mapOrderStatusID mapea el estado general de la orden
// Prioridad 1: Intentar mapear usando Status (que será shipment_status cuando existe)
// Prioridad 2: Si no encuentra con Status, intentar con OriginalStatus (financial_status) como fallback
func (uc *UseCaseOrderMapping) mapOrderStatusID(ctx context.Context, dto *dtos.ProbabilityOrderDTO) *uint {
	if dto.IntegrationType == "" {
		return nil
	}

	integrationTypeID := mappers.GetIntegrationTypeID(dto.IntegrationType)
	if integrationTypeID == 0 {
		return nil
	}

	// Prioridad 1: Intentar mapear usando Status (puede ser shipment_status, fulfillment_status, o financial_status)
	if dto.Status != "" {
		mappedStatusID, err := uc.repo.GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx, integrationTypeID, dto.Status)
		if err == nil && mappedStatusID != nil {
			return mappedStatusID
		}
	}

	// Prioridad 2: Si no encontró con Status, intentar con OriginalStatus (financial_status)
	if dto.OriginalStatus != "" {
		mappedStatusID, err := uc.repo.GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx, integrationTypeID, dto.OriginalStatus)
		if err != nil {
			uc.logger.Warn().
				Uint("integration_type_id", integrationTypeID).
				Str("status", dto.Status).
				Str("original_status", dto.OriginalStatus).
				Err(err).
				Msg("Error al buscar mapeo de estado, continuando sin status_id")
			return nil
		}
		return mappedStatusID
	}

	return nil
}

// mapPaymentStatusID mapea el estado de pago desde el DTO o desde PaymentDetails si es Shopify
func (uc *UseCaseOrderMapping) mapPaymentStatusID(ctx context.Context, dto *dtos.ProbabilityOrderDTO) *uint {
	// Si el DTO ya tiene PaymentStatusID, usarlo directamente
	if dto.PaymentStatusID != nil && *dto.PaymentStatusID > 0 {
		return dto.PaymentStatusID
	}

	// Para Shopify, extraer desde PaymentDetails
	if dto.IntegrationType == "shopify" && len(dto.PaymentDetails) > 0 {
		var paymentDetails map[string]interface{}
		if err := json.Unmarshal(dto.PaymentDetails, &paymentDetails); err != nil {
			return nil
		}

		financialStatus, ok := paymentDetails["financial_status"].(string)
		if !ok || financialStatus == "" {
			return nil
		}

		paymentStatusCode := mappers.MapShopifyFinancialStatusToPaymentStatus(financialStatus)
		mappedID, err := uc.repo.GetPaymentStatusIDByCode(ctx, paymentStatusCode)
		if err != nil || mappedID == nil {
			return nil
		}

		return mappedID
	}

	return nil
}

// mapFulfillmentStatusID mapea el estado de fulfillment desde el DTO o desde FulfillmentDetails si es Shopify
func (uc *UseCaseOrderMapping) mapFulfillmentStatusID(ctx context.Context, dto *dtos.ProbabilityOrderDTO) *uint {
	// Si el DTO ya tiene FulfillmentStatusID, usarlo directamente
	if dto.FulfillmentStatusID != nil && *dto.FulfillmentStatusID > 0 {
		return dto.FulfillmentStatusID
	}

	// Para Shopify, extraer desde FulfillmentDetails
	if dto.IntegrationType == "shopify" && len(dto.FulfillmentDetails) > 0 {
		var fulfillmentDetails map[string]interface{}
		if err := json.Unmarshal(dto.FulfillmentDetails, &fulfillmentDetails); err != nil {
			// Si no se puede parsear, usar "unfulfilled" por defecto
			return uc.getFulfillmentStatusIDByCode(ctx, "unfulfilled")
		}

		fulfillmentStatus, ok := fulfillmentDetails["fulfillment_status"].(string)
		if !ok || fulfillmentStatus == "" {
			// Si fulfillment_status es null o vacío, usar "unfulfilled"
			return uc.getFulfillmentStatusIDByCode(ctx, "unfulfilled")
		}

		fulfillmentStatusCode := mappers.MapShopifyFulfillmentStatusToFulfillmentStatus(&fulfillmentStatus)
		return uc.getFulfillmentStatusIDByCode(ctx, fulfillmentStatusCode)
	}

	return nil
}

// getFulfillmentStatusIDByCode obtiene el ID de un estado de fulfillment por su código
func (uc *UseCaseOrderMapping) getFulfillmentStatusIDByCode(ctx context.Context, code string) *uint {
	mappedID, err := uc.repo.GetFulfillmentStatusIDByCode(ctx, code)
	if err != nil || mappedID == nil {
		return nil
	}
	return mappedID
}

// syncIsPaidFromPaymentStatus sincroniza IsPaid basado en PaymentStatusID
func (uc *UseCaseOrderMapping) syncIsPaidFromPaymentStatus(ctx context.Context, order *entities.ProbabilityOrder, paymentStatusID *uint) {
	if paymentStatusID == nil {
		return
	}

	paidStatusID, err := uc.repo.GetPaymentStatusIDByCode(ctx, "paid")
	if err != nil || paidStatusID == nil || *paidStatusID != *paymentStatusID {
		return
	}

	order.IsPaid = true
	if order.PaidAt == nil {
		now := time.Now()
		order.PaidAt = &now
	}
}

// ───────────────────────────────────────────
//
//	FUNCIONES DE CONSTRUCCIÓN DE ENTIDADES
//
// ───────────────────────────────────────────

// buildOrderEntity construye la entidad ProbabilityOrder desde el DTO
func (uc *UseCaseOrderMapping) buildOrderEntity(dto *dtos.ProbabilityOrderDTO, clientID *uint, statusMapping orderStatusMapping) *entities.ProbabilityOrder {
	return &entities.ProbabilityOrder{
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

		// Información del cliente
		CustomerID:    clientID,
		CustomerName:  dto.CustomerName,
		CustomerEmail: dto.CustomerEmail,
		CustomerPhone: dto.CustomerPhone,
		CustomerDNI:   dto.CustomerDNI,
		CustomerOrderCount: func() int {
			if dto.CustomerOrderCount != nil {
				return *dto.CustomerOrderCount
			}
			return 0
		}(),
		CustomerTotalSpent: func() string {
			if dto.CustomerTotalSpent != nil {
				return *dto.CustomerTotalSpent
			}
			return ""
		}(),

		// Tipo y estado
		OrderTypeID:         dto.OrderTypeID,
		OrderTypeName:       dto.OrderTypeName,
		Status:              dto.Status,
		OriginalStatus:      dto.OriginalStatus,
		StatusID:            statusMapping.OrderStatusID,
		PaymentStatusID:     statusMapping.PaymentStatusID,
		FulfillmentStatusID: statusMapping.FulfillmentStatusID,

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
}

// assignPaymentMethodID asigna el PaymentMethodID desde el primer pago
func (uc *UseCaseOrderMapping) assignPaymentMethodID(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) {
	order.PaymentMethodID = 1 // Valor por defecto

	if len(dto.Payments) == 0 {
		return
	}

	payment := dto.Payments[0]
	if payment.PaymentMethodID > 0 {
		order.PaymentMethodID = payment.PaymentMethodID
		if payment.Status == "completed" && payment.PaidAt != nil {
			order.IsPaid = true
			order.PaidAt = payment.PaidAt
		}
	}
}

// populateOrderFields popula campos JSONB y campos planos de dirección
func (uc *UseCaseOrderMapping) populateOrderFields(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) {
	// Popular campos JSONB Items si están vacíos (backward compatibility)
	if len(dto.OrderItems) > 0 && len(dto.Items) <= 4 {
		itemsJSON, err := json.Marshal(dto.OrderItems)
		if err == nil {
			order.Items = itemsJSON
		}
	}

	// Popular campos planos de dirección (flat fields)
	for _, addr := range dto.Addresses {
		if addr.Type == "shipping" {
			order.ShippingStreet = addr.Street
			order.ShippingCity = addr.City
			order.ShippingState = addr.State
			order.ShippingCountry = addr.Country
			order.ShippingPostalCode = addr.PostalCode
			order.ShippingLat = addr.Latitude
			order.ShippingLng = addr.Longitude

			if addr.Street2 != "" {
				order.ShippingStreet = fmt.Sprintf("%s %s", order.ShippingStreet, addr.Street2)
				order.Address2 = addr.Street2 // Populate for scoring
			}
			break
		}
	}
}

// ───────────────────────────────────────────
//
//	FUNCIONES DE PERSISTENCIA
//
// ───────────────────────────────────────────

// saveRelatedEntities guarda todas las entidades relacionadas (items, addresses, payments, shipments, metadata)
func (uc *UseCaseOrderMapping) saveRelatedEntities(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if err := uc.saveOrderItems(ctx, order, dto); err != nil {
		return err
	}

	if err := uc.saveAddresses(ctx, order, dto); err != nil {
		return err
	}

	if err := uc.savePayments(ctx, order, dto); err != nil {
		return err
	}

	if err := uc.saveShipments(ctx, order, dto); err != nil {
		return err
	}

	if err := uc.saveChannelMetadata(ctx, order, dto); err != nil {
		return err
	}

	return nil
}

// saveOrderItems guarda los items de la orden
func (uc *UseCaseOrderMapping) saveOrderItems(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if len(dto.OrderItems) == 0 {
		return nil
	}

	orderItems := make([]*entities.ProbabilityOrderItem, len(dto.OrderItems))
	for i, itemDTO := range dto.OrderItems {
		// Validar/Crear Producto
		product, err := uc.GetOrCreateProduct(ctx, *dto.BusinessID, itemDTO)
		if err != nil {
			return fmt.Errorf("error processing product for item %s: %w", itemDTO.ProductSKU, err)
		}

		var productID *string
		if product != nil {
			productID = &product.ID
		}

		orderItems[i] = &entities.ProbabilityOrderItem{
			OrderID:          order.ID,
			ProductID:        productID,
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
			// Precios en moneda local
			UnitPricePresentment:  itemDTO.UnitPricePresentment,
			TotalPricePresentment: itemDTO.TotalPricePresentment,
			DiscountPresentment:   itemDTO.DiscountPresentment,
			TaxPresentment:        itemDTO.TaxPresentment,
		}
	}

	return uc.repo.CreateOrderItems(ctx, orderItems)
}

// saveAddresses guarda las direcciones de la orden
func (uc *UseCaseOrderMapping) saveAddresses(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if len(dto.Addresses) == 0 {
		return nil
	}

	addresses := make([]*entities.ProbabilityAddress, len(dto.Addresses))
	for i, addrDTO := range dto.Addresses {
		addresses[i] = &entities.ProbabilityAddress{
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

	return uc.repo.CreateAddresses(ctx, addresses)
}

// savePayments guarda los pagos de la orden
func (uc *UseCaseOrderMapping) savePayments(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if len(dto.Payments) == 0 {
		return nil
	}

	payments := make([]*entities.ProbabilityPayment, len(dto.Payments))
	for i, payDTO := range dto.Payments {
		payments[i] = &entities.ProbabilityPayment{
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

	return uc.repo.CreatePayments(ctx, payments)
}

// saveShipments guarda los envíos de la orden
func (uc *UseCaseOrderMapping) saveShipments(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if len(dto.Shipments) == 0 {
		return nil
	}

	shipments := make([]*entities.ProbabilityShipment, len(dto.Shipments))
	for i, shipDTO := range dto.Shipments {
		shipments[i] = &entities.ProbabilityShipment{
			OrderID:           &order.ID,
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

	return uc.repo.CreateShipments(ctx, shipments)
}

// saveChannelMetadata guarda los metadatos del canal
func (uc *UseCaseOrderMapping) saveChannelMetadata(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if dto.ChannelMetadata == nil {
		return nil
	}

	metadata := &entities.ProbabilityOrderChannelMetadata{
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

	return uc.repo.CreateChannelMetadata(ctx, metadata)
}

// ───────────────────────────────────────────
//
//	FUNCIONES DE EVENTOS
//
// ───────────────────────────────────────────

// publishOrderEvents publica los eventos relacionados con la orden creada
func (uc *UseCaseOrderMapping) publishOrderEvents(ctx context.Context, order *entities.ProbabilityOrder) {
	if uc.redisEventPublisher == nil {
		return
	}

	// Publicar evento de orden creada
	uc.publishOrderCreatedEvent(ctx, order)

	// Calcular score directamente
	uc.calculateOrderScore(ctx, order)

	// Publicar evento para calcular score (para otros consumidores)
	uc.publishScoreCalculationEvent(ctx, order)
}

// publishOrderCreatedEvent publica el evento de orden creada
func (uc *UseCaseOrderMapping) publishOrderCreatedEvent(_ context.Context, order *entities.ProbabilityOrder) {
	eventData := entities.OrderEventData{
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,
		ExternalID:     order.ExternalID,
		CurrentStatus:  order.Status,
		CustomerEmail:  order.CustomerEmail,
		TotalAmount:    &order.TotalAmount,
		Currency:       order.Currency,
		Platform:       order.Platform,
	}

	event := entities.NewOrderEvent(entities.OrderEventTypeCreated, order.ID, eventData)
	event.BusinessID = order.BusinessID
	if order.IntegrationID > 0 {
		integrationID := order.IntegrationID
		event.IntegrationID = &integrationID
	}

	// Publicar en ambos canales (Redis + RabbitMQ) con orden completa
	helpers.PublishEventDual(context.Background(), event, order, uc.redisEventPublisher, uc.rabbitEventPublisher, uc.logger)
}

// calculateOrderScore calcula el score de la orden directamente
func (uc *UseCaseOrderMapping) calculateOrderScore(ctx context.Context, order *entities.ProbabilityOrder) {
	go func() {
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
}

// publishScoreCalculationEvent publica el evento para calcular score
func (uc *UseCaseOrderMapping) publishScoreCalculationEvent(ctx context.Context, order *entities.ProbabilityOrder) {
	scoreEventData := entities.OrderEventData{
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,
		ExternalID:     order.ExternalID,
	}

	scoreEvent := entities.NewOrderEvent(entities.OrderEventTypeScoreCalculationRequested, order.ID, scoreEventData)
	scoreEvent.BusinessID = order.BusinessID
	if order.IntegrationID > 0 {
		integrationID := order.IntegrationID
		scoreEvent.IntegrationID = &integrationID
	}

	// Publicar en ambos canales (Redis + RabbitMQ) con orden completa
	helpers.PublishEventDual(context.Background(), scoreEvent, order, uc.redisEventPublisher, uc.rabbitEventPublisher, uc.logger)
}

// ───────────────────────────────────────────
//
//	FUNCIONES DE MAPEO DE RESPUESTA
//
// ───────────────────────────────────────────

// mapOrderToResponse convierte un modelo Order a OrderResponse
func (uc *UseCaseOrderMapping) mapOrderToResponse(order *entities.ProbabilityOrder) *dtos.OrderResponse {
	return &dtos.OrderResponse{
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
		Subtotal:                order.Subtotal,
		Tax:                     order.Tax,
		Discount:                order.Discount,
		ShippingCost:            order.ShippingCost,
		TotalAmount:             order.TotalAmount,
		Currency:                order.Currency,
		CodTotal:                order.CodTotal,
		SubtotalPresentment:     order.SubtotalPresentment,
		TaxPresentment:          order.TaxPresentment,
		DiscountPresentment:     order.DiscountPresentment,
		ShippingCostPresentment: order.ShippingCostPresentment,
		TotalAmountPresentment:  order.TotalAmountPresentment,
		CurrencyPresentment:     order.CurrencyPresentment,

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
		OrderTypeID:         order.OrderTypeID,
		OrderTypeName:       order.OrderTypeName,
		Status:              order.Status,
		OriginalStatus:      order.OriginalStatus,
		StatusID:            order.StatusID,
		PaymentStatusID:     order.PaymentStatusID,
		FulfillmentStatusID: order.FulfillmentStatusID,
		OrderStatus:         order.OrderStatus,
		PaymentStatus:       order.PaymentStatus,
		FulfillmentStatus:   order.FulfillmentStatus,

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
