package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

type carrierDetailRow struct {
	Carrier         string
	Orders          int
	EnCurso         float64
	EnCursoOrders   int
	Entregado       float64
	EntregadoOrders int
	PorPagar        float64
	Recaudado       float64
	Cargo           float64
	Total           float64
}

func (r *Repository) SummaryCarrierDetail(ctx context.Context, f dtos.ReportFilter) ([]entities.CarrierDetail, error) {
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
	where := strings.Join(conds, " AND ")

	query := fmt.Sprintf(`
SELECT UPPER(TRIM(COALESCE(NULLIF(s.carrier,''),'SIN TRANSPORTADORA'))) AS carrier,
	COUNT(*) AS orders,
	COALESCE(SUM(CASE WHEN s.status IN (%[1]s) THEN o.cod_total ELSE 0 END),0) AS en_curso,
	COALESCE(SUM(CASE WHEN s.status IN (%[1]s) THEN 1 ELSE 0 END),0) AS en_curso_orders,
	COALESCE(SUM(CASE WHEN s.status = 'delivered' THEN o.cod_total ELSE 0 END),0) AS entregado,
	COALESCE(SUM(CASE WHEN s.status = 'delivered' THEN 1 ELSE 0 END),0) AS entregado_orders,
	COALESCE(SUM(CASE WHEN s.status = 'delivered' AND NOT %[2]s THEN o.cod_total ELSE 0 END),0) AS por_pagar,
	COALESCE(SUM(CASE WHEN %[2]s THEN o.cod_total ELSE 0 END),0) AS recaudado,
	COALESCE(SUM(COALESCE(s.cod_carrier_fee,0)),0) AS cargo,
	COALESCE(SUM(o.cod_total),0) AS total
FROM orders o %[3]s
WHERE %[4]s
GROUP BY 1
ORDER BY total DESC`, pendingStatuses, paidExpr, latestShipmentJoin, where)

	var rows []carrierDetailRow
	if err := r.db.Conn(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entities.CarrierDetail, len(rows))
	for i := range rows {
		out[i] = entities.CarrierDetail{
			Carrier:         rows[i].Carrier,
			Orders:          rows[i].Orders,
			EnCurso:         rows[i].EnCurso,
			EnCursoOrders:   rows[i].EnCursoOrders,
			Entregado:       rows[i].Entregado,
			EntregadoOrders: rows[i].EntregadoOrders,
			PorPagar:        rows[i].PorPagar,
			Recaudado:       rows[i].Recaudado,
			Cargo:           rows[i].Cargo,
			Total:           rows[i].Total,
		}
	}
	return out, nil
}

type historyRow struct {
	Bucket    string
	Entregado float64
	EnCurso   float64
}

func historyBucketExprs(bucket string) (truncTS string, label string) {
	const ts = "COALESCE(s.delivered_at, o.created_at)"
	switch bucket {
	case "week":
		return fmt.Sprintf("date_trunc('week', %s)", ts), "to_char(bkt, 'DD Mon')"
	case "month":
		return fmt.Sprintf("date_trunc('month', %s)", ts), "to_char(bkt, 'Mon YYYY')"
	case "quarter":
		return fmt.Sprintf("date_trunc('quarter', %s)", ts), "'Q' || to_char(bkt, 'Q') || ' ' || to_char(bkt, 'YYYY')"
	case "semester":
		return fmt.Sprintf("date_trunc('year', %[1]s) + floor((extract(month from %[1]s) - 1) / 6) * interval '6 months'", ts),
			"'S' || (CASE WHEN extract(month from bkt) <= 6 THEN '1' ELSE '2' END) || ' ' || to_char(bkt, 'YYYY')"
	case "year":
		return fmt.Sprintf("date_trunc('year', %s)", ts), "to_char(bkt, 'YYYY')"
	default:
		return fmt.Sprintf("date_trunc('day', %s)", ts), "to_char(bkt, 'DD Mon')"
	}
}

func (r *Repository) SummaryHistory(ctx context.Context, f dtos.ReportFilter) ([]entities.HistoryPoint, error) {
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
	where := strings.Join(conds, " AND ")
	truncTS, label := historyBucketExprs(f.Bucket)

	query := fmt.Sprintf(`
SELECT %[5]s AS bucket,
	COALESCE(SUM(entregado),0) AS entregado,
	COALESCE(SUM(en_curso),0) AS en_curso
FROM (
	SELECT %[4]s AS bkt,
		CASE WHEN s.status = 'delivered' THEN o.cod_total ELSE 0 END AS entregado,
		CASE WHEN s.status IN (%[1]s) THEN o.cod_total ELSE 0 END AS en_curso
	FROM orders o %[2]s
	WHERE %[3]s
) t
GROUP BY bkt
ORDER BY bkt`, pendingStatuses, latestShipmentJoin, where, truncTS, label)

	var rows []historyRow
	if err := r.db.Conn(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entities.HistoryPoint, len(rows))
	for i := range rows {
		out[i] = entities.HistoryPoint{
			Label:     rows[i].Bucket,
			Entregado: rows[i].Entregado,
			EnCurso:   rows[i].EnCurso,
		}
	}
	return out, nil
}
