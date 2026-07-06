package dtos

type Credentials struct {
	Username  string
	AccessKey string
	AccountID string
	PartnerID string
	BaseURL   string
}

type CustomerData struct {
	Name    string
	Email   string
	Phone   string
	DNI     string
	Address string
}

type ItemData struct {
	ProductID   *string
	SKU         string
	Name        string
	Description *string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
	Tax         float64
	TaxRate     *float64
	Discount    float64
}

type CreateInvoiceRequest struct {
	Customer     CustomerData
	Items        []ItemData
	Total        float64
	Subtotal     float64
	Tax          float64
	Discount     float64
	ShippingCost float64
	Currency     string
	OrderID      string
	OrderNumber  string
	IsRetry      bool
	Credentials  Credentials
	Config       map[string]interface{}
}

type AuditData struct {
	RequestURL     string
	RequestPayload interface{}
	ResponseStatus int
	ResponseBody   string
}

type CreateInvoiceResult struct {
	InvoiceNumber  string
	ExternalID     string
	CUFE           string
	QRCode         string
	Total          string
	IssuedAt       string
	AlreadyExisted bool
	ProviderInfo   map[string]interface{}
	AuditData      *AuditData
}
