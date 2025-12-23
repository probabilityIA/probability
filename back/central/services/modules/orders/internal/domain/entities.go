package domain

import (
	"crypto/rand"
	"time"

	"gorm.io/datatypes"
)

// ───────────────────────────────────────────
//
//	ORDER ENTITIES
//
// ───────────────────────────────────────────

// ProbabilityOrder representa una orden que se guarda en la base de datos
type ProbabilityOrder struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Identificadores de integración
	BusinessID         *uint   `json:"business_id"`
	IntegrationID      uint    `json:"integration_id"`
	IntegrationType    string  `json:"integration_type"`
	IntegrationLogoURL *string `json:"integration_logo_url,omitempty"` // URL del logo del tipo de integración

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
	CodTotal     *float64 `json:"cod_total"`

	// Precios en moneda presentment (presentment_money - moneda local)
	SubtotalPresentment     float64 `json:"subtotal_presentment"`
	TaxPresentment          float64 `json:"tax_presentment"`
	DiscountPresentment     float64 `json:"discount_presentment"`
	ShippingCostPresentment float64 `json:"shipping_cost_presentment"`
	TotalAmountPresentment  float64 `json:"total_amount_presentment"`
	CurrencyPresentment     string  `json:"currency_presentment"`

	// Información del cliente
	CustomerID    *uint  `json:"customer_id"`
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
	ShippingLat        *float64 `json:"shipping_lat"`
	ShippingLng        *float64 `json:"shipping_lng"`

	// Información de pago
	PaymentMethodID uint       `json:"payment_method_id"`
	IsPaid          bool       `json:"is_paid"`
	PaidAt          *time.Time `json:"paid_at"`

	// Información de envío/logística
	TrackingNumber      *string    `json:"tracking_number"`
	TrackingLink        *string    `json:"tracking_link"`
	GuideID             *string    `json:"guide_id"`
	GuideLink           *string    `json:"guide_link"`
	DeliveryDate        *time.Time `json:"delivery_date"`
	DeliveredAt         *time.Time `json:"delivered_at"`
	DeliveryProbability *float64   `json:"delivery_probability"`

	// Información de fulfillment
	WarehouseID   *uint  `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"`
	DriverID      *uint  `json:"driver_id"`
	DriverName    string `json:"driver_name"`
	IsLastMile    bool   `json:"is_last_mile"`

	// Dimensiones y peso
	Weight *float64 `json:"weight"`
	Height *float64 `json:"height"`
	Width  *float64 `json:"width"`
	Length *float64 `json:"length"`
	Boxes  *string  `json:"boxes"`

	// Tipo y estado
	OrderTypeID    *uint            `json:"order_type_id"`
	OrderTypeName  string           `json:"order_type_name"`
	Status         string           `json:"status"`
	OriginalStatus string           `json:"original_status"`
	StatusID       *uint            `json:"status_id"`              // ID del estado mapeado en Probability (FK a order_statuses)
	OrderStatus    *OrderStatusInfo `json:"order_status,omitempty"` // Información del estado de Probability si está cargado

	// Estados independientes
	PaymentStatusID     *uint                  `json:"payment_status_id"`     // FK a payment_statuses
	FulfillmentStatusID *uint                  `json:"fulfillment_status_id"` // FK a fulfillment_statuses
	PaymentStatus       *PaymentStatusInfo     `json:"payment_status,omitempty"`
	FulfillmentStatus   *FulfillmentStatusInfo `json:"fulfillment_status,omitempty"`

	// Información adicional
	Notes    *string `json:"notes"`
	Coupon   *string `json:"coupon"`
	Approved *bool   `json:"approved"`
	UserID   *uint   `json:"user_id"`
	UserName string  `json:"user_name"`

	// Novedades
	IsConfirmed bool    `json:"is_confirmed"`
	Novelty     *string `json:"novelty"`

	// Facturación
	Invoiceable     bool    `json:"invoiceable"`
	InvoiceURL      *string `json:"invoice_url"`
	InvoiceID       *string `json:"invoice_id"`
	InvoiceProvider *string `json:"invoice_provider"`

	// Enlaces Externos
	OrderStatusURL string `json:"order_status_url,omitempty"`

	// Datos estructurados (JSONB)
	Items              datatypes.JSON `json:"items"`
	Metadata           datatypes.JSON `json:"metadata"`
	FinancialDetails   datatypes.JSON `json:"financial_details"`
	ShippingDetails    datatypes.JSON `json:"shipping_details"`
	PaymentDetails     datatypes.JSON `json:"payment_details"`
	FulfillmentDetails datatypes.JSON `json:"fulfillment_details"`

	// Timestamps
	OccurredAt time.Time `json:"occurred_at"`
	ImportedAt time.Time `json:"imported_at"`

	// Relaciones
	OrderItems      []ProbabilityOrderItem            `json:"order_items"`
	Addresses       []ProbabilityAddress              `json:"addresses"`
	Payments        []ProbabilityPayment              `json:"payments"`
	Shipments       []ProbabilityShipment             `json:"shipments"`
	ChannelMetadata []ProbabilityOrderChannelMetadata `json:"channel_metadata"`
	NegativeFactors datatypes.JSON                    `json:"negative_factors"`

	// Campos auxiliares para cálculo de score (No persistir)
	CustomerOrderCount int    `json:"-" gorm:"-"`
	Address2           string `json:"-" gorm:"-"`
}

// ProbabilityOrderItem representa un item de la orden que se guarda en la base de datos
type ProbabilityOrderItem struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	OrderID string `json:"order_id"`

	ProductID    *string `json:"product_id"`
	ProductSKU   string  `json:"product_sku"`
	ProductName  string  `json:"product_name"`
	ProductTitle string  `json:"product_title"`
	VariantID    *string `json:"variant_id"`

	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
	Currency   string  `json:"currency"`

	Discount float64  `json:"discount"`
	Tax      float64  `json:"tax"`
	TaxRate  *float64 `json:"tax_rate"`

	// Precios en moneda presentment (presentment_money - moneda local)
	UnitPricePresentment  float64 `json:"unit_price_presentment"`
	TotalPricePresentment float64 `json:"total_price_presentment"`
	DiscountPresentment   float64 `json:"discount_presentment"`
	TaxPresentment        float64 `json:"tax_presentment"`

	ImageURL          *string        `json:"image_url"`
	ProductURL        *string        `json:"product_url"`
	Weight            *float64       `json:"weight"`
	RequiresShipping  bool           `json:"requires_shipping"`
	IsGiftCard        bool           `json:"is_gift_card"`
	FulfillmentStatus *string        `json:"fulfillment_status"`
	Metadata          datatypes.JSON `json:"metadata"`
}

// ProbabilityAddress representa una dirección que se guarda en la base de datos
type ProbabilityAddress struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Type    string `json:"type"`
	OrderID string `json:"order_id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Phone     string `json:"phone"`

	Street     string `json:"street"`
	Street2    string `json:"street2"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`

	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`

	Instructions *string        `json:"instructions"`
	IsDefault    bool           `json:"is_default"`
	Metadata     datatypes.JSON `json:"metadata"`
}

// ProbabilityPayment representa un pago que se guarda en la base de datos
type ProbabilityPayment struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	OrderID         string `json:"order_id"`
	PaymentMethodID uint   `json:"payment_method_id"`

	Amount       float64  `json:"amount"`
	Currency     string   `json:"currency"`
	ExchangeRate *float64 `json:"exchange_rate"`

	Status      string     `json:"status"`
	PaidAt      *time.Time `json:"paid_at"`
	ProcessedAt *time.Time `json:"processed_at"`

	TransactionID    *string `json:"transaction_id"`
	PaymentReference *string `json:"payment_reference"`
	Gateway          *string `json:"gateway"`

	RefundAmount  *float64       `json:"refund_amount"`
	RefundedAt    *time.Time     `json:"refunded_at"`
	FailureReason *string        `json:"failure_reason"`
	Metadata      datatypes.JSON `json:"metadata"`
}

// ProbabilityShipment representa un envío que se guarda en la base de datos
type ProbabilityShipment struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	OrderID string `json:"order_id"`

	TrackingNumber *string `json:"tracking_number"`
	TrackingURL    *string `json:"tracking_url"`
	Carrier        *string `json:"carrier"`
	CarrierCode    *string `json:"carrier_code"`

	GuideID  *string `json:"guide_id"`
	GuideURL *string `json:"guide_url"`

	Status      string     `json:"status"`
	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`

	ShippingAddressID *uint `json:"shipping_address_id"`

	ShippingCost  *float64 `json:"shipping_cost"`
	InsuranceCost *float64 `json:"insurance_cost"`
	TotalCost     *float64 `json:"total_cost"`

	Weight *float64 `json:"weight"`
	Height *float64 `json:"height"`
	Width  *float64 `json:"width"`
	Length *float64 `json:"length"`

	WarehouseID   *uint  `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"`
	DriverID      *uint  `json:"driver_id"`
	DriverName    string `json:"driver_name"`
	IsLastMile    bool   `json:"is_last_mile"`

	EstimatedDelivery *time.Time     `json:"estimated_delivery"`
	DeliveryNotes     *string        `json:"delivery_notes"`
	Metadata          datatypes.JSON `json:"metadata"`
}

// ProbabilityOrderChannelMetadata representa metadata del canal que se guarda en la base de datos
type ProbabilityOrderChannelMetadata struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	OrderID string `json:"order_id"`

	ChannelSource string `json:"channel_source"`
	IntegrationID uint   `json:"integration_id"`

	RawData datatypes.JSON `json:"raw_data"`

	Version     string     `json:"version"`
	ReceivedAt  time.Time  `json:"received_at"`
	ProcessedAt *time.Time `json:"processed_at"`
	IsLatest    bool       `json:"is_latest"`

	LastSyncedAt *time.Time `json:"last_synced_at"`
	SyncStatus   string     `json:"sync_status"`
}

// ───────────────────────────────────────────
//
//	CATALOG ENTITIES
//
// ───────────────────────────────────────────

// Product representa un producto en el dominio
type Product struct {
	ID         string     `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	BusinessID uint       `json:"business_id"`
	SKU        string     `json:"sku"`
	Name       string     `json:"name"`
	ExternalID string     `json:"external_id"`
}

// Client representa un cliente en el dominio
type Client struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	BusinessID uint       `json:"business_id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	Dni        *string    `json:"dni"`
}

// ToDomainProduct convierte un modelo de BD a dominio
func ToDomainProduct(p interface{}) *Product {
	// Nota: Esto se implementará correctamente en el mapper,
	// aquí solo definimos la estructura.
	return &Product{}
}

// ───────────────────────────────────────────
//
//	ORDER STATUS
//
// ───────────────────────────────────────────

// OrderStatus define los posibles estados de una orden en Probability
type OrderStatus string

const (
	// OrderStatusPending - Orden recibida, pendiente de procesamiento
	OrderStatusPending OrderStatus = "pending"

	// OrderStatusProcessing - Orden en proceso de preparación
	OrderStatusProcessing OrderStatus = "processing"

	// OrderStatusCompleted - Orden completada exitosamente
	OrderStatusCompleted OrderStatus = "completed"

	// OrderStatusCancelled - Orden cancelada
	OrderStatusCancelled OrderStatus = "cancelled"

	// OrderStatusFailed - Orden fallida
	OrderStatusFailed OrderStatus = "failed"

	// OrderStatusRefunded - Orden reembolsada
	OrderStatusRefunded OrderStatus = "refunded"

	// OrderStatusOnHold - Orden en espera
	OrderStatusOnHold OrderStatus = "on_hold"

	// OrderStatusShipped - Orden enviada
	OrderStatusShipped OrderStatus = "shipped"

	// OrderStatusDelivered - Orden entregada
	OrderStatusDelivered OrderStatus = "delivered"
)

// IsValid verifica si el estado es válido
func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusProcessing, OrderStatusCompleted,
		OrderStatusCancelled, OrderStatusFailed, OrderStatusRefunded,
		OrderStatusOnHold, OrderStatusShipped, OrderStatusDelivered:
		return true
	}
	return false
}

// String retorna la representación en string del estado
func (s OrderStatus) String() string {
	return string(s)
}

// CanTransitionTo verifica si se puede transicionar al estado objetivo
func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
	// Definir las transiciones válidas
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending: {
			OrderStatusProcessing,
			OrderStatusCancelled,
			OrderStatusOnHold,
		},
		OrderStatusProcessing: {
			OrderStatusCompleted,
			OrderStatusCancelled,
			OrderStatusOnHold,
			OrderStatusShipped,
		},
		OrderStatusOnHold: {
			OrderStatusPending,
			OrderStatusProcessing,
			OrderStatusCancelled,
		},
		OrderStatusShipped: {
			OrderStatusDelivered,
			OrderStatusFailed,
		},
		OrderStatusDelivered: {
			OrderStatusRefunded,
		},
		OrderStatusCompleted: {
			OrderStatusRefunded,
		},
	}

	allowedTargets, exists := validTransitions[s]
	if !exists {
		return false
	}

	for _, allowed := range allowedTargets {
		if allowed == target {
			return true
		}
	}
	return false
}

// ───────────────────────────────────────────
//
//	ORDER EVENTS
//
// ───────────────────────────────────────────

// OrderEventType define los tipos de eventos relacionados con órdenes
type OrderEventType string

const (
	// Eventos de ciclo de vida de la orden
	OrderEventTypeCreated         OrderEventType = "order.created"
	OrderEventTypeUpdated         OrderEventType = "order.updated"
	OrderEventTypeStatusChanged   OrderEventType = "order.status_changed"
	OrderEventTypeCancelled       OrderEventType = "order.cancelled"
	OrderEventTypeDelivered       OrderEventType = "order.delivered"
	OrderEventTypeShipped         OrderEventType = "order.shipped"
	OrderEventTypePaymentReceived OrderEventType = "order.payment_received"
	OrderEventTypeRefunded        OrderEventType = "order.refunded"
	OrderEventTypeFailed          OrderEventType = "order.failed"
	OrderEventTypeOnHold          OrderEventType = "order.on_hold"
	OrderEventTypeProcessing      OrderEventType = "order.processing"

	// Eventos de cálculo de score
	OrderEventTypeScoreCalculationRequested OrderEventType = "order.score_calculation_requested"
)

// OrderEvent representa un evento relacionado con una orden
type OrderEvent struct {
	ID            string                 `json:"id"`
	Type          OrderEventType         `json:"type"`
	OrderID       string                 `json:"order_id"`
	BusinessID    *uint                  `json:"business_id,omitempty"`
	IntegrationID *uint                  `json:"integration_id,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          OrderEventData         `json:"data"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// OrderEventData contiene los datos específicos del evento de orden
type OrderEventData struct {
	// Información básica de la orden
	OrderNumber    string `json:"order_number,omitempty"`
	InternalNumber string `json:"internal_number,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`

	// Cambios de estado
	PreviousStatus string `json:"previous_status,omitempty"`
	CurrentStatus  string `json:"current_status,omitempty"`

	// Información adicional
	CustomerEmail string                 `json:"customer_email,omitempty"`
	TotalAmount   *float64               `json:"total_amount,omitempty"`
	Currency      string                 `json:"currency,omitempty"`
	Platform      string                 `json:"platform,omitempty"`
	Extra         map[string]interface{} `json:"extra,omitempty"`
}

// NewOrderEvent crea un nuevo evento de orden
func NewOrderEvent(eventType OrderEventType, orderID string, data OrderEventData) *OrderEvent {
	return &OrderEvent{
		ID:        generateEventID(),
		Type:      eventType,
		OrderID:   orderID,
		Timestamp: time.Now(),
		Data:      data,
		Metadata:  make(map[string]interface{}),
	}
}

// generateEventID genera un ID único para el evento
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString genera una cadena aleatoria
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

// ───────────────────────────────────────────
//
//	ORDER ERROR ENTITY
//
// ───────────────────────────────────────────

// OrderError representa un error ocurrido durante el procesamiento de una orden
type OrderError struct {
	ID              uint
	ExternalID      string
	IntegrationID   uint
	BusinessID      *uint
	IntegrationType string
	Platform        string
	ErrorType       string
	ErrorMessage    string
	ErrorStack      *string
	RawData         datatypes.JSON
	Status          string
	ResolvedAt      *time.Time
	ResolvedBy      *uint
	Resolution      *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
