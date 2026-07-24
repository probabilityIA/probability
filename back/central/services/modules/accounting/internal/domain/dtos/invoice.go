package dtos

import "time"

type InvoiceItemDTO struct {
	Description string
	Quantity    float64
	UnitPrice   float64
}

type CreateInvoiceDTO struct {
	BusinessID uint
	ConceptID  uint
	IssueDate  time.Time
	DueDate    *time.Time
	Notes      string
	EmailTo    string
	TaxIDs     []uint
	Items      []InvoiceItemDTO
}

type UpdateInvoiceDTO struct {
	ID        uint
	ConceptID uint
	IssueDate time.Time
	DueDate   *time.Time
	Notes     string
	EmailTo   string
	TaxIDs    []uint
	Items     []InvoiceItemDTO
}

type ListInvoicesParams struct {
	BusinessID *uint
	Status     string
	Page       int
	PageSize   int
}

func (p ListInvoicesParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type SendInvoiceDTO struct {
	ID      uint
	EmailTo string
}
