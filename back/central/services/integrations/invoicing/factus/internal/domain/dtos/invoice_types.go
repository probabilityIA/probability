package dtos

// Credentials contiene las credenciales de autenticación de Factus
type Credentials struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	BaseURL      string // URL base de la API (opcional, usa el default del cliente si está vacío)
}

// CustomerData datos del cliente para la API de Factus
type CustomerData struct {
	Name    string
	Email   string
	Phone   string
	DNI     string
	Address string
}

// ItemData datos de un item para la API de Factus
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

// CreateInvoiceRequest datos tipados para crear una factura en Factus
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

// ProcessInvoiceRequest es el input del caso de uso para procesar una solicitud de facturación.
// No contiene credenciales — el use case las obtiene y descifra desde la base de datos.
type ProcessInvoiceRequest struct {
	InvoiceID     uint
	Operation     string
	CorrelationID string
	IntegrationID uint
	Customer      CustomerData
	Items         []ItemData
	Total         float64
	Subtotal      float64
	Tax           float64
	Discount      float64
	ShippingCost  float64
	Currency      string
	OrderID       string
	Config        map[string]interface{}
}

// ProcessInvoiceResult es el output del caso de uso.
// Se retorna siempre (incluso en error) para propagar el AuditData hacia el consumer.
type ProcessInvoiceResult struct {
	InvoiceNumber string
	ExternalID    string
	CUFE          string
	QRCode        string
	Total         string
	IssuedAt      string
	AuditData     *AuditData
}

// AuditData captura el request/response HTTP para logging y debugging
type AuditData struct {
	RequestURL     string
	RequestPayload interface{}
	ResponseStatus int
	ResponseBody   string
}

// CreateInvoiceResult resultado de crear una factura en Factus
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
