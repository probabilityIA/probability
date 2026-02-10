package messages

// BulkInvoiceJobMessage representa un mensaje RabbitMQ para procesar una factura en un job masivo
type BulkInvoiceJobMessage struct {
	JobID         string `json:"job_id"`
	OrderID       string `json:"order_id"`
	BusinessID    uint   `json:"business_id"`
	IsManual      bool   `json:"is_manual"`
	CreatedBy     *uint  `json:"created_by,omitempty"`
	AttemptNumber int    `json:"attempt_number"`
}
