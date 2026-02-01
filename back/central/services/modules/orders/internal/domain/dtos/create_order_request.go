package dtos

import "time"

// CreateOrderRequest representa la solicitud para crear una orden
// ✅ DTO PURO - SIN TAGS (ni json, ni binding, ni gorm)
type CreateOrderRequest struct {
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
	TrackingNumber *string
	TrackingLink   *string
	GuideID        *string
	GuideLink      *string
	DeliveryDate   *time.Time
	DeliveredAt    *time.Time

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
}
