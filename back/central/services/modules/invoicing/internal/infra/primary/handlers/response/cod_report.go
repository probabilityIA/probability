package response

type CODReport struct {
	PeriodLabel   string             `json:"period_label"`
	TotalInvoices int                `json:"total_invoices"`
	TotalAmount   float64            `json:"total_amount"`
	CODCount      int                `json:"cod_count"`
	CODAmount     float64            `json:"cod_amount"`
	NonCODCount   int                `json:"non_cod_count"`
	NonCODAmount  float64            `json:"non_cod_amount"`
	ByAccount     []AccountBreakdown `json:"by_account"`
}

type AccountBreakdown struct {
	AccountNumber string  `json:"account_number"`
	IsCOD         bool    `json:"is_cod"`
	Count         int     `json:"count"`
	Amount        float64 `json:"amount"`
}
