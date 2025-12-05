package usecases

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ───────────────────────────────────────────
//
//	CREATE ORDER
//
// ───────────────────────────────────────────

// CreateOrder crea una nueva orden
func (uc *UseCases) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.OrderResponse, error) {
	// Validar que no exista una orden con el mismo external_id para la misma integración
	exists, err := uc.repo.OrderExists(ctx, req.ExternalID, req.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("error checking if order exists: %w", err)
	}
	if exists {
		return nil, errors.New("order with this external_id already exists for this integration")
	}

	// Crear el modelo de orden
	order := &models.Order{
		// Identificadores de integración
		BusinessID:      req.BusinessID,
		IntegrationID:   req.IntegrationID,
		IntegrationType: req.IntegrationType,

		// Identificadores de la orden
		Platform:       req.Platform,
		ExternalID:     req.ExternalID,
		OrderNumber:    req.OrderNumber,
		InternalNumber: req.InternalNumber,

		// Información financiera
		Subtotal:     req.Subtotal,
		Tax:          req.Tax,
		Discount:     req.Discount,
		ShippingCost: req.ShippingCost,
		TotalAmount:  req.TotalAmount,
		Currency:     req.Currency,
		CodTotal:     req.CodTotal,

		// Información del cliente
		CustomerID:    req.CustomerID,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
		CustomerDNI:   req.CustomerDNI,

		// Dirección de envío
		ShippingStreet:     req.ShippingStreet,
		ShippingCity:       req.ShippingCity,
		ShippingState:      req.ShippingState,
		ShippingCountry:    req.ShippingCountry,
		ShippingPostalCode: req.ShippingPostalCode,
		ShippingLat:        req.ShippingLat,
		ShippingLng:        req.ShippingLng,

		// Información de pago
		PaymentMethodID: req.PaymentMethodID,
		IsPaid:          req.IsPaid,
		PaidAt:          req.PaidAt,

		// Información de envío/logística
		TrackingNumber: req.TrackingNumber,
		TrackingLink:   req.TrackingLink,
		GuideID:        req.GuideID,
		GuideLink:      req.GuideLink,
		DeliveryDate:   req.DeliveryDate,
		DeliveredAt:    req.DeliveredAt,

		// Información de fulfillment
		WarehouseID:   req.WarehouseID,
		WarehouseName: req.WarehouseName,
		DriverID:      req.DriverID,
		DriverName:    req.DriverName,
		IsLastMile:    req.IsLastMile,

		// Dimensiones y peso
		Weight: req.Weight,
		Height: req.Height,
		Width:  req.Width,
		Length: req.Length,
		Boxes:  req.Boxes,

		// Tipo y estado
		OrderTypeID:    req.OrderTypeID,
		OrderTypeName:  req.OrderTypeName,
		Status:         req.Status,
		OriginalStatus: req.OriginalStatus,

		// Información adicional
		Notes:    req.Notes,
		Coupon:   req.Coupon,
		Approved: req.Approved,
		UserID:   req.UserID,
		UserName: req.UserName,

		// Facturación
		Invoiceable:     req.Invoiceable,
		InvoiceURL:      req.InvoiceURL,
		InvoiceID:       req.InvoiceID,
		InvoiceProvider: req.InvoiceProvider,

		// Datos estructurados
		Items:              req.Items,
		Metadata:           req.Metadata,
		FinancialDetails:   req.FinancialDetails,
		ShippingDetails:    req.ShippingDetails,
		PaymentDetails:     req.PaymentDetails,
		FulfillmentDetails: req.FulfillmentDetails,

		// Timestamps
		OccurredAt: req.OccurredAt,
		ImportedAt: req.ImportedAt,
	}

	// Guardar en la base de datos
	if err := uc.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	// Retornar la respuesta
	return mapOrderToResponse(order), nil
}

// ───────────────────────────────────────────
//
//	GET ORDER BY ID
//
// ───────────────────────────────────────────

// GetOrderByID obtiene una orden por su ID
func (uc *UseCases) GetOrderByID(ctx context.Context, id string) (*domain.OrderResponse, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	order, err := uc.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	return mapOrderToResponse(order), nil
}

// ───────────────────────────────────────────
//
//	LIST ORDERS
//
// ───────────────────────────────────────────

// ListOrders obtiene una lista paginada de órdenes con filtros
func (uc *UseCases) ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*domain.OrdersListResponse, error) {
	// Validar paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Obtener órdenes del repositorio
	orders, total, err := uc.repo.ListOrders(ctx, page, pageSize, filters)
	if err != nil {
		return nil, fmt.Errorf("error listing orders: %w", err)
	}

	// Mapear a respuestas
	orderResponses := make([]domain.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = *mapOrderToResponse(&order)
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &domain.OrdersListResponse{
		Data:       orderResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// ───────────────────────────────────────────
//
//	UPDATE ORDER
//
// ───────────────────────────────────────────

// UpdateOrder actualiza una orden existente
func (uc *UseCases) UpdateOrder(ctx context.Context, id string, req *domain.UpdateOrderRequest) (*domain.OrderResponse, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	// Obtener la orden existente
	order, err := uc.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	// Actualizar solo los campos proporcionados
	if req.Subtotal != nil {
		order.Subtotal = *req.Subtotal
	}
	if req.Tax != nil {
		order.Tax = *req.Tax
	}
	if req.Discount != nil {
		order.Discount = *req.Discount
	}
	if req.ShippingCost != nil {
		order.ShippingCost = *req.ShippingCost
	}
	if req.TotalAmount != nil {
		order.TotalAmount = *req.TotalAmount
	}
	if req.Currency != nil {
		order.Currency = *req.Currency
	}
	if req.CodTotal != nil {
		order.CodTotal = req.CodTotal
	}

	// Información del cliente
	if req.CustomerName != nil {
		order.CustomerName = *req.CustomerName
	}
	if req.CustomerEmail != nil {
		order.CustomerEmail = *req.CustomerEmail
	}
	if req.CustomerPhone != nil {
		order.CustomerPhone = *req.CustomerPhone
	}
	if req.CustomerDNI != nil {
		order.CustomerDNI = *req.CustomerDNI
	}

	// Dirección de envío
	if req.ShippingStreet != nil {
		order.ShippingStreet = *req.ShippingStreet
	}
	if req.ShippingCity != nil {
		order.ShippingCity = *req.ShippingCity
	}
	if req.ShippingState != nil {
		order.ShippingState = *req.ShippingState
	}
	if req.ShippingCountry != nil {
		order.ShippingCountry = *req.ShippingCountry
	}
	if req.ShippingPostalCode != nil {
		order.ShippingPostalCode = *req.ShippingPostalCode
	}
	if req.ShippingLat != nil {
		order.ShippingLat = req.ShippingLat
	}
	if req.ShippingLng != nil {
		order.ShippingLng = req.ShippingLng
	}

	// Información de pago
	if req.PaymentMethodID != nil {
		order.PaymentMethodID = *req.PaymentMethodID
	}
	if req.IsPaid != nil {
		order.IsPaid = *req.IsPaid
	}
	if req.PaidAt != nil {
		order.PaidAt = req.PaidAt
	}

	// Información de envío/logística
	if req.TrackingNumber != nil {
		order.TrackingNumber = req.TrackingNumber
	}
	if req.TrackingLink != nil {
		order.TrackingLink = req.TrackingLink
	}
	if req.GuideID != nil {
		order.GuideID = req.GuideID
	}
	if req.GuideLink != nil {
		order.GuideLink = req.GuideLink
	}
	if req.DeliveryDate != nil {
		order.DeliveryDate = req.DeliveryDate
	}
	if req.DeliveredAt != nil {
		order.DeliveredAt = req.DeliveredAt
	}

	// Información de fulfillment
	if req.WarehouseID != nil {
		order.WarehouseID = req.WarehouseID
	}
	if req.WarehouseName != nil {
		order.WarehouseName = *req.WarehouseName
	}
	if req.DriverID != nil {
		order.DriverID = req.DriverID
	}
	if req.DriverName != nil {
		order.DriverName = *req.DriverName
	}
	if req.IsLastMile != nil {
		order.IsLastMile = *req.IsLastMile
	}

	// Dimensiones y peso
	if req.Weight != nil {
		order.Weight = req.Weight
	}
	if req.Height != nil {
		order.Height = req.Height
	}
	if req.Width != nil {
		order.Width = req.Width
	}
	if req.Length != nil {
		order.Length = req.Length
	}
	if req.Boxes != nil {
		order.Boxes = req.Boxes
	}

	// Tipo y estado
	if req.OrderTypeID != nil {
		order.OrderTypeID = req.OrderTypeID
	}
	if req.OrderTypeName != nil {
		order.OrderTypeName = *req.OrderTypeName
	}
	if req.Status != nil {
		order.Status = *req.Status
	}
	if req.OriginalStatus != nil {
		order.OriginalStatus = *req.OriginalStatus
	}

	// Información adicional
	if req.Notes != nil {
		order.Notes = req.Notes
	}
	if req.Coupon != nil {
		order.Coupon = req.Coupon
	}
	if req.Approved != nil {
		order.Approved = req.Approved
	}
	if req.UserID != nil {
		order.UserID = req.UserID
	}
	if req.UserName != nil {
		order.UserName = *req.UserName
	}

	// Facturación
	if req.Invoiceable != nil {
		order.Invoiceable = *req.Invoiceable
	}
	if req.InvoiceURL != nil {
		order.InvoiceURL = req.InvoiceURL
	}
	if req.InvoiceID != nil {
		order.InvoiceID = req.InvoiceID
	}
	if req.InvoiceProvider != nil {
		order.InvoiceProvider = req.InvoiceProvider
	}

	// Datos estructurados
	if req.Items != nil {
		order.Items = req.Items
	}
	if req.Metadata != nil {
		order.Metadata = req.Metadata
	}
	if req.FinancialDetails != nil {
		order.FinancialDetails = req.FinancialDetails
	}
	if req.ShippingDetails != nil {
		order.ShippingDetails = req.ShippingDetails
	}
	if req.PaymentDetails != nil {
		order.PaymentDetails = req.PaymentDetails
	}
	if req.FulfillmentDetails != nil {
		order.FulfillmentDetails = req.FulfillmentDetails
	}

	// Guardar cambios
	if err := uc.repo.UpdateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	return mapOrderToResponse(order), nil
}

// ───────────────────────────────────────────
//
//	DELETE ORDER
//
// ───────────────────────────────────────────

// DeleteOrder elimina (soft delete) una orden
func (uc *UseCases) DeleteOrder(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("order ID is required")
	}

	// Verificar que la orden existe
	_, err := uc.repo.GetOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting order: %w", err)
	}

	// Eliminar la orden
	if err := uc.repo.DeleteOrder(ctx, id); err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}

	return nil
}

// ───────────────────────────────────────────
//
//	HELPER FUNCTIONS
//
// ───────────────────────────────────────────

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

		// Dirección de envío
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
