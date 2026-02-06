package dtos

// BulkCreateInvoicesDTO representa la petición para crear facturas masivamente
type BulkCreateInvoicesDTO struct {
	OrderIDs []string `json:"order_ids" validate:"required,min=1,max=100"`
}

// BulkCreateResult representa el resultado de la creación masiva de facturas
type BulkCreateResult struct {
	Created int                 `json:"created"`
	Failed  int                 `json:"failed"`
	Results []BulkInvoiceResult `json:"results"`
}

// BulkInvoiceResult representa el resultado individual de cada factura
type BulkInvoiceResult struct {
	OrderID   string  `json:"order_id"`
	Success   bool    `json:"success"`
	InvoiceID *uint   `json:"invoice_id,omitempty"`
	Error     *string `json:"error,omitempty"`
}
