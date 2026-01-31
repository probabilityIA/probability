package request

// CreditNoteRequest representa la solicitud de creación de nota de crédito
type CreditNoteRequest struct {
	InvoiceNumber string  `json:"invoice_number"` // Número de la factura a la que aplica
	Referer       string  `json:"referer"`        // NIT del facturador
	Amount        float64 `json:"amount"`         // Monto de la nota de crédito
	Reason        string  `json:"reason"`         // Razón de la nota de crédito
	Description   string  `json:"description,omitempty"` // Descripción detallada
	NoteType      string  `json:"note_type"`      // Tipo: "01" anulación, "02" corrección, "03" devolución
}
