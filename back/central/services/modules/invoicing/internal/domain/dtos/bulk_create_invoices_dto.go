package dtos

// BulkCreateInvoicesDTO representa la petición para crear facturas masivamente
type BulkCreateInvoicesDTO struct {
	OrderIDs   []string
	BusinessID *uint // Enviado por super admin desde el frontend
}

// BulkCreateResult representa el resultado de la creación masiva de facturas (DEPRECADO - Síncrono)
// DEPRECATED: Usar sistema de jobs asíncronos en su lugar
type BulkCreateResult struct {
	Created int                 `json:"created"`
	Failed  int                 `json:"failed"`
	Results []BulkInvoiceResult `json:"results"`
}

// BulkInvoiceResult representa el resultado individual de cada factura (DEPRECADO - Síncrono)
// DEPRECATED: Usar sistema de jobs asíncronos en su lugar
type BulkInvoiceResult struct {
	OrderID   string  `json:"order_id"`
	Success   bool    `json:"success"`
	InvoiceID *uint   `json:"invoice_id,omitempty"`
	Error     *string `json:"error,omitempty"`
}

// BulkInvoiceJobMessage representa un mensaje RabbitMQ para procesar una factura en un job masivo
type BulkInvoiceJobMessage struct {
	JobID         string `json:"job_id"`
	OrderID       string `json:"order_id"`
	BusinessID    uint   `json:"business_id"`
	IsManual      bool   `json:"is_manual"`
	CreatedBy     *uint  `json:"created_by,omitempty"`
	AttemptNumber int    `json:"attempt_number"`
}
