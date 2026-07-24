package request

type InvoiceItemRequest struct {
	Description string  `json:"description" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
}

type CreateInvoiceRequest struct {
	BusinessID uint                 `json:"business_id" binding:"required"`
	ConceptID  uint                 `json:"concept_id" binding:"required"`
	IssueDate  string               `json:"issue_date" binding:"required"`
	DueDate    string               `json:"due_date"`
	Notes      string               `json:"notes"`
	EmailTo    string               `json:"email_to"`
	TaxIDs     []uint               `json:"tax_ids"`
	Items      []InvoiceItemRequest `json:"items" binding:"required"`
}

type UpdateInvoiceRequest struct {
	ConceptID uint                 `json:"concept_id" binding:"required"`
	IssueDate string               `json:"issue_date" binding:"required"`
	DueDate   string               `json:"due_date"`
	Notes     string               `json:"notes"`
	EmailTo   string               `json:"email_to"`
	TaxIDs    []uint               `json:"tax_ids"`
	Items     []InvoiceItemRequest `json:"items" binding:"required"`
}

type SendInvoiceRequest struct {
	EmailTo string `json:"email_to"`
}
