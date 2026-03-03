package usecaseupdateorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// updateAdditionalFields actualiza campos adicionales de la orden
func (uc *UseCaseUpdateOrder) updateAdditionalFields(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	changed := false

	if dto.OrderTypeID != nil && (order.OrderTypeID == nil || *order.OrderTypeID != *dto.OrderTypeID) {
		order.OrderTypeID = dto.OrderTypeID
		changed = true
	}

	if dto.OrderTypeName != "" && order.OrderTypeName != dto.OrderTypeName {
		order.OrderTypeName = dto.OrderTypeName
		changed = true
	}

	if dto.Notes != nil && (order.Notes == nil || *order.Notes != *dto.Notes) {
		order.Notes = dto.Notes
		changed = true
	}

	if dto.Coupon != nil && (order.Coupon == nil || *order.Coupon != *dto.Coupon) {
		order.Coupon = dto.Coupon
		changed = true
	}

	if dto.Approved != nil && (order.Approved == nil || *order.Approved != *dto.Approved) {
		order.Approved = dto.Approved
		changed = true
	}

	if dto.UserID != nil && (order.UserID == nil || *order.UserID != *dto.UserID) {
		order.UserID = dto.UserID
		changed = true
	}

	if dto.UserName != "" && order.UserName != dto.UserName {
		order.UserName = dto.UserName
		changed = true
	}

	// Actualizar campos de facturación
	if order.Invoiceable != dto.Invoiceable {
		order.Invoiceable = dto.Invoiceable
		changed = true
	}

	if dto.InvoiceURL != nil && (order.InvoiceURL == nil || *order.InvoiceURL != *dto.InvoiceURL) {
		order.InvoiceURL = dto.InvoiceURL
		changed = true
	}

	if dto.InvoiceID != nil && (order.InvoiceID == nil || *order.InvoiceID != *dto.InvoiceID) {
		order.InvoiceID = dto.InvoiceID
		changed = true
	}

	if dto.InvoiceProvider != nil && (order.InvoiceProvider == nil || *order.InvoiceProvider != *dto.InvoiceProvider) {
		order.InvoiceProvider = dto.InvoiceProvider
		changed = true
	}

	if dto.OrderStatusURL != "" && order.OrderStatusURL != dto.OrderStatusURL {
		order.OrderStatusURL = dto.OrderStatusURL
		changed = true
	}

	return changed
}
