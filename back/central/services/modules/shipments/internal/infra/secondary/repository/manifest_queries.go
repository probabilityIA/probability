package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (r *Repository) ListPendingForManifest(ctx context.Context, filter domain.ManifestFilter) ([]domain.ManifestShipmentRow, error) {
	if filter.BusinessID == 0 {
		return []domain.ManifestShipmentRow{}, nil
	}

	businessIDs := []uint{filter.BusinessID}
	if filter.IncludeChildren {
		var childIDs []uint
		if err := r.db.Conn(ctx).
			Table("business").
			Select("id").
			Where("parent_business_id = ? AND deleted_at IS NULL", filter.BusinessID).
			Scan(&childIDs).Error; err == nil {
			businessIDs = append(businessIDs, childIDs...)
		}
	}

	type row struct {
		ShipmentID         uint
		OrderID            *string
		OrderNumber        string
		TrackingNumber     *string
		Carrier            *string
		CarrierCode        *string
		CustomerName       string
		CustomerDNI        string
		ShippingStreet     string
		ShippingCity       string
		ShippingState      string
		Weight             *float64
		TotalAmount        float64
		CodTotal           *float64
		BusinessID         *uint
		BusinessName       string
		WarehouseName      *string
		ShipmentCreatedAt  *time.Time
		OrderCreatedAt     *time.Time
		ShipmentStatus     string
		OrderStatus        string
	}

	q := r.db.Conn(ctx).
		Table("shipments AS s").
		Select(`s.id AS shipment_id,
			s.order_id,
			COALESCE(o.order_number, '') AS order_number,
			s.tracking_number,
			s.carrier,
			s.carrier_code,
			COALESCE(o.customer_name, '') AS customer_name,
			COALESCE(o.customer_dni, '') AS customer_dni,
			COALESCE(o.shipping_street, '') AS shipping_street,
			COALESCE(o.shipping_city, '') AS shipping_city,
			COALESCE(o.shipping_state, '') AS shipping_state,
			s.weight,
			COALESCE(o.total_amount, 0) AS total_amount,
			o.cod_total,
			o.business_id,
			COALESCE(b.name, '') AS business_name,
			w.name AS warehouse_name,
			s.created_at AS shipment_created_at,
			o.created_at AS order_created_at,
			COALESCE(s.status, '') AS shipment_status,
			COALESCE(o.status, '') AS order_status`).
		Joins("LEFT JOIN orders o ON o.id = s.order_id").
		Joins("LEFT JOIN business b ON b.id = o.business_id").
		Joins("LEFT JOIN warehouses w ON w.id = s.warehouse_id").
		Where("s.deleted_at IS NULL").
		Where("s.status = ?", "pending").
		Where("o.business_id IN ?", businessIDs)

	if filter.Carrier != "" {
		q = q.Where("s.carrier ILIKE ?", "%"+filter.Carrier+"%")
	}

	q = q.Order("s.carrier ASC, s.created_at DESC")

	var rows []row
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]domain.ManifestShipmentRow, 0, len(rows))
	for _, r := range rows {
		item := domain.ManifestShipmentRow{
			ShipmentID:         r.ShipmentID,
			OrderID:            r.OrderID,
			OrderNumber:        r.OrderNumber,
			CustomerName:       r.CustomerName,
			CustomerDocument:   r.CustomerDNI,
			DestinationAddress: r.ShippingStreet,
			DestinationCity:    r.ShippingCity,
			DestinationState:   r.ShippingState,
			BusinessName:       r.BusinessName,
		}
		if r.TrackingNumber != nil {
			item.TrackingNumber = *r.TrackingNumber
		}
		if r.Carrier != nil {
			item.Carrier = *r.Carrier
		}
		if r.CarrierCode != nil {
			item.CarrierCode = *r.CarrierCode
		}
		if r.Weight != nil {
			item.Weight = *r.Weight
		}
		if r.CodTotal != nil {
			item.CodTotal = *r.CodTotal
		}
		item.DeclaredValue = r.TotalAmount
		if r.BusinessID != nil {
			item.BusinessID = *r.BusinessID
		}
		if r.WarehouseName != nil {
			item.WarehouseName = *r.WarehouseName
		}
		item.ShipmentCreatedAt = r.ShipmentCreatedAt
		item.OrderCreatedAt = r.OrderCreatedAt
		item.ShipmentStatus = r.ShipmentStatus
		item.OrderStatus = r.OrderStatus
		out = append(out, item)
	}
	return out, nil
}

func (r *Repository) GetBusinessForManifest(ctx context.Context, businessID uint) (*domain.ManifestBusinessInfo, error) {
	type row struct {
		ID               uint
		Name             string
		Code             string
		Address          string
		ParentBusinessID *uint
	}
	var b row
	err := r.db.Conn(ctx).
		Table("business").
		Select("id, name, code, address, parent_business_id").
		Where("id = ? AND deleted_at IS NULL", businessID).
		Scan(&b).Error
	if err != nil {
		return nil, err
	}
	if b.ID == 0 {
		return nil, nil
	}
	return &domain.ManifestBusinessInfo{
		ID:       b.ID,
		Name:     b.Name,
		Code:     b.Code,
		Address:  b.Address,
		ParentID: b.ParentBusinessID,
	}, nil
}

func (r *Repository) GetChildBusinessIDs(ctx context.Context, parentID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Conn(ctx).
		Table("business").
		Select("id").
		Where("parent_business_id = ? AND deleted_at IS NULL", parentID).
		Scan(&ids).Error
	return ids, err
}
