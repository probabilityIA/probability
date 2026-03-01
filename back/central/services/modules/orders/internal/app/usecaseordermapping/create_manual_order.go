package usecaseordermapping

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
)

// CreateManualOrder crea una orden manual pasando por el pipeline completo de MapAndSaveOrder.
// Aplica defaults, convierte CreateOrderRequest → ProbabilityOrderDTO y delega a MapAndSaveOrder.
func (uc *UseCaseOrderMapping) CreateManualOrder(ctx context.Context, req *dtos.CreateOrderRequest) (*dtos.OrderResponse, error) {
	if req.BusinessID == nil || *req.BusinessID == 0 {
		return nil, fmt.Errorf("business_id is required")
	}

	// 1. Aplicar defaults para órdenes manuales
	if err := uc.applyManualDefaults(ctx, req); err != nil {
		return nil, fmt.Errorf("error applying manual defaults: %w", err)
	}

	// 2. Convertir CreateOrderRequest → ProbabilityOrderDTO
	dto := uc.mapCreateRequestToDTO(req)

	// 3. Delegar al pipeline completo
	return uc.MapAndSaveOrder(ctx, dto)
}

// applyManualDefaults aplica valores por defecto para órdenes manuales:
// IntegrationID, ExternalID, OrderNumber, PaymentMethodID, Platform, IntegrationType.
func (uc *UseCaseOrderMapping) applyManualDefaults(ctx context.Context, req *dtos.CreateOrderRequest) error {
	// Auto IntegrationID
	if req.IntegrationID == 0 {
		intID, err := uc.repo.GetPlatformIntegrationIDByBusinessID(ctx, *req.BusinessID)
		if err != nil {
			intID, err = uc.repo.GetFirstIntegrationIDByBusinessID(ctx, *req.BusinessID)
			if err != nil {
				uc.logger.Warn().Err(err).Msg("No default integration found for manual order, proceeding with ID 0")
				req.IntegrationID = 0
			} else {
				req.IntegrationID = intID
			}
		} else {
			req.IntegrationID = intID
		}
	}

	// Auto ExternalID
	if req.ExternalID == "" {
		platform := req.Platform
		if platform == "" {
			platform = "manual"
		}
		req.ExternalID = fmt.Sprintf("%s-%d", platform, time.Now().UnixNano())
	}

	// Auto OrderNumber
	if req.OrderNumber == "" || req.OrderNumber == "AUTO" {
		lastNum, err := uc.repo.GetLastManualOrderNumber(ctx, *req.BusinessID)
		if err != nil {
			return fmt.Errorf("error getting last manual order number: %w", err)
		}
		req.OrderNumber = fmt.Sprintf("prob-%04d", lastNum+1)
	}

	// Defaults
	if req.PaymentMethodID == 0 {
		req.PaymentMethodID = 1
	}
	if req.Platform == "" {
		req.Platform = "manual"
	}
	if req.IntegrationType == "" {
		req.IntegrationType = "platform"
	}

	return nil
}

// mapCreateRequestToDTO convierte un CreateOrderRequest plano en un ProbabilityOrderDTO
// con arrays de entidades relacionadas (Addresses, Payments, Shipments).
func (uc *UseCaseOrderMapping) mapCreateRequestToDTO(req *dtos.CreateOrderRequest) *dtos.ProbabilityOrderDTO {
	now := time.Now()

	dto := &dtos.ProbabilityOrderDTO{
		IsManualOrder: true,

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
		CustomerID:        req.CustomerID,
		CustomerName:      req.CustomerName,
		CustomerFirstName: req.CustomerFirstName,
		CustomerLastName:  req.CustomerLastName,
		CustomerEmail:     req.CustomerEmail,
		CustomerPhone:     req.CustomerPhone,
		CustomerDNI:       req.CustomerDNI,

		// Tipo y estado
		OrderTypeID:    req.OrderTypeID,
		OrderTypeName:  req.OrderTypeName,
		Status:         req.Status,
		OriginalStatus: req.OriginalStatus,
		StatusID:       req.StatusID,

		// Estados independientes
		PaymentStatusID:     req.PaymentStatusID,
		FulfillmentStatusID: req.FulfillmentStatusID,

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

		// Datos estructurados (JSONB)
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

	// Defaults de timestamps
	if dto.OccurredAt.IsZero() {
		dto.OccurredAt = now
	}
	if dto.ImportedAt.IsZero() {
		dto.ImportedAt = now
	}

	// Shipping fields → Address
	if req.ShippingStreet != "" || req.ShippingCity != "" {
		dto.Addresses = append(dto.Addresses, dtos.ProbabilityAddressDTO{
			Type:       "shipping",
			FirstName:  req.CustomerFirstName,
			LastName:   req.CustomerLastName,
			Phone:      req.CustomerPhone,
			Street:     req.ShippingStreet,
			City:       req.ShippingCity,
			State:      req.ShippingState,
			Country:    req.ShippingCountry,
			PostalCode: req.ShippingPostalCode,
			Latitude:   req.ShippingLat,
			Longitude:  req.ShippingLng,
		})
	}

	// Payment fields → Payment
	if req.PaymentMethodID > 0 || req.IsPaid {
		paymentStatus := "pending"
		if req.IsPaid {
			paymentStatus = "completed"
		}
		dto.Payments = append(dto.Payments, dtos.ProbabilityPaymentDTO{
			PaymentMethodID: req.PaymentMethodID,
			Amount:          req.TotalAmount,
			Currency:        req.Currency,
			Status:          paymentStatus,
			PaidAt:          req.PaidAt,
		})
	}

	// Logistics fields → Shipment
	if req.TrackingNumber != nil || req.GuideID != nil || req.WarehouseID != nil {
		dto.Shipments = append(dto.Shipments, dtos.ProbabilityShipmentDTO{
			TrackingNumber:    req.TrackingNumber,
			TrackingURL:       req.TrackingLink,
			GuideID:           req.GuideID,
			GuideURL:          req.GuideLink,
			Status:            "pending",
			ShippedAt:         req.DeliveryDate,
			DeliveredAt:       req.DeliveredAt,
			ShippingCost:      &req.ShippingCost,
			Weight:            req.Weight,
			Height:            req.Height,
			Width:             req.Width,
			Length:            req.Length,
			WarehouseID:       req.WarehouseID,
			WarehouseName:     req.WarehouseName,
			DriverID:          req.DriverID,
			DriverName:        req.DriverName,
			IsLastMile:        req.IsLastMile,
			EstimatedDelivery: req.DeliveryDate,
		})
	}

	return dto
}
