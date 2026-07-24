package entities

import "time"

const (
	InvoiceStatusDraft     = "DRAFT"
	InvoiceStatusSent      = "SENT"
	InvoiceStatusPaid      = "PAID"
	InvoiceStatusCancelled = "CANCELLED"

	SourceInvoice = "INVOICE"

	DianStatusNone      = "NONE"
	DianStatusValidated = "VALIDATED"
)

type InvoiceItem struct {
	ID          uint
	Description string
	Quantity    float64
	UnitPrice   float64
	Amount      float64
}

type Invoice struct {
	ID           uint
	Number       string
	BusinessID   uint
	BusinessName string
	ConceptID    uint
	ConceptName  string
	IssueDate    time.Time
	DueDate      *time.Time
	Status       string
	Notes        string
	Subtotal     float64
	TaxTotal     float64
	Total        float64
	TaxDetail    []TaxLine
	EmailTo      string
	SentAt       *time.Time
	PaidAt       *time.Time
	Items        []InvoiceItem

	CustomerDocument string
	CustomerPhone    string
	CustomerAddress  string
	DianStatus       string
	Cufe             string
	DianNumber       string
	DianExternalID   string
	DianQR           string
	DianEmittedAt    *time.Time
	IsTest           bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type DianEmitResult struct {
	Number     string
	CUFE       string
	QRCode     string
	ExternalID string
	IssuedAt   string
}
