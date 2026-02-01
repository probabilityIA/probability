package response

import (
	"time"

	"gorm.io/datatypes"
)

// Order representa la respuesta HTTP de una orden
// ✅ DTO HTTP - CON TAGS (json + datatypes.JSON)
type Order struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Identificadores de integración
	BusinessID         *uint   `json:"business_id"`
	IntegrationID      uint    `json:"integration_id"`
	IntegrationType    string  `json:"integration_type"`
	IntegrationLogoURL *string `json:"integration_logo_url,omitempty"`

	// Identificadores de la orden
	Platform       string `json:"platform"`
	ExternalID     string `json:"external_id"`
	OrderNumber    string `json:"order_number"`
	InternalNumber string `json:"internal_number"`

	// Información financiera
	Subtotal     float64  `json:"subtotal"`
	Tax          float64  `json:"tax"`
	Discount     float64  `json:"discount"`
	ShippingCost float64  `json:"shipping_cost"`
	TotalAmount  float64  `json:"total_amount"`
	Currency     string   `json:"currency"`
	CodTotal     *float64 `json:"cod_total,omitempty"`

	// Precios en moneda presentment
	SubtotalPresentment     float64 `json:"subtotal_presentment,omitempty"`
	TaxPresentment          float64 `json:"tax_presentment,omitempty"`
	DiscountPresentment     float64 `json:"discount_presentment,omitempty"`
	ShippingCostPresentment float64 `json:"shipping_cost_presentment,omitempty"`
	TotalAmountPresentment  float64 `json:"total_amount_presentment,omitempty"`
	CurrencyPresentment     string  `json:"currency_presentment,omitempty"`

	// Información del cliente
	CustomerID    *uint  `json:"customer_id,omitempty"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	CustomerDNI   string `json:"customer_dni"`

	// Dirección de envío
	ShippingStreet     string   `json:"shipping_street"`
	ShippingCity       string   `json:"shipping_city"`
	ShippingState      string   `json:"shipping_state"`
	ShippingCountry    string   `json:"shipping_country"`
	ShippingPostalCode string   `json:"shipping_postal_code"`
	ShippingLat        *float64 `json:"shipping_lat,omitempty"`
	ShippingLng        *float64 `json:"shipping_lng,omitempty"`

	// Información de pago
	PaymentMethodID uint       `json:"payment_method_id"`
	IsPaid          bool       `json:"is_paid"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`

	// Información de envío/logística
	TrackingNumber      *string    `json:"tracking_number,omitempty"`
	TrackingLink        *string    `json:"tracking_link,omitempty"`
	GuideID             *string    `json:"guide_id,omitempty"`
	GuideLink           *string    `json:"guide_link,omitempty"`
	DeliveryDate        *time.Time `json:"delivery_date,omitempty"`
	DeliveredAt         *time.Time `json:"delivered_at,omitempty"`
	DeliveryProbability *float64   `json:"delivery_probability,omitempty"`

	// Información de fulfillment
	WarehouseID   *uint  `json:"warehouse_id,omitempty"`
	WarehouseName string `json:"warehouse_name"`
	DriverID      *uint  `json:"driver_id,omitempty"`
	DriverName    string `json:"driver_name"`
	IsLastMile    bool   `json:"is_last_mile"`

	// Dimensiones y peso
	Weight *float64 `json:"weight,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Width  *float64 `json:"width,omitempty"`
	Length *float64 `json:"length,omitempty"`
	Boxes  *string  `json:"boxes,omitempty"`

	// Tipo y estado
	OrderTypeID    *uint            `json:"order_type_id,omitempty"`
	OrderTypeName  string           `json:"order_type_name"`
	Status         string           `json:"status"`
	OriginalStatus string           `json:"original_status"`
	StatusID       *uint            `json:"status_id,omitempty"`
	OrderStatus    *OrderStatusInfo `json:"order_status,omitempty"`

	// Estados independientes
	PaymentStatusID     *uint                  `json:"payment_status_id,omitempty"`
	FulfillmentStatusID *uint                  `json:"fulfillment_status_id,omitempty"`
	PaymentStatus       *PaymentStatusInfo     `json:"payment_status,omitempty"`
	FulfillmentStatus   *FulfillmentStatusInfo `json:"fulfillment_status,omitempty"`

	// Información adicional
	Notes    *string `json:"notes,omitempty"`
	Coupon   *string `json:"coupon,omitempty"`
	Approved *bool   `json:"approved,omitempty"`
	UserID   *uint   `json:"user_id,omitempty"`
	UserName string  `json:"user_name"`

	// Novedades
	IsConfirmed *bool   `json:"is_confirmed"`
	Novelty     *string `json:"novelty"`

	// Facturación
	Invoiceable     bool    `json:"invoiceable"`
	InvoiceURL      *string `json:"invoice_url,omitempty"`
	InvoiceID       *string `json:"invoice_id,omitempty"`
	InvoiceProvider *string `json:"invoice_provider,omitempty"`

	// Enlaces Externos
	OrderStatusURL string `json:"order_status_url,omitempty"`

	// Datos estructurados (JSONB) - usando datatypes.JSON
	Items              datatypes.JSON `json:"items,omitempty"`
	Metadata           datatypes.JSON `json:"metadata,omitempty"`
	FinancialDetails   datatypes.JSON `json:"financial_details,omitempty"`
	ShippingDetails    datatypes.JSON `json:"shipping_details,omitempty"`
	PaymentDetails     datatypes.JSON `json:"payment_details,omitempty"`
	FulfillmentDetails datatypes.JSON `json:"fulfillment_details,omitempty"`
	NegativeFactors    []string       `json:"negative_factors,omitempty"`

	// Timestamps
	OccurredAt time.Time `json:"occurred_at"`
	ImportedAt time.Time `json:"imported_at"`
}

// OrderStatusInfo contiene información del estado de orden
type OrderStatusInfo struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Color       string `json:"color"`
}

// PaymentStatusInfo contiene información del estado de pago
type PaymentStatusInfo struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Color       string `json:"color"`
}

// FulfillmentStatusInfo contiene información del estado de fulfillment
type FulfillmentStatusInfo struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Color       string `json:"color"`
}
