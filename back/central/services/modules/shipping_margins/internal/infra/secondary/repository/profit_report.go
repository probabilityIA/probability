package repository

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
)

func (r *Repository) ProfitReport(ctx context.Context, params dtos.ProfitReportParams) (*dtos.ProfitReportResponse, error) {
	type row struct {
		Carrier             string
		Shipments           int
		CarrierCostTotal    float64
		CustomerChargeTotal float64
		ProfitTotal         float64
	}

	conn := r.db.Conn(ctx)

	q := conn.Table("shipments AS s").
		Select(`COALESCE(NULLIF(s.carrier, ''), 'Sin transportadora') AS carrier,
COUNT(*) AS shipments,
COALESCE(SUM(s.carrier_cost), 0) AS carrier_cost_total,
COALESCE(SUM(s.total_cost), 0) AS customer_charge_total,
COALESCE(SUM(s.total_cost - s.carrier_cost), 0) AS profit_total`).
		Joins("JOIN orders o ON o.id = s.order_id").
		Where("s.deleted_at IS NULL").
		Where("s.tracking_number IS NOT NULL AND s.tracking_number <> ''").
		Where("s.total_cost IS NOT NULL").
		Group("COALESCE(NULLIF(s.carrier, ''), 'Sin transportadora')")

	if params.BusinessID > 0 {
		q = q.Where("o.business_id = ?", params.BusinessID)
	}
	if params.From != nil {
		q = q.Where("s.created_at >= ?", *params.From)
	}
	if params.To != nil {
		q = q.Where("s.created_at < ?", *params.To)
	}
	if params.Carrier != "" {
		q = q.Where("LOWER(COALESCE(s.carrier, '')) = ?", strings.ToLower(params.Carrier))
	}

	var rows []row
	if err := q.Order("carrier ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}

	resp := &dtos.ProfitReportResponse{Rows: make([]dtos.ProfitReportRow, 0, len(rows))}
	for _, r := range rows {
		item := dtos.ProfitReportRow{
			Carrier:             r.Carrier,
			CarrierCode:         strings.ToLower(r.Carrier),
			Shipments:           r.Shipments,
			CarrierCostTotal:    r.CarrierCostTotal,
			CustomerChargeTotal: r.CustomerChargeTotal,
			ProfitTotal:         r.ProfitTotal,
		}
		resp.Rows = append(resp.Rows, item)
		resp.Totals.Shipments += item.Shipments
		resp.Totals.CarrierCostTotal += item.CarrierCostTotal
		resp.Totals.CustomerChargeTotal += item.CustomerChargeTotal
		resp.Totals.ProfitTotal += item.ProfitTotal
	}
	resp.Totals.Carrier = "TOTAL"
	resp.Totals.CarrierCode = "total"
	return resp, nil
}
