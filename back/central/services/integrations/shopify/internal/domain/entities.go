package domain

import (
	"time"
)

// ShopifyOrderDTO represents the raw structure of an order from Shopify
// This is used for internal mapping and processing, not for DB storage.
type ShopifyOrderDTO struct {
	ID                string                 `json:"id"`
	OrderNumber       string                 `json:"order_number"`
	TotalPrice        float64                `json:"total_price"`
	Currency          string                 `json:"currency"`
	PaymentType       string                 `json:"payment_type"`
	CustomerName      string                 `json:"customer_name"`
	CustomerEmail     string                 `json:"customer_email"`
	Phone             string                 `json:"phone"`
	Country           string                 `json:"country"`
	Province          string                 `json:"province"`
	City              string                 `json:"city"`
	Address           string                 `json:"address"`
	AddressComplement string                 `json:"address_complement"`
	FinancialStatus   string                 `json:"financial_status"`
	FulfillmentStatus string                 `json:"fulfillment_status"`
	CreatedAt         time.Time              `json:"created_at"`
	RawData           map[string]interface{} `json:"raw_data"`
}

// UnifiedOrder represents the standardized order structure for the system
// This is what we publish to the queue for the Orders module to consume.
type UnifiedOrder struct {
	// Identificadores de la integración
	BusinessID      *uint  `json:"business_id"`      // ID del negocio (null = global)
	IntegrationID   uint   `json:"integration_id"`   // ID de la integración específica
	IntegrationType string `json:"integration_type"` // Tipo de integración: "shopify", "whatsapp", etc.

	// Identificadores de la orden
	Platform    string `json:"platform"`     // Plataforma origen: "shopify", "whatsapp", etc.
	ExternalID  string `json:"external_id"`  // ID de la orden en la plataforma externa
	OrderNumber string `json:"order_number"` // Número de orden visible

	// Información financiera
	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`

	// Información del cliente y envío
	Customer        UnifiedCustomer `json:"customer"`
	ShippingAddress UnifiedAddress  `json:"shipping_address"`

	// Estado
	Status         string `json:"status"`          // Estado mapeado interno
	OriginalStatus string `json:"original_status"` // Estado original de la plataforma

	// Items de la orden
	Items []UnifiedOrderItem `json:"items"`

	// Metadata adicional (puede contener datos específicos de cada plataforma)
	Metadata map[string]interface{} `json:"metadata"`

	// Timestamps
	OccurredAt time.Time `json:"occurred_at"` // Cuándo ocurrió la orden en la plataforma
	ImportedAt time.Time `json:"imported_at"` // Cuándo se importó a Probability
}

type UnifiedCustomer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UnifiedAddress struct {
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PostalCode  string `json:"postal_code"`
	Coordinates *struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"coordinates,omitempty"`
}

type UnifiedOrderItem struct {
	ExternalID string  `json:"external_id"`
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}
