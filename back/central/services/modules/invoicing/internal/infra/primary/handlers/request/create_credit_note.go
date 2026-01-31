package request

// CreateCreditNote es el request para crear una nota de crédito
type CreateCreditNote struct {
	InvoiceID uint    `json:"invoice_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Reason    string  `json:"reason" binding:"required,min=3,max=500"`
	NoteType  string  `json:"note_type" binding:"required,oneof=cancellation correction full_refund partial_refund"`
}
