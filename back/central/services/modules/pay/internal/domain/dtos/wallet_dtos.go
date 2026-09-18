package dtos

// RechargeWalletDTO datos para solicitar una recarga
type RechargeWalletDTO struct {
	BusinessID uint
	Amount     float64
	Reference  string // Motivo/razón de la recarga
	Concept    string
}

// ManualDebitDTO datos para débito manual (admin)
type ManualDebitDTO struct {
	BusinessID uint
	Amount     float64
	Reference  string
	Concept    string
	UserID     *uint
}

// DebitForGuideDTO datos para débito por generación de guía
type DebitForGuideDTO struct {
	BusinessID     uint
	Amount         float64
	TrackingNumber string
	UserID         *uint
}

// FinancialStatsDTO datos para obtener estadísticas financieras
type FinancialStatsDTO struct {
	BusinessID *uint  // nil para todos los negocios
	StartDate  string // YYYY-MM-DD
	EndDate    string // YYYY-MM-DD
	Month      string // YYYY-MM (alternativa rápida)
}

// BusinessFinancialStats ingresos por negocio
type BusinessFinancialStats struct {
	BusinessID         uint    `json:"business_id"`
	BusinessName       string  `json:"business_name"`
	SubscriptionIncome float64 `json:"subscription_income"`
	GuideIncome        float64 `json:"guide_income"`
	GuideCount         int     `json:"guide_count"`
	TotalIncome        float64 `json:"total_income"`
}

// FinancialStatsResponse respuesta de estadísticas financieras
type FinancialStatsResponse struct {
	Period      PeriodInfo               `json:"period"`
	TotalIncome float64                  `json:"total_income"`
	Businesses  []BusinessFinancialStats `json:"businesses"`
}

// PeriodInfo información del período de la consulta
type PeriodInfo struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// SpendSummaryDTO filtro para el resumen de gasto de un negocio (propia plata,
// nunca margen de Probability)
type SpendSummaryDTO struct {
	BusinessID uint
	StartDate  string // YYYY-MM-DD
	EndDate    string // YYYY-MM-DD
}

// ConceptTotal total gastado por concepto (GUIDE, SUBSCRIPTION, etc.)
type ConceptTotal struct {
	Concept string  `json:"concept"`
	Count   int64   `json:"count"`
	Amount  float64 `json:"amount"`
}

// TransactionFilterDTO filtro paginado de transacciones de un negocio
type TransactionFilterDTO struct {
	BusinessID uint
	StartDate  string // YYYY-MM-DD
	EndDate    string // YYYY-MM-DD
	Concept    string
	Page       int
	PageSize   int
}
