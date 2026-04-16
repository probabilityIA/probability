package dtos

// Credentials contiene las credenciales de autenticación de World Office
type Credentials struct {
	Username    string
	Password    string
	CompanyCode string // Código de la empresa en World Office
	BaseURL     string // URL base de la API (opcional, usa el default del cliente si está vacío)
}

// CustomerData datos del cliente para la API de World Office
type CustomerData struct {
	Name    string
	Email   string
	Phone   string
	DNI     string
	Address string
}

// ItemData datos de un item para la API de World Office
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

// CreateInvoiceRequest datos tipados para crear una factura en World Office
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
	Credentials  Credentials
	Config       map[string]interface{}
}

// AuditData captura el request/response HTTP para logging y debugging
type AuditData struct {
	RequestURL     string
	RequestPayload interface{}
	ResponseStatus int
	ResponseBody   string
}

// CreateInvoiceResult resultado de crear una factura en World Office
// Se retorna siempre (incluso en error) para incluir audit data
type CreateInvoiceResult struct {
	InvoiceNumber string
	ExternalID    string
	CUFE          string
	QRCode        string
	Total         string
	IssuedAt      string
	ProviderInfo  map[string]interface{}
	AuditData     *AuditData
}
