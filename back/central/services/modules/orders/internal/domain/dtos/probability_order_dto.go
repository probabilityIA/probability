package dtos

import "time"

// ProbabilityOrderDTO representa la estructura de orden en la lógica de negocio que todas las integraciones
// deben enviar después de mapear sus datos específicos
// ✅ DTO PURO - SIN TAGS (ni json, ni binding, ni gorm)
type ProbabilityOrderDTO struct {
	// Identificadores de integración
	BusinessID      *uint
	IntegrationID   uint
	IntegrationType string

	// Identificadores de la orden
	Platform       string
	ExternalID     string
	OrderNumber    string
	InternalNumber string

	// Información financiera
	Subtotal     float64
	Tax          float64
	Discount     float64
	ShippingCost float64
	TotalAmount  float64
	Currency     string
	CodTotal     *float64

	// Precios en moneda presentment (presentment_money - moneda local)
	SubtotalPresentment     float64
	TaxPresentment          float64
	DiscountPresentment     float64
	ShippingCostPresentment float64
	TotalAmountPresentment  float64
	CurrencyPresentment     string

	// Información del cliente
	CustomerID         *uint
	CustomerName       string
	CustomerEmail      string
	CustomerPhone      string
	CustomerDNI        string
	CustomerOrderCount *int
	CustomerTotalSpent *string

	// Tipo y estado
	OrderTypeID    *uint
	OrderTypeName  string
	Status         string
	OriginalStatus string
	StatusID       *uint

	// Estados independientes
	PaymentStatusID     *uint
	FulfillmentStatusID *uint

	// Información adicional
	Notes    *string
	Coupon   *string
	Approved *bool
	UserID   *uint
	UserName string

	// Facturación
	Invoiceable     bool
	InvoiceURL      *string
	InvoiceID       *string
	InvoiceProvider *string

	// Enlaces Externos
	OrderStatusURL string

	// Timestamps
	OccurredAt time.Time
	ImportedAt time.Time

	// Datos estructurados (JSONB) - Para compatibilidad - almacenados como []byte
	Items              []byte
	Metadata           []byte
	FinancialDetails   []byte
	ShippingDetails    []byte
	PaymentDetails     []byte
	FulfillmentDetails []byte

	// Tablas relacionadas
	OrderItems      []ProbabilityOrderItemDTO
	Addresses       []ProbabilityAddressDTO
	Payments        []ProbabilityPaymentDTO
	Shipments       []ProbabilityShipmentDTO
	ChannelMetadata *ProbabilityChannelMetadataDTO
}

// ProbabilityOrderItemDTO representa un item/producto de la orden
type ProbabilityOrderItemDTO struct {
	ProductID    *string
	ProductSKU   string
	ProductName  string
	ProductTitle string
	VariantID    *string
	Quantity     int
	UnitPrice    float64
	TotalPrice   float64
	Currency     string
	Discount     float64
	Tax          float64
	TaxRate      *float64

	// Precios en moneda presentment (presentment_money - moneda local)
	UnitPricePresentment  float64
	TotalPricePresentment float64
	DiscountPresentment   float64
	TaxPresentment        float64
	ImageURL              *string
	ProductURL            *string
	Weight                *float64
	Metadata              []byte
}

// ProbabilityAddressDTO representa una dirección (envío o facturación)
type ProbabilityAddressDTO struct {
	Type         string
	FirstName    string
	LastName     string
	Company      string
	Phone        string
	Street       string
	Street2      string
	City         string
	State        string
	Country      string
	PostalCode   string
	Latitude     *float64
	Longitude    *float64
	Instructions *string
	Metadata     []byte
}

// ProbabilityPaymentDTO representa un pago de la orden
type ProbabilityPaymentDTO struct {
	PaymentMethodID  uint
	Amount           float64
	Currency         string
	ExchangeRate     *float64
	Status           string
	PaidAt           *time.Time
	ProcessedAt      *time.Time
	TransactionID    *string
	PaymentReference *string
	Gateway          *string
	RefundAmount     *float64
	RefundedAt       *time.Time
	FailureReason    *string
	Metadata         []byte
}

// ProbabilityShipmentDTO representa un envío de la orden
type ProbabilityShipmentDTO struct {
	TrackingNumber    *string
	TrackingURL       *string
	Carrier           *string
	CarrierCode       *string
	GuideID           *string
	GuideURL          *string
	Status            string
	ShippedAt         *time.Time
	DeliveredAt       *time.Time
	ShippingAddressID *uint
	ShippingCost      *float64
	InsuranceCost     *float64
	TotalCost         *float64
	Weight            *float64
	Height            *float64
	Width             *float64
	Length            *float64
	WarehouseID       *uint
	WarehouseName     string
	DriverID          *uint
	DriverName        string
	IsLastMile        bool
	EstimatedDelivery *time.Time
	DeliveryNotes     *string
	Metadata          []byte
}

// ProbabilityChannelMetadataDTO representa los datos crudos del canal
type ProbabilityChannelMetadataDTO struct {
	ChannelSource string
	RawData       []byte
	Version       string
	ReceivedAt    time.Time
	ProcessedAt   *time.Time
	IsLatest      bool
	LastSyncedAt  *time.Time
	SyncStatus    string
}
