package usecaseordermapping

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	integrationevents "github.com/secamc93/probability/back/central/services/integrations/events"
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
	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "D",
			"location":     "map-order.go:26",
			"message":      "MapAndSaveOrder - Checking if order exists",
			"data": map[string]interface{}{
				"external_id":    dto.ExternalID,
				"order_number":   dto.OrderNumber,
				"integration_id": dto.IntegrationID,
			},
			"timestamp": time.Now().UnixMilli(),
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion
	exists, err := uc.repo.OrderExists(ctx, dto.ExternalID, dto.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("error checking if order exists: %w", err)
	}
	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "D",
			"location":     "map-order.go:30",
			"message":      "MapAndSaveOrder - OrderExists result",
			"data": map[string]interface{}{
				"external_id":    dto.ExternalID,
				"order_number":   dto.OrderNumber,
				"integration_id": dto.IntegrationID,
				"exists":         exists,
			},
			"timestamp": time.Now().UnixMilli(),
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion
	if exists {
		// #region agent log
		if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			logData, _ := json.Marshal(map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "D",
				"location":     "map-order.go:35",
				"message":      "MapAndSaveOrder - Order exists, UPDATING instead of creating",
				"data": map[string]interface{}{
					"external_id":    dto.ExternalID,
					"order_number":   dto.OrderNumber,
					"integration_id": dto.IntegrationID,
				},
				"timestamp": time.Now().UnixMilli(),
			})
			f.WriteString(string(logData) + "\n")
			f.Close()
		}
		// #endregion
		existingOrder, err := uc.repo.GetOrderByExternalID(ctx, dto.ExternalID, dto.IntegrationID)
		if err != nil {
			return nil, fmt.Errorf("error getting existing order: %w", err)
		}
		return uc.UpdateOrder(ctx, existingOrder, dto)
	}
	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "D",
			"location":     "map-order.go:37",
			"message":      "MapAndSaveOrder - Order does NOT exist, will CREATE new",
			"data": map[string]interface{}{
				"external_id":    dto.ExternalID,
				"order_number":   dto.OrderNumber,
				"integration_id": dto.IntegrationID,
			},
			"timestamp": time.Now().UnixMilli(),
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	// 1.5. Validar/Crear Cliente
	client, err := uc.GetOrCreateCustomer(ctx, *dto.BusinessID, dto)
	if err != nil {
		// Publicar evento de orden rechazada por error en cliente
		integrationevents.PublishSyncOrderRejected(
			ctx,
			dto.IntegrationID,
			dto.BusinessID,
			"", // orderID aún no existe
			dto.OrderNumber,
			dto.ExternalID,
			dto.Platform,
			"Error al procesar cliente",
			err.Error(),
		)
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
	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "E",
			"location":     "map-order.go:64",
			"message":      "MapAndSaveOrder - Attempting to CREATE order in database",
			"data": map[string]interface{}{
				"external_id":    dto.ExternalID,
				"order_number":   dto.OrderNumber,
				"integration_id": dto.IntegrationID,
			},
			"timestamp": time.Now().UnixMilli(),
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion
	if err := uc.repo.CreateOrder(ctx, order); err != nil {
		// #region agent log
		if f, err2 := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			logData, _ := json.Marshal(map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "E",
				"location":     "map-order.go:66",
				"message":      "MapAndSaveOrder - ERROR creating order in database",
				"data": map[string]interface{}{
					"external_id":    dto.ExternalID,
					"order_number":   dto.OrderNumber,
					"integration_id": dto.IntegrationID,
					"error":          err.Error(),
				},
				"timestamp": time.Now().UnixMilli(),
			})
			f.WriteString(string(logData) + "\n")
			f.Close()
		}
		// #endregion

		// Publicar evento de orden rechazada
		integrationevents.PublishSyncOrderRejected(
			ctx,
			dto.IntegrationID,
			dto.BusinessID,
			"", // orderID aún no existe
			dto.OrderNumber,
			dto.ExternalID,
			dto.Platform,
			"Error al guardar en base de datos",
			err.Error(),
		)

		return nil, fmt.Errorf("error creating order: %w", err)
	}
	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "E",
			"location":     "map-order.go:68",
			"message":      "MapAndSaveOrder - Order CREATED successfully in database",
			"data": map[string]interface{}{
				"order_id":       order.ID,
				"external_id":    dto.ExternalID,
				"order_number":   dto.OrderNumber,
				"integration_id": dto.IntegrationID,
			},
			"timestamp": time.Now().UnixMilli(),
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	fmt.Printf("[MapAndSaveOrder] Orden %s guardada exitosamente. Publicando evento para calcular score...\n", order.ID)

	// Publicar evento de orden creada exitosamente
	integrationevents.PublishSyncOrderCreated(
		ctx,
		dto.IntegrationID,
		dto.BusinessID,
		order.ID,
		dto.OrderNumber,
		dto.ExternalID,
		dto.Platform,
		dto.CustomerEmail,
		dto.Currency,
		order.Status,
		order.CreatedAt,
		&order.TotalAmount,
	)

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
func (uc *UseCaseOrderMapping) mapOrderStatuses(ctx context.Context, dto *domain.ProbabilityOrderDTO) orderStatusMapping {
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
func (uc *UseCaseOrderMapping) mapOrderStatusID(ctx context.Context, dto *domain.ProbabilityOrderDTO) *uint {
	if dto.IntegrationType == "" {
		return nil
	}

	integrationTypeID := getIntegrationTypeID(dto.IntegrationType)
	if integrationTypeID == 0 {
		return nil
	}

	// Prioridad 1: Intentar mapear usando Status (puede ser shipment_status, fulfillment_status, o financial_status)
	if dto.Status != "" {
		mappedStatusID, err := uc.orderStatusRepository.GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx, integrationTypeID, dto.Status)
		if err == nil && mappedStatusID != nil {
			return mappedStatusID
		}
	}

	// Prioridad 2: Si no encontró con Status, intentar con OriginalStatus (financial_status)
	if dto.OriginalStatus != "" {
		mappedStatusID, err := uc.orderStatusRepository.GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx, integrationTypeID, dto.OriginalStatus)
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
func (uc *UseCaseOrderMapping) mapPaymentStatusID(ctx context.Context, dto *domain.ProbabilityOrderDTO) *uint {
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

		paymentStatusCode := mapShopifyFinancialStatusToPaymentStatus(financialStatus)
		mappedID, err := uc.paymentStatusRepository.GetPaymentStatusIDByCode(ctx, paymentStatusCode)
		if err != nil || mappedID == nil {
			return nil
		}

		return mappedID
	}

	return nil
}

// mapFulfillmentStatusID mapea el estado de fulfillment desde el DTO o desde FulfillmentDetails si es Shopify
func (uc *UseCaseOrderMapping) mapFulfillmentStatusID(ctx context.Context, dto *domain.ProbabilityOrderDTO) *uint {
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

		fulfillmentStatusCode := mapShopifyFulfillmentStatusToFulfillmentStatus(&fulfillmentStatus)
		return uc.getFulfillmentStatusIDByCode(ctx, fulfillmentStatusCode)
	}

	return nil
}

// getFulfillmentStatusIDByCode obtiene el ID de un estado de fulfillment por su código
func (uc *UseCaseOrderMapping) getFulfillmentStatusIDByCode(ctx context.Context, code string) *uint {
	mappedID, err := uc.fulfillmentStatusRepository.GetFulfillmentStatusIDByCode(ctx, code)
	if err != nil || mappedID == nil {
		return nil
	}
	return mappedID
}

// syncIsPaidFromPaymentStatus sincroniza IsPaid basado en PaymentStatusID
func (uc *UseCaseOrderMapping) syncIsPaidFromPaymentStatus(ctx context.Context, order *domain.ProbabilityOrder, paymentStatusID *uint) {
	if paymentStatusID == nil {
		return
	}

	paidStatusID, err := uc.paymentStatusRepository.GetPaymentStatusIDByCode(ctx, "paid")
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
func (uc *UseCaseOrderMapping) buildOrderEntity(dto *domain.ProbabilityOrderDTO, clientID *uint, statusMapping orderStatusMapping) *domain.ProbabilityOrder {
	return &domain.ProbabilityOrder{
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
func (uc *UseCaseOrderMapping) assignPaymentMethodID(order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) {
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
func (uc *UseCaseOrderMapping) populateOrderFields(order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) {
	// Popular campos JSONB Items si están vacíos (backward compatibility)
	if len(dto.OrderItems) > 0 && len(dto.Items) <= 4 {
		itemsJSON, err := json.Marshal(dto.OrderItems)
		if err == nil {
			order.Items = datatypes.JSON(itemsJSON)
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
				fmt.Printf("[DEBUG ADDRESS] Found Street2: '%s' for Order %s. Appending to ShippingStreet.\n", addr.Street2, order.OrderNumber)
				order.ShippingStreet = fmt.Sprintf("%s %s", order.ShippingStreet, addr.Street2)
				order.Address2 = addr.Street2 // Populate for scoring
			} else {
				fmt.Printf("[DEBUG ADDRESS] Street2 is EMPTY for Order %s. Address2 will optionally be empty.\n", order.OrderNumber)
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
func (uc *UseCaseOrderMapping) saveRelatedEntities(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) error {
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
func (uc *UseCaseOrderMapping) saveOrderItems(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) error {
	if len(dto.OrderItems) == 0 {
		return nil
	}

	orderItems := make([]*domain.ProbabilityOrderItem, len(dto.OrderItems))
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

		orderItems[i] = &domain.ProbabilityOrderItem{
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
func (uc *UseCaseOrderMapping) saveAddresses(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) error {
	if len(dto.Addresses) == 0 {
		return nil
	}

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

	return uc.repo.CreateAddresses(ctx, addresses)
}

// savePayments guarda los pagos de la orden
func (uc *UseCaseOrderMapping) savePayments(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) error {
	if len(dto.Payments) == 0 {
		return nil
	}

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

	return uc.repo.CreatePayments(ctx, payments)
}

// saveShipments guarda los envíos de la orden
func (uc *UseCaseOrderMapping) saveShipments(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) error {
	if len(dto.Shipments) == 0 {
		return nil
	}

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

	return uc.repo.CreateShipments(ctx, shipments)
}

// saveChannelMetadata guarda los metadatos del canal
func (uc *UseCaseOrderMapping) saveChannelMetadata(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) error {
	if dto.ChannelMetadata == nil {
		return nil
	}

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

	return uc.repo.CreateChannelMetadata(ctx, metadata)
}

// ───────────────────────────────────────────
//
//	FUNCIONES DE EVENTOS
//
// ───────────────────────────────────────────

// publishOrderEvents publica los eventos relacionados con la orden creada
func (uc *UseCaseOrderMapping) publishOrderEvents(ctx context.Context, order *domain.ProbabilityOrder) {
	if uc.eventPublisher == nil {
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
func (uc *UseCaseOrderMapping) publishOrderCreatedEvent(_ context.Context, order *domain.ProbabilityOrder) {
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
}

// calculateOrderScore calcula el score de la orden directamente
func (uc *UseCaseOrderMapping) calculateOrderScore(ctx context.Context, order *domain.ProbabilityOrder) {
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
}

// publishScoreCalculationEvent publica el evento para calcular score
func (uc *UseCaseOrderMapping) publishScoreCalculationEvent(ctx context.Context, order *domain.ProbabilityOrder) {
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

// ───────────────────────────────────────────
//
//	FUNCIONES DE MAPEO DE RESPUESTA
//
// ───────────────────────────────────────────

// mapOrderToResponse convierte un modelo Order a OrderResponse
func (uc *UseCaseOrderMapping) mapOrderToResponse(order *domain.ProbabilityOrder) *domain.OrderResponse {
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
