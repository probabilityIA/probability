package request

// CreateBillCustomer datos del cliente para POST /v1/bills/validate
type CreateBillCustomer struct {
	Identification           string `json:"identification"`
	DV                       string `json:"dv,omitempty"`
	Company                  string `json:"company,omitempty"`
	TradeName                string `json:"trade_name,omitempty"`
	Names                    string `json:"names"`
	Address                  string `json:"address,omitempty"`
	Email                    string `json:"email,omitempty"`
	Phone                    string `json:"phone,omitempty"`
	LegalOrganizationID      string `json:"legal_organization_id"`
	TributeID                string `json:"tribute_id"`
	IdentificationDocumentID string `json:"identification_document_id"`
	MunicipalityID           string `json:"municipality_id"`
}

// CreateBillItem item de línea para POST /v1/bills/validate
type CreateBillItem struct {
	SchemeID       string  `json:"scheme_id"`
	Note           string  `json:"note,omitempty"`
	CodeReference  string  `json:"code_reference"`
	Name           string  `json:"name"`
	Quantity       int     `json:"quantity"`
	DiscountRate   float64 `json:"discount_rate"`
	Price          float64 `json:"price"`
	TaxRate        string  `json:"tax_rate"`
	UnitMeasureID  int     `json:"unit_measure_id"`
	StandardCodeID int     `json:"standard_code_id"`
	IsExcluded     int     `json:"is_excluded"`
	TributeID      int     `json:"tribute_id"`
}

// CreateBillBody body completo de POST /v1/bills/validate
type CreateBillBody struct {
	NumberingRangeID  int                `json:"numbering_range_id"`
	ReferenceCode     string             `json:"reference_code"`
	Observation       string             `json:"observation,omitempty"`
	PaymentForm       string             `json:"payment_form"`
	PaymentDueDate    string             `json:"payment_due_date,omitempty"`
	PaymentMethodCode string             `json:"payment_method_code"`
	OperationType     int                `json:"operation_type"`
	SendEmail         bool               `json:"send_email"`
	Document          string             `json:"document,omitempty"`
	Customer          CreateBillCustomer `json:"customer"`
	Items             []CreateBillItem   `json:"items"`
}
