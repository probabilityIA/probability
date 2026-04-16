package messages

// InvoiceEvent representa un evento de factura (con tags para JSON)
type InvoiceEvent struct {
	EventType     string  `json:"event_type"`      // "invoice.created", "invoice.cancelled", "invoice.failed"
	InvoiceID     uint    `json:"invoice_id"`
	OrderID       string  `json:"order_id"`
	BusinessID    uint    `json:"business_id"`
	InvoiceNumber string  `json:"invoice_number"`
	TotalAmount   float64 `json:"total_amount"`
	Status        string  `json:"status"`
	ErrorMessage  string  `json:"error_message,omitempty"`
	Timestamp     string  `json:"timestamp"`
}

// CreditNoteEvent representa un evento de nota de crédito (con tags para JSON)
type CreditNoteEvent struct {
	EventType        string  `json:"event_type"` // "credit_note.created"
	CreditNoteID     uint    `json:"credit_note_id"`
	InvoiceID        uint    `json:"invoice_id"`
	BusinessID       uint    `json:"business_id"`
	CreditNoteNumber string  `json:"credit_note_number"`
	Amount           float64 `json:"amount"`
	NoteType         string  `json:"note_type"`
	Timestamp        string  `json:"timestamp"`
}
