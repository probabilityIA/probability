package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (r *Repository) resolveManifestBusinessIDs(ctx context.Context, businessID uint, includeChildren bool) []uint {
	ids := []uint{businessID}
	if includeChildren {
		var childIDs []uint
		if err := r.db.Conn(ctx).
			Table("business").
			Select("id").
			Where("parent_business_id = ? AND deleted_at IS NULL", businessID).
			Scan(&childIDs).Error; err == nil {
			ids = append(ids, childIDs...)
		}
	}
	return ids
}

func (r *Repository) ListPendingCarriers(ctx context.Context, businessID uint, includeChildren bool) ([]domain.ManifestCarrierCount, error) {
	if businessID == 0 {
		return []domain.ManifestCarrierCount{}, nil
	}
	businessIDs := r.resolveManifestBusinessIDs(ctx, businessID, includeChildren)

	type row struct {
		Carrier *string
		Count   int64
	}
	var rows []row
	err := r.db.Conn(ctx).
		Table("shipments AS s").
		Select("s.carrier, COUNT(*) AS count").
		Joins("LEFT JOIN orders o ON o.id = s.order_id").
		Where("s.deleted_at IS NULL").
		Where("s.status = ?", "pending").
		Where("o.business_id IN ?", businessIDs).
		Group("s.carrier").
		Order("count DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.ManifestCarrierCount, 0, len(rows))
	for _, r := range rows {
		name := ""
		if r.Carrier != nil {
			name = *r.Carrier
		}
		if name == "" {
			name = "Sin asignar"
		}
		out = append(out, domain.ManifestCarrierCount{Carrier: name, Count: r.Count})
	}
	return out, nil
}

func (r *Repository) ListPendingForManifest(ctx context.Context, filter domain.ManifestFilter) ([]domain.ManifestShipmentRow, int64, error) {
	if filter.BusinessID == 0 {
		return []domain.ManifestShipmentRow{}, 0, nil
	}

	businessIDs := r.resolveManifestBusinessIDs(ctx, filter.BusinessID, filter.IncludeChildren)

	base := r.db.Conn(ctx).
		Table("shipments AS s").
		Joins("LEFT JOIN orders o ON o.id = s.order_id").
		Where("s.deleted_at IS NULL").
		Where("s.status = ?", "pending").
		Where("o.business_id IN ?", businessIDs)

	if filter.Carrier != "" {
		if filter.Carrier == "Sin asignar" {
			base = base.Where("(s.carrier IS NULL OR s.carrier = '')")
		} else {
			base = base.Where("s.carrier = ?", filter.Carrier)
		}
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
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

	q := base.
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
		Joins("LEFT JOIN business b ON b.id = o.business_id").
		Joins("LEFT JOIN warehouses w ON w.id = s.warehouse_id").
		Order("s.created_at DESC")

	if filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		if offset < 0 {
			offset = 0
		}
		q = q.Offset(offset).Limit(filter.PageSize)
	}

	var rows []row
	if err := q.Scan(&rows).Error; err != nil {
		return nil, 0, err
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
			ShipmentCreatedAt:  r.ShipmentCreatedAt,
			OrderCreatedAt:     r.OrderCreatedAt,
			ShipmentStatus:     r.ShipmentStatus,
			OrderStatus:        r.OrderStatus,
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
		out = append(out, item)
	}
	return out, total, nil
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
