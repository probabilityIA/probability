package response

import "time"

// CreditNoteResponse representa la respuesta de creación de nota de crédito
type CreditNoteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    *CreditNoteData `json:"data,omitempty"`
}

// CreditNoteData contiene los datos de la nota de crédito creada
type CreditNoteData struct {
	NoteID     string    `json:"note_id"`
	NoteNumber string    `json:"note_number"`
	CUFE       string    `json:"cufe"`
	PDFURL     string    `json:"pdf_url"`
	XMLURL     string    `json:"xml_url"`
	NoteURL    string    `json:"note_url"`
	IssuedAt   time.Time `json:"issued_at"`
	Status     string    `json:"status"`
}
