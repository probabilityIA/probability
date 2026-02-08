package domain

import "time"

// ───────────────────────────────────────────
//
//	INVOICE EVENT TYPES
//
// ───────────────────────────────────────────

// InvoiceEventType define los tipos de eventos de facturación
type InvoiceEventType string

const (
	InvoiceEventTypeCreated       InvoiceEventType = "invoice.created"
	InvoiceEventTypeFailed        InvoiceEventType = "invoice.failed"
	InvoiceEventTypeCancelled     InvoiceEventType = "invoice.cancelled"
	CreditNoteEventTypeCreated    InvoiceEventType = "credit_note.created"
	BulkJobEventTypeProgress      InvoiceEventType = "bulk_job.progress"
	BulkJobEventTypeCompleted     InvoiceEventType = "bulk_job.completed"
)

// IsValid verifica si el tipo de evento es válido
func (t InvoiceEventType) IsValid() bool {
	switch t {
	case InvoiceEventTypeCreated, InvoiceEventTypeFailed, InvoiceEventTypeCancelled,
		CreditNoteEventTypeCreated, BulkJobEventTypeProgress, BulkJobEventTypeCompleted:
		return true
	}
	return false
}

// ───────────────────────────────────────────
//
//	INVOICE EVENT STRUCTURES
//
// ───────────────────────────────────────────

// InvoiceEvent representa un evento de facturación recibido desde Redis
type InvoiceEvent struct {
	ID         string                 `json:"id"`
	Type       InvoiceEventType       `json:"event_type"`
	BusinessID uint                   `json:"business_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
}
