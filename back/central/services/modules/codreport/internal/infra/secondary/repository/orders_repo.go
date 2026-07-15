package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

type codOrderRow struct {
	OrderID       string
	OrderNumber   string
	CustomerName  string
	Carrier       string
	CodTotal      float64
	CodCarrierFee float64
	ShippingCost  float64
	Currency      string
	Status        string
	Collected     bool
	ShipmentID    uint
	HasGuide      bool
	Paid          bool
	CreatedAt     time.Time
	DeliveredAt   *time.Time
}

const hasGuideExpr = `(COALESCE(NULLIF(s.guide_id,''),'') <> '' OR COALESCE(NULLIF(s.guide_url,''),'') <> '' OR COALESCE(NULLIF(s.probability_guide_url,''),'') <> '')`

const paidExpr = `EXISTS (SELECT 1 FROM cod_payment_cut_order cpo WHERE cpo.order_id = o.id AND cpo.deleted_at IS NULL)`

func (r *Repository) ListCodOrders(ctx context.Context, f dtos.OrdersFilter) ([]entities.CodOrder, int64, error) {
	conds := []string{"o.deleted_at IS NULL", "o.cod_total > 0", "o.business_id = ?"}
	args := []any{f.BusinessID}

	if !f.StartDate.IsZero() && !f.EndDate.IsZero() {
		conds = append(conds, "COALESCE(s.delivered_at, o.created_at) BETWEEN ? AND ?")
		args = append(args, f.StartDate, f.EndDate)
	}
	if f.Carrier != "" {
		conds = append(conds, "UPPER(TRIM(COALESCE(s.carrier,''))) = ?")
		args = append(args, strings.ToUpper(strings.TrimSpace(f.Carrier)))
	}
	if f.Collected != nil {
		if *f.Collected {
			conds = append(conds, paidExpr)
		} else {
			conds = append(conds, "NOT "+paidExpr)
		}
	}
	if f.Search != "" {
		conds = append(conds, "(o.order_number ILIKE ? OR o.customer_name ILIKE ?)")
		like := "%" + strings.TrimSpace(f.Search) + "%"
		args = append(args, like, like)
	}
	if f.HasGuide != nil {
		if *f.HasGuide {
			conds = append(conds, hasGuideExpr)
		} else {
			conds = append(conds, "NOT "+hasGuideExpr)
		}
	}

	where := strings.Join(conds, " AND ")

	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM orders o %s WHERE %s`, latestShipmentJoin, where)
	var total int64
	if err := r.db.Conn(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	page := f.Page
	if page < 1 {
		page = 1
	}
	pageSize := f.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	listSQL := fmt.Sprintf(`
SELECT o.id AS order_id, o.order_number, o.customer_name, o.cod_total, o.currency, o.created_at,
	s.id AS shipment_id,
	UPPER(TRIM(COALESCE(NULLIF(s.carrier,''),'SIN TRANSPORTADORA'))) AS carrier,
	COALESCE(s.shipping_cost,0) AS shipping_cost,
	COALESCE(s.cod_carrier_fee,0) AS cod_carrier_fee,
	s.status, s.delivered_at,
	(s.status = 'delivered') AS collected,
	`+paidExpr+` AS paid,
	`+hasGuideExpr+` AS has_guide
FROM orders o %s
WHERE %s
ORDER BY COALESCE(s.delivered_at, o.created_at) DESC
LIMIT ? OFFSET ?`, latestShipmentJoin, where)

	listArgs := append(append([]any{}, args...), pageSize, offset)
	var rows []codOrderRow
	if err := r.db.Conn(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	out := make([]entities.CodOrder, len(rows))
	for i := range rows {
		out[i] = entities.CodOrder{
			OrderID:       rows[i].OrderID,
			OrderNumber:   rows[i].OrderNumber,
			ShipmentID:    rows[i].ShipmentID,
			HasGuide:      rows[i].HasGuide,
			CustomerName:  rows[i].CustomerName,
			Carrier:       rows[i].Carrier,
			CodTotal:      rows[i].CodTotal,
			CodCarrierFee: rows[i].CodCarrierFee,
			ShippingCost:  rows[i].ShippingCost,
			Currency:      rows[i].Currency,
			Status:        rows[i].Status,
			Collected:     rows[i].Collected,
			Paid:          rows[i].Paid,
			CreatedAt:     rows[i].CreatedAt,
			DeliveredAt:   rows[i].DeliveredAt,
		}
	}
	return out, total, nil
}
