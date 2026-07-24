package entities

import "time"

const (
	KindIncome  = "INCOME"
	KindExpense = "EXPENSE"

	SourceGuideMargin    = "GUIDE_MARGIN"
	SourceSubscription   = "SUBSCRIPTION"
	SourceWalletRecharge = "WALLET_RECHARGE"
	SourceCodPayout      = "COD_PAYOUT"
	SourceManual         = "MANUAL"
)

type Concept struct {
	ID           uint
	Code         string
	Name         string
	Description  string
	Kind         string
	IsRealIncome bool
	IsAutomatic  bool
	SourceType   string
	IsActive     bool
	Taxes        []ConceptTax
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Tax struct {
	ID          uint
	Code        string
	Name        string
	Description string
	RatePercent float64
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ConceptTax struct {
	ID          uint
	ConceptID   uint
	TaxID       uint
	TaxCode     string
	TaxName     string
	RatePercent float64
	IsActive    bool
}

type TaxLine struct {
	TaxCode     string
	TaxName     string
	RatePercent float64
	Amount      float64
}

type Entry struct {
	ID           uint
	ConceptID    uint
	ConceptCode  string
	ConceptName  string
	BusinessID   *uint
	BusinessName string
	EntryDate    time.Time
	Amount       float64
	Kind         string
	SourceType   string
	SourceID     string
	Description  string
	IsAutomatic  bool
	TaxTotal     float64
	TaxDetail    []TaxLine
	CreatedAt    time.Time
}
