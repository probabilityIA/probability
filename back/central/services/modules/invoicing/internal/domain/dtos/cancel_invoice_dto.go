package dtos

// CancelInvoiceDTO contiene los datos para cancelar una factura
type CancelInvoiceDTO struct {
	// ID de la factura a cancelar
	InvoiceID uint

	// Razón de la cancelación
	Reason string

	// ID del usuario que cancela
	CancelledByUserID uint
}
