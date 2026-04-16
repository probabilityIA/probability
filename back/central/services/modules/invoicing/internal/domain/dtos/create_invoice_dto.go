package dtos

// CreateInvoiceDTO contiene los datos necesarios para crear una factura
type CreateInvoiceDTO struct {
	// ID de la orden a facturar
	OrderID string

	// ID del proveedor de facturación a usar (opcional, se puede obtener de la config)
	InvoicingProviderID *uint

	// Notas adicionales (opcional)
	Notes *string

	// Si es generación manual (por defecto false = automática)
	IsManual bool

	// ID del usuario que crea la factura manualmente (solo si IsManual = true)
	CreatedByUserID *uint
}
