package dtos

const MaxBulkInvoiceOrders = 1000

type BulkCreateInvoicesDTO struct {
	OrderIDs   []string
	BusinessID *uint
}

type BulkCreateResult struct {
	Created int                 `json:"created"`
	Failed  int                 `json:"failed"`
	Results []BulkInvoiceResult `json:"results"`
}

type BulkInvoiceResult struct {
	OrderID   string  `json:"order_id"`
	Success   bool    `json:"success"`
	InvoiceID *uint   `json:"invoice_id,omitempty"`
	Error     *string `json:"error,omitempty"`
}

type BulkInvoiceJobMessage struct {
	JobID         string `json:"job_id"`
	OrderID       string `json:"order_id"`
	BusinessID    uint   `json:"business_id"`
	IsManual      bool   `json:"is_manual"`
	CreatedBy     *uint  `json:"created_by,omitempty"`
	AttemptNumber int    `json:"attempt_number"`
}
