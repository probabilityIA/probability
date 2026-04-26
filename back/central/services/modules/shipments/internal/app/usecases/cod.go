package usecases

import (
	"context"
	"math"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (uc *UseCases) ListCODShipments(ctx context.Context, filter domain.CODFilter) (*domain.ShipmentsListResponse, error) {
	shipments, total, err := uc.repo.ListCODShipments(ctx, filter)
	if err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	out := make([]domain.ShipmentResponse, len(shipments))
	for i := range shipments {
		out[i] = mapShipmentToCODResponse(&shipments[i])
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	return &domain.ShipmentsListResponse{
		Data:       out,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (uc *UseCases) CollectCOD(ctx context.Context, shipmentID uint, notes string) (*domain.ShipmentResponse, error) {
	shipment, err := uc.repo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, err
	}
	if shipment.Status != "delivered" {
		return nil, domain.ErrShipmentNotDelivered
	}
	if shipment.OrderID == nil || *shipment.OrderID == "" {
		return nil, domain.ErrOrderIDRequired
	}

	info, err := uc.repo.GetOrderCODInfo(ctx, *shipment.OrderID)
	if err != nil {
		return nil, err
	}
	if info.CodTotal == nil || *info.CodTotal <= 0 {
		return nil, domain.ErrOrderNotCOD
	}
	if info.IsPaid {
		return nil, domain.ErrOrderAlreadyPaid
	}

	if err := uc.repo.MarkOrderPaidCOD(ctx, info.OrderID, *info.CodTotal, info.PaymentMethodID, notes); err != nil {
		return nil, err
	}

	updated, err := uc.repo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, err
	}
	resp := mapShipmentToCODResponse(updated)
	return &resp, nil
}

func mapShipmentToCODResponse(s *domain.Shipment) domain.ShipmentResponse {
	return domain.ShipmentResponse{
		ID:                 s.ID,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
		DeletedAt:          s.DeletedAt,
		OrderID:            s.OrderID,
		ClientName:         s.ClientName,
		DestinationAddress: s.DestinationAddress,
		TrackingNumber:     s.TrackingNumber,
		TrackingURL:        s.TrackingURL,
		Carrier:            s.Carrier,
		CarrierCode:        s.CarrierCode,
		GuideID:            s.GuideID,
		GuideURL:           s.GuideURL,
		Status:             s.Status,
		ShippedAt:          s.ShippedAt,
		DeliveredAt:        s.DeliveredAt,
		ShippingAddressID:  s.ShippingAddressID,
		ShippingCost:       s.ShippingCost,
		InsuranceCost:      s.InsuranceCost,
		TotalCost:          s.TotalCost,
		Weight:             s.Weight,
		Height:             s.Height,
		Width:              s.Width,
		Length:             s.Length,
		WarehouseID:        s.WarehouseID,
		WarehouseName:      s.WarehouseName,
		DriverID:           s.DriverID,
		DriverName:         s.DriverName,
		IsLastMile:         s.IsLastMile,
		EstimatedDelivery:  s.EstimatedDelivery,
		DeliveryNotes:      s.DeliveryNotes,
		Metadata:           s.Metadata,
		CustomerName:       s.CustomerName,
		CustomerEmail:      s.CustomerEmail,
		CustomerPhone:      s.CustomerPhone,
		CustomerDNI:        s.CustomerDNI,
		OrderNumber:        s.OrderNumber,
		CodTotal:           s.CodTotal,
		IsPaid:             s.IsPaid,
		PaidAt:             s.PaidAt,
		PaymentMethodCode:  s.PaymentMethodCode,
		OrderTotalAmount:   s.OrderTotalAmount,
		OrderCurrency:      s.OrderCurrency,
	}
}
