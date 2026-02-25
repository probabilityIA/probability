package dtos

// Credentials contiene las credenciales de autenticación de Softpymes
type Credentials struct {
	APIKey    string
	APISecret string
}

// CustomerData datos del cliente para la API de Softpymes
type CustomerData struct {
	Name    string
	Email   string
	Phone   string
	DNI     string
	Address string
}

// ItemData datos de un item para la API de Softpymes
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

// CreateInvoiceRequest datos tipados para crear una factura en Softpymes
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
	// IsRetry indica que es un reintento: activa la consulta de idempotencia en Softpymes
	IsRetry bool
	// OrderCreatedAt es la fecha de creación de la orden (YYYY-MM-DD, zona Colombia).
	// Se usa como dateFrom en la consulta de idempotencia. Si está vacío se usa la fecha actual.
	OrderCreatedAt string
}

// AuditData captura el request/response HTTP para logging y debugging
type AuditData struct {
	RequestURL     string
	RequestPayload interface{}
	ResponseStatus int
	ResponseBody   string
}

// CreateInvoiceResult resultado de crear una factura en Softpymes
// Se retorna siempre (incluso en error) para incluir audit data
type CreateInvoiceResult struct {
	InvoiceNumber string
	ExternalID    string
	IssuedAt      string
	ProviderInfo  map[string]interface{}
	AuditData     *AuditData
}
