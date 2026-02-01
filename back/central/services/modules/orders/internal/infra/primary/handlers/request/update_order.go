package request

import (
	"time"

	"gorm.io/datatypes"
)

// UpdateOrder representa la petición HTTP para actualizar una orden
// ✅ DTO HTTP - CON TAGS (json + binding)
type UpdateOrder struct {
	// Información financiera
	Subtotal     *float64 `json:"subtotal" binding:"omitempty,min=0"`
	Tax          *float64 `json:"tax" binding:"omitempty,min=0"`
	Discount     *float64 `json:"discount" binding:"omitempty,min=0"`
	ShippingCost *float64 `json:"shipping_cost" binding:"omitempty,min=0"`
	TotalAmount  *float64 `json:"total_amount" binding:"omitempty,min=0"`
	Currency     *string  `json:"currency" binding:"omitempty,max=10"`
	CodTotal     *float64 `json:"cod_total"`

	// Información del cliente
	CustomerName       *string `json:"customer_name" binding:"omitempty,max=255"`
	CustomerEmail      *string `json:"customer_email" binding:"omitempty,max=255"`
	CustomerPhone      *string `json:"customer_phone" binding:"omitempty,max=32"`
	CustomerDNI        *string `json:"customer_dni" binding:"omitempty,max=64"`
	CustomerOrderCount *int    `json:"customer_order_count"`
	CustomerTotalSpent *string `json:"customer_total_spent"`

	// Dirección de envío
	ShippingStreet     *string  `json:"shipping_street" binding:"omitempty,max=255"`
	ShippingCity       *string  `json:"shipping_city" binding:"omitempty,max=128"`
	ShippingState      *string  `json:"shipping_state" binding:"omitempty,max=128"`
	ShippingCountry    *string  `json:"shipping_country" binding:"omitempty,max=128"`
	ShippingPostalCode *string  `json:"shipping_postal_code" binding:"omitempty,max=32"`
	ShippingLat        *float64 `json:"shipping_lat"`
	ShippingLng        *float64 `json:"shipping_lng"`

	// Información de pago
	PaymentMethodID *uint      `json:"payment_method_id"`
	IsPaid          *bool      `json:"is_paid"`
	PaidAt          *time.Time `json:"paid_at"`

	// Información de envío/logística
	TrackingNumber *string    `json:"tracking_number"`
	TrackingLink   *string    `json:"tracking_link"`
	GuideID        *string    `json:"guide_id"`
	GuideLink      *string    `json:"guide_link"`
	DeliveryDate   *time.Time `json:"delivery_date"`
	DeliveredAt    *time.Time `json:"delivered_at"`

	// Información de fulfillment
	WarehouseID   *uint   `json:"warehouse_id"`
	WarehouseName *string `json:"warehouse_name" binding:"omitempty,max=128"`
	DriverID      *uint   `json:"driver_id"`
	DriverName    *string `json:"driver_name" binding:"omitempty,max=255"`
	IsLastMile    *bool   `json:"is_last_mile"`

	// Dimensiones y peso
	Weight *float64 `json:"weight"`
	Height *float64 `json:"height"`
	Width  *float64 `json:"width"`
	Length *float64 `json:"length"`
	Boxes  *string  `json:"boxes"`

	// Tipo y estado
	OrderTypeID    *uint   `json:"order_type_id"`
	OrderTypeName  *string `json:"order_type_name" binding:"omitempty,max=64"`
	Status         *string `json:"status" binding:"omitempty,max=64"`
	OriginalStatus *string `json:"original_status" binding:"omitempty,max=64"`
	StatusID       *uint   `json:"status_id" binding:"omitempty"`

	// Estados independientes
	PaymentStatusID     *uint `json:"payment_status_id" binding:"omitempty"`
	FulfillmentStatusID *uint `json:"fulfillment_status_id" binding:"omitempty"`

	// Información adicional
	Notes    *string `json:"notes"`
	Coupon   *string `json:"coupon"`
	Approved *bool   `json:"approved"`
	UserID   *uint   `json:"user_id"`
	UserName *string `json:"user_name" binding:"omitempty,max=255"`

	// Novedades
	IsConfirmed        *bool   `json:"is_confirmed"`
	ConfirmationStatus *string `json:"confirmation_status"`
	Novelty            *string `json:"novelty"`

	// Facturación
	Invoiceable     *bool   `json:"invoiceable"`
	InvoiceURL      *string `json:"invoice_url"`
	InvoiceID       *string `json:"invoice_id"`
	InvoiceProvider *string `json:"invoice_provider"`

	// Datos estructurados (JSONB)
	Items              datatypes.JSON `json:"items"`
	Metadata           datatypes.JSON `json:"metadata"`
	FinancialDetails   datatypes.JSON `json:"financial_details"`
	ShippingDetails    datatypes.JSON `json:"shipping_details"`
	PaymentDetails     datatypes.JSON `json:"payment_details"`
	FulfillmentDetails datatypes.JSON `json:"fulfillment_details"`
}
