package dtos

// GetInvoiceDTO contiene los parámetros para obtener una factura
type GetInvoiceDTO struct {
	// ID de la factura
	InvoiceID uint

	// Si se deben incluir los items
	IncludeItems bool
}
