package dtos

import "time"

type InvoiceItemDTO struct {
	Description string
	Quantity    float64
	UnitPrice   float64
}

type CreateInvoiceDTO struct {
	BusinessID       uint
	ConceptID        uint
	IssueDate        time.Time
	DueDate          *time.Time
	Notes            string
	EmailTo          string
	CustomerDocument string
	CustomerPhone    string
	CustomerAddress  string
	IsTest           bool
	TaxIDs           []uint
	Items            []InvoiceItemDTO
}

type UpdateInvoiceDTO struct {
	ID               uint
	ConceptID        uint
	IssueDate        time.Time
	DueDate          *time.Time
	Notes            string
	EmailTo          string
	CustomerDocument string
	CustomerPhone    string
	CustomerAddress  string
	IsTest           bool
	TaxIDs           []uint
	Items            []InvoiceItemDTO
}

type ListInvoicesParams struct {
	BusinessID *uint
	Status     string
	IsTest     *bool
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

type EmitInvoiceDianDTO struct {
	ID               uint
	CustomerDocument string
	CustomerPhone    string
	CustomerAddress  string
}

type DianConfigDTO struct {
	ApiURL            string
	ClientID          string
	ClientSecret      string
	Username          string
	Password          string
	NumberingRangeID  int
	MunicipalityID    string
	DocumentIDType    string
	LegalOrgID        string
	PaymentMethodCode string
}

type DianConfigStatus struct {
	Configured        bool   `json:"configured"`
	IntegrationID     uint   `json:"integration_id"`
	ApiURL            string `json:"api_url"`
	Username          string `json:"username"`
	NumberingRangeID  int    `json:"numbering_range_id"`
	MunicipalityID    string `json:"municipality_id"`
	DocumentIDType    string `json:"identification_document_id"`
	LegalOrgID        string `json:"legal_organization_id"`
	PaymentMethodCode string `json:"payment_method_code"`
}

type ClientProfileDTO struct {
	BusinessID        uint    `json:"business_id"`
	Configured        bool    `json:"configured"`
	DocumentType      string  `json:"document_type"`
	DocumentNumber    string  `json:"document_number"`
	DV                string  `json:"dv"`
	PersonType        string  `json:"person_type"`
	MunicipalityID    string  `json:"municipality_id"`
	Address           string  `json:"address"`
	Phone             string  `json:"phone"`
	BillingEmail      string  `json:"billing_email"`
	AppliesIVA        bool    `json:"applies_iva"`
	IVARate           float64 `json:"iva_rate"`
	AppliesReteFuente bool    `json:"applies_retefuente"`
	ReteFuenteRate    float64 `json:"retefuente_rate"`
	AppliesReteICA    bool    `json:"applies_reteica"`
	ReteICARate       float64 `json:"reteica_rate"`
	SuggestedTaxIDs   []uint  `json:"suggested_tax_ids"`
}
