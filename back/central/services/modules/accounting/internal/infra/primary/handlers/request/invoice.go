package request

type InvoiceItemRequest struct {
	Description string  `json:"description" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
}

type CreateInvoiceRequest struct {
	BusinessID       uint                 `json:"business_id" binding:"required"`
	ConceptID        uint                 `json:"concept_id" binding:"required"`
	IssueDate        string               `json:"issue_date" binding:"required"`
	DueDate          string               `json:"due_date"`
	Notes            string               `json:"notes"`
	EmailTo          string               `json:"email_to"`
	CustomerDocument string               `json:"customer_document"`
	CustomerPhone    string               `json:"customer_phone"`
	CustomerAddress  string               `json:"customer_address"`
	IsTest           bool                 `json:"is_test"`
	TaxIDs           []uint               `json:"tax_ids"`
	Items            []InvoiceItemRequest `json:"items" binding:"required"`
}

type UpdateInvoiceRequest struct {
	ConceptID        uint                 `json:"concept_id" binding:"required"`
	IssueDate        string               `json:"issue_date" binding:"required"`
	DueDate          string               `json:"due_date"`
	Notes            string               `json:"notes"`
	EmailTo          string               `json:"email_to"`
	CustomerDocument string               `json:"customer_document"`
	CustomerPhone    string               `json:"customer_phone"`
	CustomerAddress  string               `json:"customer_address"`
	IsTest           bool                 `json:"is_test"`
	TaxIDs           []uint               `json:"tax_ids"`
	Items            []InvoiceItemRequest `json:"items" binding:"required"`
}

type SendInvoiceRequest struct {
	EmailTo string `json:"email_to"`
}

type EmitInvoiceDianRequest struct {
	CustomerDocument string `json:"customer_document"`
	CustomerPhone    string `json:"customer_phone"`
	CustomerAddress  string `json:"customer_address"`
}

type DianConfigRequest struct {
	ApiURL           string `json:"api_url"`
	ClientID         string `json:"client_id" binding:"required"`
	ClientSecret     string `json:"client_secret" binding:"required"`
	Username         string `json:"username" binding:"required"`
	Password         string `json:"password" binding:"required"`
	NumberingRangeID int    `json:"numbering_range_id"`
	MunicipalityID   string `json:"municipality_id"`
	DocumentIDType   string `json:"identification_document_id"`
	LegalOrgID       string `json:"legal_organization_id"`
}
