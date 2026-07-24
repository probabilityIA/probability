package dtos

import "time"

type CreateConceptDTO struct {
	Code         string
	Name         string
	Description  string
	Kind         string
	IsRealIncome bool
}

type UpdateConceptDTO struct {
	ID           uint
	Name         string
	Description  string
	IsRealIncome bool
	IsActive     bool
}

type CreateTaxDTO struct {
	Code        string
	Name        string
	Description string
	RatePercent float64
	Kind        string
}

type UpdateTaxDTO struct {
	ID          uint
	Name        string
	Description string
	RatePercent float64
	Kind        string
	IsActive    bool
}

type SetConceptTaxDTO struct {
	ConceptID uint
	TaxID     uint
	IsActive  bool
}

type CreateEntryDTO struct {
	ConceptID   uint
	BusinessID  *uint
	EntryDate   time.Time
	Amount      float64
	Description string
}

type ListEntriesParams struct {
	From      *time.Time
	To        *time.Time
	ConceptID *uint
	Kind      string
	Page      int
	PageSize  int
}

func (p ListEntriesParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type SyncCandidate struct {
	SourceType  string
	SourceID    string
	BusinessID  *uint
	Amount      float64
	EntryDate   time.Time
	Description string
}

type SyncResult struct {
	Created  int            `json:"created"`
	BySource map[string]int `json:"by_source"`
}

type ReportParams struct {
	From time.Time
	To   time.Time
}

type ReportConceptRow struct {
	ConceptID    uint    `json:"concept_id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Kind         string  `json:"kind"`
	IsRealIncome bool    `json:"is_real_income"`
	EntriesCount int64   `json:"entries_count"`
	Amount       float64 `json:"amount"`
	TaxTotal     float64 `json:"tax_total"`
}

type ReportTaxRow struct {
	TaxCode string  `json:"tax_code"`
	TaxName string  `json:"tax_name"`
	Amount  float64 `json:"amount"`
}

type ReportTotals struct {
	RealIncome  float64 `json:"real_income"`
	CashIn      float64 `json:"cash_in"`
	RealExpense float64 `json:"real_expense"`
	CashOut     float64 `json:"cash_out"`
	TaxTotal    float64 `json:"tax_total"`
	Net         float64 `json:"net"`
}

type ReportResponse struct {
	From      string             `json:"from"`
	To        string             `json:"to"`
	Totals    ReportTotals       `json:"totals"`
	ByConcept []ReportConceptRow `json:"by_concept"`
	ByTax     []ReportTaxRow     `json:"by_tax"`
}
