package dtos

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// OrderResponse representa la respuesta de una orden
// ✅ DTO PURO - SIN TAGS
type OrderResponse struct {
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
	OrderStatus    *entities.OrderStatusInfo

	// Estados independientes
	PaymentStatusID     *uint
	FulfillmentStatusID *uint
	PaymentStatus       *entities.PaymentStatusInfo
	FulfillmentStatus   *entities.FulfillmentStatusInfo

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
	NegativeFactors    []string

	// Timestamps
	OccurredAt time.Time
	ImportedAt time.Time
}
