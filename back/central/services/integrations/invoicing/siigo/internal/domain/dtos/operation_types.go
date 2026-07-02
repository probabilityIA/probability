package dtos

type InvoiceDetail struct {
	ID                     string
	Name                   string
	Prefix                 string
	Number                 int
	Date                   string
	CustomerID             string
	CustomerIdentification string
	CustomerBranchOffice   int
	Total                  float64
	Balance                float64
	Status                 string
	StampStatus            string
	CUFE                   string
	PublicURL              string
}

type StampError struct {
	Code    string
	Message string
}

type AnnulInvoiceResult struct {
	AuditData *AuditData
}

type ProductItem struct {
	ID                string
	Code              string
	Name              string
	Description       string
	Price             float64
	StockControl      bool
	AvailableQuantity float64
	Warehouses        []ProductWarehouseStock
}

type ProductWarehouseStock struct {
	ID       int
	Name     string
	Quantity float64
}

type WarehouseItem struct {
	ID   int
	Name string
}

type PaymentTypeItem struct {
	ID   int
	Name string
	Type string
}

type CreateCashReceiptRequest struct {
	InvoiceNumber string
	Credentials   Credentials
	Config        map[string]interface{}
}

type CreateCashReceiptResult struct {
	ReceiptID    string
	ReceiptName  string
	ProviderInfo map[string]interface{}
	AuditData    *AuditData
}

type CreateCreditNoteRequest struct {
	InvoiceExternalID string
	InvoiceNumber     string
	Amount            float64
	Reason            string
	NoteType          string
	CustomerDNI       string
	Credentials       Credentials
	Config            map[string]interface{}
}

type CreateCreditNoteResult struct {
	CreditNoteID     string
	CreditNoteNumber string
	CUFE             string
	ProviderInfo     map[string]interface{}
	AuditData        *AuditData
}
