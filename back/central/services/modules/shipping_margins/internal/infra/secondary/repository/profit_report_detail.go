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
		ServiceType    string
		CustomerCharge float64
		CarrierCost    float64
		Status         string
		CreatedAt      time.Time
	}

	conn := r.db.Conn(ctx)

	where := []string{
		"s.deleted_at IS NULL",
		"s.tracking_number IS NOT NULL AND s.tracking_number <> ''",
		"s.total_cost IS NOT NULL",
		"s.status <> 'cancelled'",
	}
	args := []interface{}{}

	if params.BusinessID > 0 {
		where = append(where, "o.business_id = ?")
		args = append(args, params.BusinessID)
	}
	if params.From != nil {
		where = append(where, "s.created_at >= ?")
		args = append(args, *params.From)
	}
	if params.To != nil {
		where = append(where, "s.created_at < ?")
		args = append(args, *params.To)
	}
	if params.Carrier != "" {
		where = append(where, "LOWER(COALESCE(s.carrier, '')) = ?")
		args = append(args, strings.ToLower(params.Carrier))
	}

	whereSQL := strings.Join(where, " AND ")

	guideCharge := "GREATEST(COALESCE(s.total_cost, 0) - COALESCE(s.cod_probability_margin, 0), 0)"

	baseSQL := `
SELECT shipment_id, order_number, tracking_number, carrier, service_type,
       customer_charge, carrier_cost, status, created_at
FROM (
    SELECT s.id AS shipment_id,
           COALESCE(o.order_number, '') AS order_number,
           COALESCE(s.tracking_number, '') AS tracking_number,
           COALESCE(NULLIF(s.carrier, ''), 'Sin transportadora') AS carrier,
           'guide' AS service_type,
           ` + guideCharge + ` AS customer_charge,
           COALESCE(s.carrier_cost, 0) AS carrier_cost,
           COALESCE(s.status, '') AS status,
           s.created_at AS created_at,
           1 AS service_order
    FROM shipments s
    JOIN orders o ON o.id = s.order_id
    WHERE ` + whereSQL + `

    UNION ALL

    SELECT s.id AS shipment_id,
           COALESCE(o.order_number, '') AS order_number,
           COALESCE(s.tracking_number, '') AS tracking_number,
           COALESCE(NULLIF(s.carrier, ''), 'Sin transportadora') AS carrier,
           'cod' AS service_type,
           COALESCE(s.cod_probability_margin, 0) AS customer_charge,
           0 AS carrier_cost,
           COALESCE(s.status, '') AS status,
           s.created_at AS created_at,
           2 AS service_order
    FROM shipments s
    JOIN orders o ON o.id = s.order_id
    WHERE ` + whereSQL + ` AND COALESCE(s.cod_probability_margin, 0) > 0
) AS combined
ORDER BY created_at DESC, shipment_id DESC, service_order ASC
`

	countSQL := `
SELECT COUNT(*) FROM (
    SELECT 1
    FROM shipments s
    JOIN orders o ON o.id = s.order_id
    WHERE ` + whereSQL + `
    UNION ALL
    SELECT 1
    FROM shipments s
    JOIN orders o ON o.id = s.order_id
    WHERE ` + whereSQL + ` AND COALESCE(s.cod_probability_margin, 0) > 0
) AS combined
`

	countArgs := append([]interface{}{}, args...)
	countArgs = append(countArgs, args...)

	var total int64
	if err := conn.Raw(countSQL, countArgs...).Scan(&total).Error; err != nil {
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

	queryArgs := append([]interface{}{}, args...)
	queryArgs = append(queryArgs, args...)
	queryArgs = append(queryArgs, pageSize, offset)

	var rows []row
	err := conn.Raw(baseSQL+" LIMIT ? OFFSET ?", queryArgs...).Scan(&rows).Error
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
			ServiceType:    r.ServiceType,
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
