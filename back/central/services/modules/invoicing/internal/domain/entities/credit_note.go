package entities

import "time"

// CreditNote representa una nota de crédito asociada a una factura
// Entidad PURA de dominio - SIN TAGS de infraestructura
type CreditNote struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Relaciones (solo IDs)
	InvoiceID  uint
	BusinessID uint

	// Identificadores
	CreditNoteNumber string
	ExternalID       *string
	InternalNumber   string

	// Tipo de nota de crédito
	NoteType string

	// Información financiera
	Amount   float64
	Currency string

	// Razón y descripción
	Reason      string
	Description *string

	// Estado
	Status string

	// Timestamps
	IssuedAt    *time.Time
	CancelledAt *time.Time

	// URLs y archivos
	NoteURL *string
	PDFURL  *string
	XMLURL  *string
	CUFE    *string

	// Metadata
	Metadata         map[string]interface{}
	ProviderResponse map[string]interface{}

	// Usuario que creó la nota
	CreatedByID uint
}
