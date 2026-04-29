package repository

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
)

func (r *Repository) ProfitReportDetail(ctx context.Context, params dtos.ProfitReportDetailParams) (*dtos.ProfitReportDetailResponse, error) {
	type row struct {
		ShipmentID     uint
		OrderNumber    string
		TrackingNumber string
		Carrier        string
		CustomerCharge float64
		CarrierCost    float64
		Status         string
		CreatedAt      time.Time
	}

	conn := r.db.Conn(ctx)

	base := conn.Table("shipments AS s").
		Joins("JOIN orders o ON o.id = s.order_id").
		Where("s.deleted_at IS NULL").
		Where("s.tracking_number IS NOT NULL AND s.tracking_number <> ''").
		Where("s.total_cost IS NOT NULL").
		Where("s.status <> ?", "cancelled")

	if params.BusinessID > 0 {
		base = base.Where("o.business_id = ?", params.BusinessID)
	}
	if params.From != nil {
		base = base.Where("s.created_at >= ?", *params.From)
	}
	if params.To != nil {
		base = base.Where("s.created_at < ?", *params.To)
	}
	if params.Carrier != "" {
		base = base.Where("LOWER(COALESCE(s.carrier, '')) = ?", strings.ToLower(params.Carrier))
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}

	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var rows []row
	err := base.
		Select(`s.id AS shipment_id,
COALESCE(o.order_number, '') AS order_number,
COALESCE(s.tracking_number, '') AS tracking_number,
COALESCE(NULLIF(s.carrier, ''), 'Sin transportadora') AS carrier,
COALESCE(s.total_cost, 0) AS customer_charge,
COALESCE(s.carrier_cost, 0) AS carrier_cost,
COALESCE(s.status, '') AS status,
s.created_at AS created_at`).
		Order("s.created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	data := make([]dtos.ProfitReportDetailRow, 0, len(rows))
	for _, r := range rows {
		data = append(data, dtos.ProfitReportDetailRow{
			ShipmentID:     r.ShipmentID,
			OrderNumber:    r.OrderNumber,
			TrackingNumber: r.TrackingNumber,
			Carrier:        r.Carrier,
			CustomerCharge: r.CustomerCharge,
			CarrierCost:    r.CarrierCost,
			Profit:         r.CustomerCharge - r.CarrierCost,
			Status:         r.Status,
			CreatedAt:      r.CreatedAt,
		})
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return &dtos.ProfitReportDetailResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
