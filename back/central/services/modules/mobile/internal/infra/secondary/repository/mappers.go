package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func toOrderSummary(order models.Order) entities.OrderSummary {
	return entities.OrderSummary{
		ID:              order.ID,
		OrderNumber:     order.OrderNumber,
		InternalNumber:  order.InternalNumber,
		IntegrationID:   order.IntegrationID,
		IntegrationType: order.IntegrationType,
		Platform:        order.Platform,
		Status:          order.Status,
		StatusID:        order.StatusID,
		IsPaid:          order.IsPaid,
		IsCod:           order.IsCod,
		CodTotal:        order.CodTotal,
		Subtotal:        order.Subtotal,
		Tax:             order.Tax,
		Discount:        order.Discount,
		ShippingCost:    order.ShippingCost,
		TotalAmount:     order.TotalAmount,
		Currency:        order.Currency,
		CustomerName:    order.CustomerName,
		CustomerEmail:   order.CustomerEmail,
		CustomerPhone:   order.CustomerPhone,
		ShippingStreet:  order.ShippingStreet,
		ShippingCity:    order.ShippingCity,
		WarehouseName:   order.WarehouseName,
		UserName:        order.UserName,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func toOrderLines(items []models.OrderItem) []entities.OrderLine {
	lines := make([]entities.OrderLine, 0, len(items))
	for _, item := range items {
		lines = append(lines, entities.OrderLine{
			ID:         item.ID,
			ProductID:  item.ProductID,
			SKU:        item.ProductSKU,
			Name:       item.ProductName,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
		})
	}
	return lines
}

func toShipmentSummary(shipment models.Shipment) *entities.ShipmentSummary {
	return &entities.ShipmentSummary{
		ID:                   shipment.ID,
		TrackingNumber:       shipment.TrackingNumber,
		TrackingURL:          shipment.TrackingURL,
		Carrier:              shipment.Carrier,
		GuideURL:             shipment.GuideURL,
		Status:               shipment.Status,
		CarrierStatus:        shipment.CarrierStatus,
		DestinationCity:      shipment.DestinationCity,
		InsuranceCost:        shipment.InsuranceCost,
		TotalCost:            shipment.TotalCost,
		CarrierCost:          shipment.CarrierCost,
		AppliedMargin:        shipment.AppliedMargin,
		CodCarrierFee:        shipment.CodCarrierFee,
		CodProbabilityMargin: shipment.CodProbabilityMargin,
		CreatedAt:            shipment.CreatedAt,
	}
}

func toInvoiceSummary(invoice models.Invoice) *entities.InvoiceSummary {
	return &entities.InvoiceSummary{
		ID:            invoice.ID,
		InvoiceNumber: invoice.InvoiceNumber,
		Status:        invoice.Status,
		TotalAmount:   invoice.TotalAmount,
		Currency:      invoice.Currency,
		InvoiceURL:    invoice.InvoiceURL,
		CUFE:          invoice.CUFE,
		IssuedAt:      invoice.IssuedAt,
		CreatedAt:     invoice.CreatedAt,
	}
}
