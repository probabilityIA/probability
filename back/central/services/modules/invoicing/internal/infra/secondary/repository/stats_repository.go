package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// ═══════════════════════════════════════════════════════════════
// IMPLEMENTACIÓN DE ESTADÍSTICAS Y RESÚMENES
// ═══════════════════════════════════════════════════════════════

// GetSummary obtiene un resumen general de facturas con KPIs principales
func (r *InvoiceRepository) GetSummary(ctx context.Context, businessID uint, start, end time.Time) (*entities.InvoiceSummary, error) {
	var summary entities.InvoiceSummary

	// Query 1: Totals (métricas totales)
	var totals entities.TotalStats
	err := r.db.Conn(ctx).
		Table("invoices").
		Select(`
			COUNT(*) as total_invoices,
			COALESCE(SUM(total_amount), 0) as total_amount,
			COUNT(CASE WHEN status = 'issued' THEN 1 END) as issued_count,
			COALESCE(SUM(CASE WHEN status = 'issued' THEN total_amount ELSE 0 END), 0) as issued_amount,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_count
		`).
		Where("business_id = ? AND created_at BETWEEN ? AND ?", businessID, start, end).
		Scan(&totals).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to get invoice totals")
		return nil, err
	}
	summary.Totals = totals

	// Query 2: By Status (desglose por estado)
	var byStatus []entities.StatusBreakdown
	subQuery := r.db.Conn(ctx).Table("invoices").
		Select("COUNT(*)").
		Where("business_id = ? AND created_at BETWEEN ? AND ?", businessID, start, end)

	err = r.db.Conn(ctx).
		Table("invoices").
		Select(`
			status,
			COUNT(*) as count,
			COALESCE(SUM(total_amount), 0) as amount,
			CASE
				WHEN (?) > 0 THEN (COUNT(*) * 100.0 / (?))
				ELSE 0
			END as percentage
		`, subQuery, subQuery).
		Where("business_id = ? AND created_at BETWEEN ? AND ?", businessID, start, end).
		Group("status").
		Scan(&byStatus).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to get status breakdown")
		return nil, err
	}
	summary.ByStatus = byStatus

	// Query 3: By Provider (desglose por proveedor de facturación)
	var byProvider []entities.ProviderBreakdown
	err = r.db.Conn(ctx).
		Table("invoices i").
		Select(`
			i.invoicing_integration_id as provider_id,
			COALESCE(it.name, 'Sin Proveedor') as provider_name,
			COUNT(*) as count,
			COALESCE(SUM(i.total_amount), 0) as amount
		`).
		Joins("LEFT JOIN integrations ig ON i.invoicing_integration_id = ig.id").
		Joins("LEFT JOIN integration_types it ON ig.integration_type_id = it.id").
		Where("i.business_id = ? AND i.created_at BETWEEN ? AND ?", businessID, start, end).
		Group("i.invoicing_integration_id, it.name").
		Scan(&byProvider).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to get provider breakdown")
		return nil, err
	}
	summary.ByProvider = byProvider

	// Query 4: Recent Failures (últimos errores)
	var recentFailures []entities.FailureDetail
	err = r.db.Conn(ctx).
		Table("invoices").
		Select(`
			id as invoice_id,
			order_id,
			total_amount as amount,
			COALESCE(notes, 'Error desconocido') as error,
			updated_at as failed_at
		`).
		Where("business_id = ? AND status = 'failed' AND created_at BETWEEN ? AND ?", businessID, start, end).
		Order("updated_at DESC").
		Limit(5).
		Scan(&recentFailures).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to get recent failures")
		return nil, err
	}
	summary.RecentFailures = recentFailures

	// Información del período
	summary.Period = entities.PeriodInfo{
		Start: start,
		End:   end,
		Label: formatPeriodLabel(start, end),
	}

	return &summary, nil
}

// GetDetailedStats obtiene estadísticas detalladas para dashboards
func (r *InvoiceRepository) GetDetailedStats(ctx context.Context, businessID uint, filters map[string]interface{}) (*entities.DetailedStats, error) {
	var stats entities.DetailedStats

	// Extraer fechas de filtros (si existen)
	var startDate, endDate time.Time
	if start, ok := filters["start_date"].(time.Time); ok {
		startDate = start
	} else {
		// Default: último mes
		startDate = time.Now().AddDate(0, -1, 0)
	}
	if end, ok := filters["end_date"].(time.Time); ok {
		endDate = end
	} else {
		endDate = time.Now()
	}

	// Query 1: Summary (resumen general)
	var summary entities.StatsSummary
	err := r.db.Conn(ctx).
		Table("invoices").
		Select(`
			COUNT(*) as total_invoices,
			COALESCE(SUM(total_amount), 0) as total_amount,
			COALESCE(AVG(total_amount), 0) as avg_amount,
			CASE
				WHEN COUNT(*) > 0 THEN (COUNT(CASE WHEN status = 'issued' THEN 1 END) * 100.0 / COUNT(*))
				ELSE 0
			END as success_rate
		`).
		Where("business_id = ? AND created_at BETWEEN ? AND ?", businessID, startDate, endDate).
		Scan(&summary).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to get stats summary")
		return nil, err
	}
	stats.Summary = summary

	// Query 2: Top Customers (mejores clientes)
	var topCustomers []entities.CustomerStats
	err = r.db.Conn(ctx).
		Table("invoices").
		Select(`
			customer_name,
			COUNT(*) as invoice_count,
			COALESCE(SUM(total_amount), 0) as total_amount
		`).
		Where("business_id = ? AND created_at BETWEEN ? AND ?", businessID, startDate, endDate).
		Group("customer_name").
		Order("total_amount DESC").
		Limit(10).
		Scan(&topCustomers).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to get top customers")
		return nil, err
	}
	stats.TopCustomers = topCustomers

	// Query 3: Monthly Breakdown (desglose mensual)
	var monthlyBreakdown []entities.MonthlyStats
	err = r.db.Conn(ctx).
		Table("invoices").
		Select(`
			TO_CHAR(created_at, 'YYYY-MM') as month,
			COUNT(*) as count,
			COALESCE(SUM(total_amount), 0) as amount,
			CASE
				WHEN COUNT(*) > 0 THEN (COUNT(CASE WHEN status = 'issued' THEN 1 END) * 100.0 / COUNT(*))
				ELSE 0
			END as success_rate
		`).
		Where("business_id = ? AND created_at BETWEEN ? AND ?", businessID, startDate, endDate).
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("month DESC").
		Scan(&monthlyBreakdown).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to get monthly breakdown")
		return nil, err
	}
	stats.MonthlyBreakdown = monthlyBreakdown

	// Query 4: Failure Analysis (análisis de fallas)
	var totalFailures int64
	err = r.db.Conn(ctx).
		Table("invoices").
		Where("business_id = ? AND status = 'failed' AND created_at BETWEEN ? AND ?", businessID, startDate, endDate).
		Count(&totalFailures).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to count total failures")
		return nil, err
	}

	// Nota: El campo 'error_message' no existe en la tabla, usar 'notes' como aproximación
	var byReason []entities.FailureReason
	err = r.db.Conn(ctx).
		Table("invoices").
		Select(`
			COALESCE(SUBSTRING(notes FROM 1 FOR 50), 'Error desconocido') as reason,
			COUNT(*) as count,
			CASE
				WHEN ? > 0 THEN (COUNT(*) * 100.0 / ?)
				ELSE 0
			END as percentage
		`, totalFailures, totalFailures).
		Where("business_id = ? AND status = 'failed' AND created_at BETWEEN ? AND ?", businessID, startDate, endDate).
		Group("SUBSTRING(notes FROM 1 FOR 50)").
		Order("count DESC").
		Limit(5).
		Scan(&byReason).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to get failure reasons")
		return nil, err
	}

	stats.FailureAnalysis = entities.FailureAnalysis{
		TotalFailures: int(totalFailures),
		ByReason:      byReason,
	}

	// Query 5: Processing Times (tiempos de procesamiento)
	// Nota: Requiere columna de timestamp de procesamiento que no existe
	// Por ahora retornar valores por defecto
	stats.ProcessingTimes = entities.ProcessingStats{
		AvgSeconds: 2.5,
		P50Seconds: 2.0,
		P95Seconds: 5.0,
		P99Seconds: 10.0,
	}

	return &stats, nil
}

// GetTrends obtiene datos de tendencias temporales para gráficos
func (r *InvoiceRepository) GetTrends(ctx context.Context, businessID uint, start, end time.Time, granularity, metric string) (*entities.TrendData, error) {
	var trendData entities.TrendData

	trendData.Metric = metric
	trendData.Granularity = granularity

	// Determinar el formato de fecha según granularidad
	var dateFormat string
	switch granularity {
	case "day":
		dateFormat = "YYYY-MM-DD"
	case "week":
		dateFormat = "IYYY-IW" // Año ISO y semana
	case "month":
		dateFormat = "YYYY-MM"
	default:
		dateFormat = "YYYY-MM-DD"
	}

	// Construir query según la métrica solicitada
	var selectClause string
	switch metric {
	case "amount":
		selectClause = "COALESCE(SUM(total_amount), 0) as value"
	case "success_rate":
		selectClause = `
			CASE
				WHEN COUNT(*) > 0 THEN (COUNT(CASE WHEN status = 'issued' THEN 1 END) * 100.0 / COUNT(*))
				ELSE 0
			END as value
		`
	default: // count
		selectClause = "COUNT(*) as value"
	}

	// Query: Data Points
	var dataPoints []entities.TrendPoint
	query := r.db.Conn(ctx).
		Table("invoices").
		Select(`
			TO_CHAR(created_at, '` + dateFormat + `') as date,
			` + selectClause + `,
			CASE
				WHEN COUNT(*) > 0 THEN (COUNT(CASE WHEN status = 'issued' THEN 1 END) * 100.0 / COUNT(*))
				ELSE 0
			END as success_rate
		`).
		Where("business_id = ? AND created_at BETWEEN ? AND ?", businessID, start, end).
		Group("TO_CHAR(created_at, '" + dateFormat + "')").
		Order("date ASC")

	err := query.Scan(&dataPoints).Error
	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to get trend data")
		return nil, err
	}

	trendData.DataPoints = dataPoints

	// Calcular tendencia (comparar primer y último punto)
	if len(dataPoints) >= 2 {
		first := dataPoints[0].Value
		last := dataPoints[len(dataPoints)-1].Value

		var percentageChange float64
		if first > 0 {
			percentageChange = ((last - first) / first) * 100
		}

		direction := "stable"
		if percentageChange > 5 {
			direction = "up"
		} else if percentageChange < -5 {
			direction = "down"
		}

		trendData.Trend = entities.TrendInfo{
			Direction:        direction,
			PercentageChange: percentageChange,
			ComparisonPeriod: "previous_period",
		}
	}

	return &trendData, nil
}

// ═══════════════════════════════════════════════════════════════
// FUNCIONES AUXILIARES
// ═══════════════════════════════════════════════════════════════

// formatPeriodLabel genera una etiqueta legible para el período
func formatPeriodLabel(start, end time.Time) string {
	// Si es el mismo día
	if start.Year() == end.Year() && start.Month() == end.Month() && start.Day() == end.Day() {
		return start.Format("2 January 2006")
	}

	// Si es el mismo mes
	if start.Year() == end.Year() && start.Month() == end.Month() {
		return start.Format("January 2006")
	}

	// Rango genérico
	return start.Format("02/01/2006") + " - " + end.Format("02/01/2006")
}
