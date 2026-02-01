package entities

import "time"

// ═══════════════════════════════════════════════════════════════
// ESTADÍSTICAS Y RESÚMENES - Entidades PURAS de dominio
// ═══════════════════════════════════════════════════════════════

// InvoiceSummary representa un resumen general de facturas con KPIs principales
type InvoiceSummary struct {
	Period         PeriodInfo
	Totals         TotalStats
	ByStatus       []StatusBreakdown
	ByProvider     []ProviderBreakdown
	RecentFailures []FailureDetail
}

// PeriodInfo contiene la información del período de tiempo analizado
type PeriodInfo struct {
	Start time.Time
	End   time.Time
	Label string
}

// TotalStats representa las métricas totales de facturación
type TotalStats struct {
	TotalInvoices int
	TotalAmount   float64
	IssuedCount   int
	IssuedAmount  float64
	FailedCount   int
	PendingCount  int
}

// StatusBreakdown representa el desglose de facturas por estado
type StatusBreakdown struct {
	Status     string
	Count      int
	Amount     float64
	Percentage float64
}

// ProviderBreakdown representa el desglose de facturas por proveedor
type ProviderBreakdown struct {
	ProviderID   uint
	ProviderName string
	Count        int
	Amount       float64
}

// FailureDetail representa el detalle de una factura fallida
type FailureDetail struct {
	InvoiceID uint
	OrderID   string
	Amount    float64
	Error     string
	FailedAt  time.Time
}

// DetailedStats representa estadísticas detalladas para dashboards
type DetailedStats struct {
	Summary          StatsSummary
	TopCustomers     []CustomerStats
	MonthlyBreakdown []MonthlyStats
	FailureAnalysis  FailureAnalysis
	ProcessingTimes  ProcessingStats
}

// StatsSummary contiene el resumen estadístico general
type StatsSummary struct {
	TotalInvoices int
	TotalAmount   float64
	AvgAmount     float64
	SuccessRate   float64
}

// CustomerStats representa estadísticas de facturación por cliente
type CustomerStats struct {
	CustomerName string
	InvoiceCount int
	TotalAmount  float64
}

// MonthlyStats representa estadísticas mensuales de facturación
type MonthlyStats struct {
	Month       string
	Count       int
	Amount      float64
	SuccessRate float64
}

// FailureAnalysis representa el análisis de facturas fallidas
type FailureAnalysis struct {
	TotalFailures int
	ByReason      []FailureReason
}

// FailureReason representa la categorización de errores
type FailureReason struct {
	Reason     string
	Count      int
	Percentage float64
}

// ProcessingStats representa estadísticas de tiempos de procesamiento
type ProcessingStats struct {
	AvgSeconds float64
	P50Seconds float64
	P95Seconds float64
	P99Seconds float64
}

// TrendData representa datos de tendencias temporales para gráficos
type TrendData struct {
	Metric      string
	Granularity string
	DataPoints  []TrendPoint
	Trend       TrendInfo
}

// TrendPoint representa un punto en la serie temporal
type TrendPoint struct {
	Date        string
	Value       float64
	SuccessRate float64
}

// TrendInfo representa información sobre la tendencia
type TrendInfo struct {
	Direction        string
	PercentageChange float64
	ComparisonPeriod string
}
