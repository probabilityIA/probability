package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

// ═══════════════════════════════════════════════════════════════
// CLIENTE DE SIIGO (Secondary Port - Driven Adapter)
// ═══════════════════════════════════════════════════════════════

// ISiigoClient define las operaciones con la API de Siigo
type ISiigoClient interface {
	// TestAuthentication verifica que las credenciales sean válidas
	// baseURL es opcional: si está vacío usa la URL configurada en el cliente
	TestAuthentication(ctx context.Context, username, accessKey, accountID, partnerID, baseURL string) error

	// CreateInvoice crea una factura en Siigo
	// Retorna resultado con datos de la factura y audit data (incluso en caso de error)
	CreateInvoice(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error)

	// GetCustomerByIdentification busca un cliente en Siigo por número de identificación
	// Endpoint: GET /v1/customers?identification=xxx
	GetCustomerByIdentification(ctx context.Context, credentials dtos.Credentials, identification string) (*dtos.CustomerResult, error)

	// CreateCustomer crea un cliente en Siigo
	// Endpoint: POST /v1/customers
	CreateCustomer(ctx context.Context, credentials dtos.Credentials, req *dtos.CreateCustomerRequest) (*dtos.CustomerResult, error)

	// ListInvoices consulta la lista paginada de facturas emitidas en Siigo
	// Endpoint: GET /v1/invoices
	ListInvoices(ctx context.Context, credentials dtos.Credentials, params dtos.ListInvoicesParams) (*dtos.ListInvoicesResult, error)
}

// ═══════════════════════════════════════════════════════════════
// USE CASE DE FACTURACIÓN AUTOMÁTICA (Primary Port - Driving Adapter)
// ═══════════════════════════════════════════════════════════════

// IInvoiceUseCase define el caso de uso para procesar órdenes y crear facturas automáticamente
type IInvoiceUseCase interface {
	// ProcessOrderForInvoicing procesa un evento de orden para determinar si debe facturarse
	ProcessOrderForInvoicing(ctx context.Context, event *OrderEventMessage) error
}

// ═══════════════════════════════════════════════════════════════
// ESTRUCTURAS DE EVENTOS (Replicadas localmente)
// ═══════════════════════════════════════════════════════════════

// OrderEventMessage representa el payload de eventos de órdenes en RabbitMQ
type OrderEventMessage struct {
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	OrderID       string                 `json:"order_id"`
	BusinessID    *uint                  `json:"business_id"`
	IntegrationID *uint                  `json:"integration_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Order         *OrderSnapshot         `json:"order"`
	Changes       map[string]interface{} `json:"changes,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// OrderSnapshot representa un snapshot completo de una orden
type OrderSnapshot struct {
	ID              string              `json:"id"`
	OrderNumber     string              `json:"order_number"`
	InternalNumber  string              `json:"internal_number"`
	ExternalID      string              `json:"external_id"`
	TotalAmount     float64             `json:"total_amount"`
	Currency        string              `json:"currency"`
	PaymentMethodID uint                `json:"payment_method_id"`
	PaymentStatusID *uint               `json:"payment_status_id,omitempty"`
	Subtotal        float64             `json:"subtotal"`
	Tax             float64             `json:"tax"`
	Discount        float64             `json:"discount"`
	ShippingCost    float64             `json:"shipping_cost"`
	CustomerName    string              `json:"customer_name"`
	CustomerEmail   string              `json:"customer_email,omitempty"`
	CustomerPhone   string              `json:"customer_phone,omitempty"`
	CustomerDNI     string              `json:"customer_dni,omitempty"`
	Platform        string              `json:"platform"`
	IntegrationID   uint                `json:"integration_id"`
	Items           []OrderItemSnapshot `json:"items,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

// OrderItemSnapshot representa un item de orden
type OrderItemSnapshot struct {
	ProductID  *string  `json:"product_id,omitempty"`
	SKU        string   `json:"sku"`
	Name       string   `json:"name"`
	Quantity   int      `json:"quantity"`
	UnitPrice  float64  `json:"unit_price"`
	TotalPrice float64  `json:"total_price"`
	Tax        float64  `json:"tax"`
	TaxRate    *float64 `json:"tax_rate,omitempty"`
	Discount   float64  `json:"discount"`
}
