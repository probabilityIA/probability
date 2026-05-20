package dtos

type MarkInvoiceAsCancelledDTO struct {
	InvoiceID         uint
	Reason            string
	CancelledByUserID uint
}
