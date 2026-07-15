package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

type carrierAggRow struct {
	Carrier        string
	OrdersCount    int
	TotalCollected float64
}

func (r *Repository) AggregateByCarrier(ctx context.Context, f dtos.ReportFilter, collected bool) ([]entities.CarrierAggregate, error) {
	conds := []string{"o.deleted_at IS NULL", "o.cod_total > 0", "o.business_id = ?"}
	args := []any{f.BusinessID}

	if collected {
		conds = append(conds, paidExpr)
		if !f.StartDate.IsZero() && !f.EndDate.IsZero() {
			conds = append(conds, "EXISTS (SELECT 1 FROM cod_payment_cut_order cpo_r WHERE cpo_r.order_id = o.id AND cpo_r.deleted_at IS NULL AND cpo_r.paid_at BETWEEN ? AND ?)")
			args = append(args, f.StartDate, f.EndDate)
		}
	} else {
		conds = append(conds, "s.status = 'delivered'", "NOT "+paidExpr)
		if !f.StartDate.IsZero() && !f.EndDate.IsZero() {
			conds = append(conds, "COALESCE(s.delivered_at, s.updated_at) BETWEEN ? AND ?")
			args = append(args, f.StartDate, f.EndDate)
		}
	}
	if f.Carrier != "" {
		conds = append(conds, "UPPER(TRIM(COALESCE(s.carrier,''))) = ?")
		args = append(args, strings.ToUpper(strings.TrimSpace(f.Carrier)))
	}

	where := strings.Join(conds, " AND ")
	query := fmt.Sprintf(`
SELECT UPPER(TRIM(COALESCE(NULLIF(s.carrier,''),'SIN TRANSPORTADORA'))) AS carrier,
	COUNT(*) AS orders_count,
	COALESCE(SUM(o.cod_total),0) AS total_collected
FROM orders o %s
WHERE %s
GROUP BY 1
ORDER BY total_collected DESC`, latestShipmentJoin, where)

	var rows []carrierAggRow
	if err := r.db.Conn(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entities.CarrierAggregate, len(rows))
	for i := range rows {
		out[i] = entities.CarrierAggregate{
			Carrier:        rows[i].Carrier,
			OrdersCount:    rows[i].OrdersCount,
			TotalCollected: rows[i].TotalCollected,
		}
	}
	return out, nil
}

func (r *Repository) CutPeriodOrders(ctx context.Context, businessID uint, start, end time.Time) ([]entities.CarrierAggregate, error) {
	return r.AggregateByCarrier(ctx, dtos.ReportFilter{
		BusinessID: businessID,
		StartDate:  start,
		EndDate:    end,
	}, true)
}

type monthlyRow struct {
	Month          string
	OrdersCount    int
	TotalCollected float64
}

func (r *Repository) MonthlyHistory(ctx context.Context, businessID uint, months int) ([]entities.MonthlyPoint, error) {
	if months < 1 {
		months = 6
	}
	query := `
SELECT to_char(date_trunc('month', cpo.paid_at), 'YYYY-MM') AS month,
	COUNT(*) AS orders_count,
	COALESCE(SUM(cpo.cod_amount),0) AS total_collected
FROM cod_payment_cut_order cpo
WHERE cpo.deleted_at IS NULL AND cpo.business_id = ?
	AND cpo.paid_at >= date_trunc('month', now()) - make_interval(months => ?)
GROUP BY 1
ORDER BY 1`

	var rows []monthlyRow
	if err := r.db.Conn(ctx).Raw(query, businessID, months-1).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entities.MonthlyPoint, len(rows))
	for i := range rows {
		out[i] = entities.MonthlyPoint{
			Month:     rows[i].Month,
			Label:     monthLabel(rows[i].Month),
			Orders:    rows[i].OrdersCount,
			Collected: rows[i].TotalCollected,
		}
	}
	return out, nil
}

type weekRow struct {
	WeekStart      time.Time
	Carrier        string
	OrdersCount    int
	TotalCollected float64
}

func (r *Repository) WeeklyAggregates(ctx context.Context, businessID uint, weeks int) ([]entities.WeekAggregate, error) {
	if weeks < 1 {
		weeks = 8
	}
	query := fmt.Sprintf(`
SELECT date_trunc('week', COALESCE(s.delivered_at, s.updated_at))::date AS week_start,
	UPPER(TRIM(COALESCE(NULLIF(s.carrier,''),'SIN TRANSPORTADORA'))) AS carrier,
	COUNT(*) AS orders_count,
	COALESCE(SUM(o.cod_total),0) AS total_collected
FROM orders o %s
WHERE o.deleted_at IS NULL AND o.cod_total > 0 AND o.business_id = ?
	AND s.status = 'delivered'
	AND COALESCE(s.delivered_at, s.updated_at) >= date_trunc('week', now()) - make_interval(weeks => ?)
GROUP BY 1, 2
ORDER BY 1 DESC`, latestShipmentJoin)

	var rows []weekRow
	if err := r.db.Conn(ctx).Raw(query, businessID, weeks-1).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entities.WeekAggregate, len(rows))
	for i := range rows {
		out[i] = entities.WeekAggregate{
			WeekStart: rows[i].WeekStart,
			Carrier:   rows[i].Carrier,
			Orders:    rows[i].OrdersCount,
			Collected: rows[i].TotalCollected,
		}
	}
	return out, nil
}

var monthNamesES = []string{"", "Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"}

func monthLabel(ym string) string {
	if len(ym) != 7 {
		return ym
	}
	year := ym[0:4]
	var m int
	_, err := fmt.Sscanf(ym[5:7], "%d", &m)
	if err != nil || m < 1 || m > 12 {
		return ym
	}
	return monthNamesES[m] + " " + year
}
