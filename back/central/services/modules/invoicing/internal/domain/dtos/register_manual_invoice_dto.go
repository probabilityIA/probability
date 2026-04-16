package dtos

// RegisterManualInvoiceDTO contiene los datos para registrar una factura externa (manual)
// Se usa cuando el cliente ya tiene una factura hecha por fuera del sistema
// y quiere asociarla a una orden para que no se facture de nuevo.
type RegisterManualInvoiceDTO struct {
	// Número de factura externa
	InvoiceNumber string

	// ID de la orden a asociar
	OrderID string

	// ID del negocio (obligatorio para super admin)
	BusinessID uint
}
