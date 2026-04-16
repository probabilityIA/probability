package request

// RegisterManualInvoice es el request para registrar una factura externa
type RegisterManualInvoice struct {
	InvoiceNumber string `json:"invoice_number" binding:"required"`
	OrderID       string `json:"order_id" binding:"required"`
}
