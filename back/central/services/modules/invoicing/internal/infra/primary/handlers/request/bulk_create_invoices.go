package request

// BulkCreateInvoicesRequest representa la petición para crear facturas masivamente
type BulkCreateInvoicesRequest struct {
	OrderIDs []string `json:"order_ids" binding:"required,min=1,max=100"`
}
