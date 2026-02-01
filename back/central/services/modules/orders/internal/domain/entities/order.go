package entities

import "time"

// ProbabilityOrder representa una orden que se guarda en la base de datos
// ✅ ENTIDAD PURA - SIN TAGS (ni json, ni gorm, ni validate)
type ProbabilityOrder struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Identificadores de integración
	BusinessID         *uint
	IntegrationID      uint
	IntegrationType    string
	IntegrationLogoURL *string

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
	CustomerID    *uint
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	CustomerDNI   string

	// Dirección de envío
	ShippingStreet     string
	ShippingCity       string
	ShippingState      string
	ShippingCountry    string
	ShippingPostalCode string
	ShippingLat        *float64
	ShippingLng        *float64

	// Información de pago
	PaymentMethodID uint
	IsPaid          bool
	PaidAt          *time.Time

	// Información de envío/logística
	TrackingNumber      *string
	TrackingLink        *string
	GuideID             *string
	GuideLink           *string
	DeliveryDate        *time.Time
	DeliveredAt         *time.Time
	DeliveryProbability *float64

	// Información de fulfillment
	WarehouseID   *uint
	WarehouseName string
	DriverID      *uint
	DriverName    string
	IsLastMile    bool

	// Dimensiones y peso
	Weight *float64
	Height *float64
	Width  *float64
	Length *float64
	Boxes  *string

	// Tipo y estado
	OrderTypeID    *uint
	OrderTypeName  string
	Status         string
	OriginalStatus string
	StatusID       *uint
	OrderStatus    *OrderStatusInfo

	// Estados independientes
	PaymentStatusID     *uint
	FulfillmentStatusID *uint
	PaymentStatus       *PaymentStatusInfo
	FulfillmentStatus   *FulfillmentStatusInfo

	// Información adicional
	Notes    *string
	Coupon   *string
	Approved *bool
	UserID   *uint
	UserName string

	// Novedades
	IsConfirmed *bool
	Novelty     *string

	// Facturación
	Invoiceable     bool
	InvoiceURL      *string
	InvoiceID       *string
	InvoiceProvider *string

	// Enlaces Externos
	OrderStatusURL string

	// Datos estructurados (JSONB) - almacenados como []byte
	Items              []byte
	Metadata           []byte
	FinancialDetails   []byte
	ShippingDetails    []byte
	PaymentDetails     []byte
	FulfillmentDetails []byte

	// Timestamps
	OccurredAt time.Time
	ImportedAt time.Time

	// Relaciones
	OrderItems      []ProbabilityOrderItem
	Addresses       []ProbabilityAddress
	Payments        []ProbabilityPayment
	Shipments       []ProbabilityShipment
	ChannelMetadata []ProbabilityOrderChannelMetadata
	NegativeFactors []byte

	// Campos auxiliares para cálculo de score
	CustomerOrderCount int
	CustomerTotalSpent string
	Address2           string
}

// OrderStatusInfo contiene información básica del estado de orden de Probability
type OrderStatusInfo struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	Color       string
}

// PaymentStatusInfo contiene información del estado de pago
type PaymentStatusInfo struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	Color       string
}

// FulfillmentStatusInfo contiene información del estado de fulfillment
type FulfillmentStatusInfo struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	Color       string
}
