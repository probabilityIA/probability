package usecasecreateorder

import (
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// buildOrderEntity construye la entidad ProbabilityOrder desde el DTO
func (uc *UseCaseCreateOrder) buildOrderEntity(dto *dtos.ProbabilityOrderDTO, clientID *uint, statusMapping orderStatusMapping) *entities.ProbabilityOrder {
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
		ShippingCostPresentment:     dto.ShippingCostPresentment,
		ShippingDiscount:            dto.ShippingDiscount,
		ShippingDiscountPresentment: dto.ShippingDiscountPresentment,
		TotalAmountPresentment:  dto.TotalAmountPresentment,
		CurrencyPresentment:     dto.CurrencyPresentment,

		// Información del cliente
		CustomerID:        clientID,
		CustomerName:      dto.CustomerName,
		CustomerFirstName: dto.CustomerFirstName,
		CustomerLastName:  dto.CustomerLastName,
		CustomerEmail:     dto.CustomerEmail,
		CustomerPhone:     dto.CustomerPhone,
		CustomerDNI:       dto.CustomerDNI,
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

		// Testing
		IsTest: dto.IsTest,

		// Facturación
		Invoiceable:     dto.Invoiceable,
		InvoiceURL:      dto.InvoiceURL,
		InvoiceID:       dto.InvoiceID,
		InvoiceProvider: dto.InvoiceProvider,
		OrderStatusURL:  dto.OrderStatusURL,

		// Datos estructurados (JSONB)
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
func (uc *UseCaseCreateOrder) assignPaymentMethodID(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) {
	order.PaymentMethodID = 1 // Valor por defecto

	if len(dto.Payments) == 0 {
		return
	}

	payment := dto.Payments[0]
	if payment.PaymentMethodID > 0 {
		order.PaymentMethodID = payment.PaymentMethodID
	}
	if payment.Status == "completed" {
		order.IsPaid = true
		if payment.PaidAt != nil {
			order.PaidAt = payment.PaidAt
		} else {
			now := time.Now()
			order.PaidAt = &now
		}
	}
}

// populateOrderFields popula campos planos de dirección desde Addresses
func (uc *UseCaseCreateOrder) populateOrderFields(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) {
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
