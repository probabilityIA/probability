package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
)

// ═══════════════════════════════════════════════════════════════
// CLIENTE DE SOFTPYMES (Secondary Port - Driven Adapter)
// ═══════════════════════════════════════════════════════════════

// ISoftpymesClient define las operaciones con la API de Softpymes
type ISoftpymesClient interface {
	// TestAuthentication verifica que las credenciales sean válidas
	// referer: Identificación de la instancia del cliente (requerido por API)
	// baseURL: URL base efectiva (producción o testing); vacío usa la URL del constructor
	TestAuthentication(ctx context.Context, apiKey, apiSecret, referer, baseURL string) error

	// CreateInvoice crea una factura en Softpymes
	// Retorna resultado con datos de la factura y audit data (incluso en caso de error)
	// baseURL: URL base efectiva (producción o testing); vacío usa la URL del constructor
	CreateInvoice(ctx context.Context, req *dtos.CreateInvoiceRequest, baseURL string) (*dtos.CreateInvoiceResult, error)

	// CreateCreditNote crea una nota crédito en Softpymes
	CreateCreditNote(ctx context.Context, creditNoteData map[string]interface{}) error

	// GetDocumentByNumber obtiene un documento completo por su número
	// Usado para consulta posterior después de crear factura (esperar procesamiento DIAN)
	// Retorna el documento con todos sus detalles (items, totales, información de envío)
	// baseURL: URL base efectiva (producción o testing); vacío usa la URL del constructor
	GetDocumentByNumber(ctx context.Context, apiKey, apiSecret, referer, documentNumber, baseURL string) (map[string]interface{}, error)
}

// ═══════════════════════════════════════════════════════════════
// USE CASE DE FACTURACIÓN AUTOMÁTICA (Primary Port - Driving Adapter)
// ═══════════════════════════════════════════════════════════════

// IInvoiceUseCase define el caso de uso para procesar órdenes y crear facturas automáticamente
// Este use case es consumido por el OrderConsumer (RabbitMQ) y NO requiere base de datos
type IInvoiceUseCase interface {
	// ProcessOrderForInvoicing procesa un evento de orden para determinar si debe facturarse
	// Valida filtros, verifica duplicados y crea la factura en Softpymes si corresponde
	ProcessOrderForInvoicing(ctx context.Context, event *OrderEventMessage) error

	// TestConnection valida que las credenciales y configuración sean correctas
	// contra la API de Softpymes. Llamado desde el contrato global IIntegrationContract.
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
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
