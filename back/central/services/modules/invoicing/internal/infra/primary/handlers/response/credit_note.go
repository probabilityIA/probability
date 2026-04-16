package response

import "time"

// CreditNote es la respuesta de una nota de crédito
type CreditNote struct {
	ID               uint                   `json:"id"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	InvoiceID        uint                   `json:"invoice_id"`
	BusinessID       uint                   `json:"business_id"`
	CreditNoteNumber string                 `json:"credit_note_number"`
	ExternalID       *string                `json:"external_id,omitempty"`
	Amount           float64                `json:"amount"`
	NoteType         string                 `json:"note_type"`
	Reason           string                 `json:"reason"`
	Status           string                 `json:"status"`
	IssuedAt         *time.Time             `json:"issued_at,omitempty"`
	CUFE             *string                `json:"cufe,omitempty"`
	PDFURL           *string                `json:"pdf_url,omitempty"`
	XMLURL           *string                `json:"xml_url,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}
