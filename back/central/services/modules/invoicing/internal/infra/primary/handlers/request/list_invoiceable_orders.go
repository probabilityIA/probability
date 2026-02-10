package request

// ListInvoiceableOrdersRequest representa los parámetros de paginación
// para listar órdenes facturables
type ListInvoiceableOrdersRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}
