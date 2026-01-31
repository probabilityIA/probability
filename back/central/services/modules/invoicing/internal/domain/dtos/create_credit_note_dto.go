package dtos

// CreateCreditNoteDTO contiene los datos para crear una nota de crédito
type CreateCreditNoteDTO struct {
	// ID de la factura asociada
	InvoiceID uint

	// Tipo de nota de crédito
	NoteType string

	// Monto de la nota de crédito
	Amount float64

	// Razón de la nota de crédito
	Reason string

	// Descripción detallada (opcional)
	Description *string

	// ID del usuario que crea la nota
	CreatedByUserID uint
}
