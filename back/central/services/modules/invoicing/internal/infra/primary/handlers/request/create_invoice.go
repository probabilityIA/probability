package request

// CreateInvoice es el request para crear una factura manualmente
type CreateInvoice struct {
	OrderID  string `json:"order_id" binding:"required"`
	IsManual bool   `json:"is_manual"` // Por defecto true en creación manual
}
