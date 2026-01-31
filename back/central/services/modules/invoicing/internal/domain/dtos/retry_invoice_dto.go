package dtos

// RetryInvoiceDTO contiene los datos para reintentar una factura
type RetryInvoiceDTO struct {
	// ID de la factura a reintentar
	InvoiceID uint

	// Forzar reintento (ignorar contador de reintentos)
	Force bool
}
